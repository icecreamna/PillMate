package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ส่งออกสำหรับ Mobile (สตริงทั้งหมดสำหรับ date/time)
type MobileAppointmentResponse struct {
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

func ToMobileAppointmentResponse(m models.Appointment) MobileAppointmentResponse {
	return MobileAppointmentResponse{
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
