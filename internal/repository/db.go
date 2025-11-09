package repository

import (
	"database/sql"
	"time"

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

func (r *Repository) GetDoctors() ([]models.Doctor, error) {
	query := `SELECT id, firstName, lastName, email, specialty, experience, available FROM doctors`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []models.Doctor
	for rows.Next() {
		var doc models.Doctor
		err := rows.Scan(&doc.ID, &doc.FirstName, &doc.LastName, &doc.Email, &doc.Specialty, &doc.Experience, &doc.Available)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, doc)
	}
	return doctors, nil
}

func (r *Repository) GetAppointmentsByPatientID(patientID int) ([]models.Appointment, error) {
	// 1. UPDATED: The query now JOINS the doctors table
	query := `
		SELECT 
			a.id, a.patient_id, a.doctor_id, a.start_time, a.end_time, a.status, a.appointment_type,
			d.firstName, d.lastName, d.specialty
		FROM appointments a
		JOIN doctors d ON a.doctor_id = d.id
		WHERE a.patient_id = $1
		ORDER BY a.start_time DESC
	`

	rows, err := r.DB.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var appt models.Appointment
		// 2. UPDATED: Add variables to scan the new doctor fields
		var docFirstName, docLastName, docSpecialty string

		// 3. UPDATED: Scan all 10 columns
		err := rows.Scan(
			&appt.ID, &appt.PatientID, &appt.DoctorID, &appt.StartTime, &appt.EndTime, &appt.Status, &appt.Type,
			&docFirstName, &docLastName, &docSpecialty,
		)
		if err != nil {
			return nil, err
		}

		// 4. UPDATED: Populate the new struct fields
		appt.DoctorName = docFirstName + " " + docLastName
		appt.DoctorSpecialty = docSpecialty

		appointments = append(appointments, appt)
	}
	return appointments, nil
}

func (r *Repository) GetPrescriptionsByPatientID(patientID int) ([]models.Prescription, error) {
	query := `
		SELECT p.id, p.patient_id, p.doctor_id, p.medication, p.notes, p.file_name, p.created_at, d.firstName, d.lastName
		FROM prescriptions p
		JOIN doctors d ON p.doctor_id = d.id
		WHERE p.patient_id = $1
		ORDER BY p.created_at DESC
	`
	rows, err := r.DB.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []models.Prescription
	for rows.Next() {
		var pres models.Prescription
		var docFirstName, docLastName string
		err := rows.Scan(&pres.ID, &pres.PatientID, &pres.DoctorID, &pres.Medication, &pres.Notes, &pres.FileName, &pres.CreatedAt, &docFirstName, &docLastName)
		if err != nil {
			return nil, err
		}
		pres.DoctorName = docFirstName + " " + docLastName
		prescriptions = append(prescriptions, pres)
	}
	return prescriptions, nil
}

func (r *Repository) CreateAppointment(patientID int, doctorID int, startTime time.Time, endTime time.Time, apptType string) (int, error) {
	query := `
		INSERT INTO appointments (patient_id, doctor_id, start_time, end_time, appointment_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var newID int
	err := r.DB.QueryRow(query, patientID, doctorID, startTime, endTime, apptType).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}
