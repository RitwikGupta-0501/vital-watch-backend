package models

import "time"

type Patient struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"` // Always hide this
	CreatedAt time.Time `json:"created_at"`
}

type Doctor struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"` // Always hide this
	CreatedAt time.Time `json:"created_at"`
}

type Appointment struct {
	ID        string    `json:"id"`
	DoctorID  string    `json:"doctor_id"`
	PatientID string    `json:"patient_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"` // "scheduled", "completed", "cancelled"
}

// ... add Prescription, etc. later
