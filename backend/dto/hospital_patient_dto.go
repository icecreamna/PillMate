package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

// ====== Request DTOs (no pagination) ======

type CreateHospitalPatientDTO struct {
	IDCardNumber string    `json:"id_card_number"` // 13 digits (validate ใน handler)
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`   // 10 digits (validate ใน handler)
	BirthDay     time.Time `json:"birth_day"`      // ส่งเข้าเป็น date/time; แนะนำรูปแบบ YYYY-MM-DD
	Gender       string    `json:"gender"`         // "ชาย" | "หญิง"
}

type UpdateHospitalPatientDTO struct {
	IDCardNumber *string    `json:"id_card_number,omitempty"`
	FirstName    *string    `json:"first_name,omitempty"`
	LastName     *string    `json:"last_name,omitempty"`
	PhoneNumber  *string    `json:"phone_number,omitempty"`
	BirthDay     *time.Time `json:"birth_day,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
}

// ====== Response DTOs ======
// BirthDay เป็น string รูปแบบ "YYYY-MM-DD"

type HospitalPatientResponse struct {
	ID           uint      `json:"id"`
	IDCardNumber string    `json:"id_card_number"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	BirthDay     string    `json:"birth_day"` // <-- เปลี่ยนเป็น string
	Gender       string    `json:"gender"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewHospitalPatientResponse(m models.HospitalPatient) HospitalPatientResponse {
	return HospitalPatientResponse{
		ID:           m.ID,
		IDCardNumber: m.IDCardNumber,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		PhoneNumber:  m.PhoneNumber,
		BirthDay:     dateToYMD(m.BirthDay),
		Gender:       m.Gender,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// ====== helpers ======

func dateToYMD(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		loc = time.Local
	}
	return t.In(loc).Format("2006-01-02")
}
