package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/RitwikGupta-0501/vital-watch/internal/api"
	"github.com/RitwikGupta-0501/vital-watch/internal/repository"
)

/*
========================================
=        Database Initialization       =
========================================
*/
func init_db() *sql.DB {

	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatal("Invalid DB_PORT:", err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	// Open the database connection pool
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	// --- NEW: Add a retry loop for db.Ping() ---
	var dbErr error
	for i := 0; i < 5; i++ { // Try 5 times
		err = db.Ping()
		if err == nil {
			// Success!
			log.Println("Successfully connected to database!")
			return db
		}
		dbErr = err
		log.Println("Failed to ping database, retrying in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
	// If the loop finishes, we failed
	log.Fatal("Failed to ping database after retries:", dbErr)
	return nil
}

/*
========================================
=           Migrations Runner          =
========================================
*/
func run_migrations(db *sql.DB) {
	log.Println("Running database migrations...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	// Point to the migration files
	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}

	// Run the migrations "up"
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations finished successfully.")
}

/*
========================================
=                Main                  =
========================================
*/
func main() {
	// Initialize environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Initialize DB
	var db = init_db()
	defer db.Close()

	// Run DB migrations
	run_migrations(db)

	// Initialize repository
	repo := &repository.Repository{
		DB: db,
	}

	// Create the API Handler
	h := &api.Handler{
		Repo: repo,
	}

	// Set up Gin Server
	r := gin.Default()

	// Enable CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // TODO: Replace with the fontend's address like []string{"http://localhost:3000"}
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// -----------------------
	// -       Routes        -
	// -----------------------
	r.GET("/api/ping", h.Ping)
	r.POST("/api/register", h.RegisterUser)
	r.POST("/api/login", h.Login)

	// --- Protected Routes ---
	// Create a group that uses the AuthMiddleware
	authGroup := r.Group("/api")

	// All routes inside this block will require authentication
	authGroup.Use(api.AuthMiddleware())
	{
		authGroup.GET("/profile", h.GetUserProfile)
		authGroup.GET("/doctors", h.GetDoctors)
		authGroup.GET("/appointments", h.GetAppointments)
		authGroup.GET("/prescriptions", h.GetPrescriptions)
		authGroup.POST("/appointments", h.CreateAppointment)

		authGroup.GET("/prescriptions/:filename", h.DownloadPrescription)
	}

	// Run the server
	r.Run()
}
