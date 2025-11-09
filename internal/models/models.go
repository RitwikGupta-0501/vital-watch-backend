package models

import "time"

type Authenticatable interface {
	GetID() int
	GetHashedPassword() string
}

type Patient struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

func (p Patient) GetID() int {
	return p.ID
}

func (p Patient) GetHashedPassword() string {
	return p.HashedPassword
}

type Doctor struct {
	ID             int       `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	Specialty      string    `json:"specialty"`
	Experience     int       `json:"experience"`
	Available      bool      `json:"available"`
}

func (d Doctor) GetID() int {
	return d.ID
}
func (d Doctor) GetHashedPassword() string {
	return d.HashedPassword
}

type Appointment struct {
	ID        int       `json:"id"`
	DoctorID  int       `json:"doctor_id"`
	PatientID int       `json:"patient_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
}

type Prescription struct {
	ID         int       `json:"id"`
	PatientID  int       `json:"patient_id"`
	DoctorID   int       `json:"doctor_id"`
	Medication string    `json:"medication"`
	Notes      string    `json:"notes"`
	FileName   string    `json:"file_name"`
	CreatedAt  time.Time `json:"created_at"`

	DoctorName string `json:"doctorName,omitempty"`
}
