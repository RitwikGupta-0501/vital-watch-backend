package repository

import (
	"database/sql"

	"github.com/RitwikGupta-0501/vital-watch/internal/models"
)

// Repository is the struct that holds our database connection
type Repository struct {
	DB *sql.DB
}

func (r *Repository) CreatePatient(user models.Patient) (string, error) {
	query := `INSERT INTO patient (FirstName, LastName, HashedPassword) VALUES ($1, '$2, $3) RETURNING ID`

	var id string

	err := r.DB.QueryRow(query, user.FirstName, user.LastName, user.HashedPassword).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) CreateDoctor(user models.Doctor) (string, error) {
	query := `INSERT INTO doctor (FirstName, LastName, HashedPassword) VALUES ($1, '$2, $3) RETURNING ID`

	var id string

	err := r.DB.QueryRow(query, user.FirstName, user.LastName, user.HashedPassword).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) GetPatientByEmail(email string) (models.Patient, error) {
	query := `SELECT ID, FirstName, LastName, Email, HashedPassword FROM patients WHERE Email=$1`

	var user models.Patient

	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.Patient{}, err
	}

	return user, nil
}

func (r *Repository) GetDoctorByEmail(email string) (models.Doctor, error) {
	query := `SELECT ID, FirstName, LastName, Email, HashedPassword FROM doctors WHERE Email=$1`

	var user models.Doctor
	err := r.DB.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.Doctor{}, err
	}

	return user, nil
}
