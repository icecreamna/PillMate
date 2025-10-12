package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ====== Request DTOs (รับเป็นสตริง) ======

type DoctorCreateAppointmentDTO struct {
	IDCardNumber    string  `json:"id_card_number"`   // 13 digits
	AppointmentDate string  `json:"appointment_date"` // "YYYY-MM-DD"
	AppointmentTime string  `json:"appointment_time"` // "HH:MM"
	HospitalID      *uint   `json:"hospital_id,omitempty"` // optional (default = 1)
	DoctorID        *uint   `json:"doctor_id,omitempty"`   // optional (default = from token)
	Note            *string `json:"note,omitempty"`
}

type DoctorUpdateAppointmentDTO struct {
	IDCardNumber    *string `json:"id_card_number,omitempty"`
	AppointmentDate *string `json:"appointment_date,omitempty"` // "YYYY-MM-DD"
	AppointmentTime *string `json:"appointment_time,omitempty"` // "HH:MM"
	HospitalID      *uint   `json:"hospital_id,omitempty"`
	DoctorID        *uint   `json:"doctor_id,omitempty"`
	Note            *string `json:"note,omitempty"`
}

// ====== Response DTOs (ส่งออกเป็นสตริง) ======

type DoctorAppointmentResponse struct {
	ID              uint      `json:"id"`
	IDCardNumber    string    `json:"id_card_number"`
	AppointmentDate string    `json:"appointment_date"` // "YYYY-MM-DD"
	AppointmentTime string    `json:"appointment_time"` // "HH:MM"
	HospitalID      uint      `json:"hospital_id"`
	DoctorID        uint      `json:"doctor_id"`
	Note            *string   `json:"note,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func ToDoctorAppointmentResponse(m models.Appointment) DoctorAppointmentResponse {
	return DoctorAppointmentResponse{
		ID:              m.ID,
		IDCardNumber:    m.IDCardNumber,
		AppointmentDate: m.AppointmentDate.In(time.Local).Format("2006-01-02"),
		AppointmentTime: m.AppointmentTime.In(time.Local).Format("15:04"),
		HospitalID:      m.HospitalID,
		DoctorID:        m.DoctorID,
		Note:            m.Note,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func ToDoctorAppointmentResponses(ms []models.Appointment) []DoctorAppointmentResponse {
	out := make([]DoctorAppointmentResponse, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToDoctorAppointmentResponse(m))
	}
	return out
}
