package models

import (
	"gorm.io/gorm"
	"time"
)

type VerificationCode struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PatientID uint           `gorm:"index;constraint:OnDelete:CASCADE;" json:"patient_id"` 
	OTPCode   string         `json:"otp_code"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
