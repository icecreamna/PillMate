package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ===== Fixed timezone (Bangkok) =====
var appLoc = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Bangkok")
	return l
}()

func formatYMDBangkok(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(appLoc).Format("2006-01-02")
}

func formatHMBangkok(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(appLoc).Format("15:04")
}

// ====== Request DTOs (รับเป็นสตริง) ======

type DoctorCreateAppointmentDTO struct {
	IDCardNumber    string  `json:"id_card_number"`        // 13 digits
	AppointmentDate string  `json:"appointment_date"`      // "YYYY-MM-DD"
	AppointmentTime string  `json:"appointment_time"`      // "HH:MM"
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
		AppointmentDate: formatYMDBangkok(m.AppointmentDate), // UTC -> Bangkok -> "YYYY-MM-DD"
		AppointmentTime: formatHMBangkok(m.AppointmentTime),  // UTC -> Bangkok -> "HH:MM"
		HospitalID:      m.HospitalID,
		DoctorID:        m.DoctorID,
		Note:            m.Note,
		// แสดงเป็นเวลาไทยเพื่อความสอดคล้อง
		CreatedAt: m.CreatedAt.In(appLoc),
		UpdatedAt: m.UpdatedAt.In(appLoc),
	}
}

func ToDoctorAppointmentResponses(ms []models.Appointment) []DoctorAppointmentResponse {
	out := make([]DoctorAppointmentResponse, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToDoctorAppointmentResponse(m))
	}
	return out
}
