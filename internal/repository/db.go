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
	query := `INSERT INTO users (FirstName, LastName, Password) VALUES ($1, '$2, $3) RETURNING ID`

	var id string

	err := r.DB.QueryRow(query, user.FirstName, user.LastName, user.Password).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
