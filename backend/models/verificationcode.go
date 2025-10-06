package models

import (
	"time"
	"gorm.io/gorm"
)

type VerificationCode struct {
	ID        uint           `gorm:"primaryKey" json:"id"`

	// เก็บเป็น plaintext เหมือนเดิม แต่ "ไม่ส่งออก" ทาง JSON
	OTPCode   string         `gorm:"type:char(6);not null" json:"-"`

	// เวลาหมดอายุ + เวลาที่ใช้แล้ว (optional)
	ExpiresAt time.Time      `json:"expires_at"`
	UsedAt    *time.Time     `json:"used_at,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	PatientID uint           `gorm:"index;not null" json:"patient_id"`
	Patient   Patient        `gorm:"foreignKey:PatientID" json:"-"` // ซ่อน relation
}
