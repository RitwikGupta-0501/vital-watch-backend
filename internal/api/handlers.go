package api

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/RitwikGupta-0501/vital-watch/internal/models"
	"github.com/RitwikGupta-0501/vital-watch/internal/repository"
	"github.com/RitwikGupta-0501/vital-watch/utils"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Handler struct {
	Repo *repository.Repository
}

func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong from the api layer!"})
}

// --- Authentication Handlers ---
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userIDFloat, ok := claims["sub"].(float64)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}

			// Get User Role
			role, ok := claims["role"].(string)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims (role)"})
				return
			}

			c.Set("userID", int(userIDFloat))
			c.Set("role", role)
		}

		c.Next()
	}
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Role     string `json:"role"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var (
		user models.Authenticatable
		err  error
	)

	switch req.Role {
	case "patient":
		user, err = h.Repo.GetPatientByEmail(req.Email)
	case "doctor":
		user, err = h.Repo.GetDoctorByEmail(req.Email)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.GetHashedPassword()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.GetID(),                              // "subject" (who the token is for)
		"role": req.Role,                                  // custom claim for user role
		"iat":  time.Now().Unix(),                         // "issued at"
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(), // "expires at" (e.g., 7 days)
		"iss":  "vital-watch",                             // "issuer"
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Failed to sign token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var req struct {
		Role      string `json:"role"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var id int
	switch req.Role {
	case "patient":
		id, err = h.Repo.CreatePatient(req.FirstName, req.LastName, req.Email, hashed)
	case "doctor":
		id, err = h.Repo.CreateDoctor(req.FirstName, req.LastName, req.Email, hashed)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Role not found in context"})
		return
	}

	switch role {
	case "patient":
		patient, err := h.Repo.GetPatientByID(userID.(int))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient profile not found"})
			return
		}
		c.JSON(http.StatusOK, patient)

	case "doctor":
		doctor, err := h.Repo.GetDoctorByID(userID.(int))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Doctor profile not found"})
			return
		}
		c.JSON(http.StatusOK, doctor)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user role"})
	}
}

func (h *Handler) GetDoctors(c *gin.Context) {
	doctors, err := h.Repo.GetDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

func (h *Handler) GetAppointments(c *gin.Context) {
	patientID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	appointments, err := h.Repo.GetAppointmentsByPatientID(patientID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
		return
	}
	c.JSON(http.StatusOK, appointments)
}

func (h *Handler) GetPrescriptions(c *gin.Context) {
	patientID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	prescriptions, err := h.Repo.GetPrescriptionsByPatientID(patientID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prescriptions"})
		return
	}
	c.JSON(http.StatusOK, prescriptions)
}

func (h *Handler) CreateAppointment(c *gin.Context) {
	var req struct {
		DoctorID  int       `json:"doctor_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		Type      string    `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	patientID, ok := c.Get("userID")
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	newID, err := h.Repo.CreateAppointment(patientID.(int), req.DoctorID, req.StartTime, req.EndTime, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newID})
}
