package dto

import (
	"time"
	
	"github.com/fouradithep/pillmate/models"
)

// ส่งออกสำหรับ Mobile (สตริงทั้งหมดสำหรับ date/time)
type MobileAppointmentResponse struct {
	ID              uint    `json:"id"`
	IDCardNumber    string  `json:"id_card_number"`
	AppointmentDate string  `json:"appointment_date"` // "YYYY-MM-DD" (Bangkok)
	AppointmentTime string  `json:"appointment_time"` // "HH:MM"     (Bangkok)
	HospitalID      uint    `json:"hospital_id"`
	DoctorID        uint    `json:"doctor_id"`
	Note            *string `json:"note,omitempty"`
	CreatedAt       string  `json:"created_at"`       // RFC3339 (Bangkok)
	UpdatedAt       string  `json:"updated_at"`       // RFC3339 (Bangkok)
}

func ToMobileAppointmentResponse(m models.Appointment) MobileAppointmentResponse {
	return MobileAppointmentResponse{
		ID:              m.ID,
		IDCardNumber:    m.IDCardNumber,
		AppointmentDate: formatYMDBangkok(m.AppointmentDate), // UTC -> Bangkok -> "YYYY-MM-DD"
		AppointmentTime: formatHMBangkok(m.AppointmentTime),  // UTC -> Bangkok -> "HH:MM"
		HospitalID:      m.HospitalID,
		DoctorID:        m.DoctorID,
		Note:            m.Note,
		CreatedAt:       m.CreatedAt.In(appLoc).Format(time.RFC3339), // "YYYY-MM-DDTHH:MM:SS+07:00"
		UpdatedAt:       m.UpdatedAt.In(appLoc).Format(time.RFC3339),
	}
}

func ToMobileAppointmentResponses(ms []models.Appointment) []MobileAppointmentResponse {
	out := make([]MobileAppointmentResponse, 0, len(ms))
	for _, m := range ms {
		out = append(out, ToMobileAppointmentResponse(m))
	}
	return out
}
