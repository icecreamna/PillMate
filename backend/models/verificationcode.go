package models

import(
	"gorm.io/gorm"
	"time"
)
type VerificationCode struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	OTPCode   string         `gorm:"type:char(6);not null" json:"otp_code"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created"`
    ExpiresAt time.Time      `json:"expires"`
	UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	PatientID uint           `gorm:"not null" json:"patient_id"`
    Patient   Patient        `gorm:"foreignKey:PatientID"`
}