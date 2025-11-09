package repository

import (
	"database/sql"

	"github.com/RitwikGupta-0501/vital-watch/internal/models"
)

// Repository is the struct that holds our database connection
type Repository struct {
	DB *sql.DB
}

func (r *Repository) CreatePatient(firstName, lastName, email, hashedPassword string) (int, error) {
	query := `
		INSERT INTO patients (firstName, lastName, email, hashedPassword)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var newID int

	// Pass all 4 arguments
	err := r.DB.QueryRow(query, firstName, lastName, email, hashedPassword).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *Repository) CreateDoctor(firstName, lastName, email, hashedPassword string) (int, error) {
	query := `
		INSERT INTO doctors (firstName, lastName, email, hashedPassword)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var newID int

	err := r.DB.QueryRow(query, firstName, lastName, email, hashedPassword).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *Repository) GetPatientByEmail(email string) (models.Patient, error) {
	query := `SELECT id, firstName, lastName, email, hashedPassword FROM patients WHERE email=$1`

	var user models.Patient

	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.Patient{}, err
	}

	return user, nil
}

func (r *Repository) GetDoctorByEmail(email string) (models.Doctor, error) {
	query := `SELECT id, firstName, lastName, email, hashedPassword FROM doctors WHERE email=$1`

	var user models.Doctor
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.Doctor{}, err
	}

	return user, nil
}

func (r *Repository) GetPatientByID(id int) (models.Patient, error) {
	query := `SELECT id, firstName, lastName, email, hashedPassword FROM patients WHERE id=$1`

	var user models.Patient
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.Patient{}, err
	}

	return user, nil
}

func (r *Repository) GetDoctorByID(id int) (models.Doctor, error) {
	query := `SELECT id, firstName, lastName, email, hashedPassword FROM doctors WHERE id=$1`

	var user models.Doctor
	err := r.DB.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.Doctor{}, err
	}

	return user, nil
}
