package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database!")
	return db
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
	r.Use(cors.Default())

	// -----------------------
	// -       Routes        -
	// -----------------------
	r.GET("/api/ping", h.Ping)
	r.POST("/api/register", h.RegisterUser)
	r.POST("/api/login", h.Login)

	// --- Protected Routes ---
	// Create a group that uses the AuthMiddleware
	authGroup := r.Group("/api")
	authGroup.Use(api.AuthMiddleware()) // Apply the middleware
	{
		// All routes inside this block are now protected

		// authGroup.GET("/patient", h.GetPatient)
		// authGroup.GET("/doctors", h.GetDoctors)
		// authGroup.GET("/appointments", h.GetAppointments)
		// authGroup.GET("/prescriptions", h.GetPrescriptions)
	}

	// Run the server
	r.Run()
}
