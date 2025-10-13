package models

import (
	"gorm.io/gorm"
	"time"
)

// บันทึกอาการ
type Symptom struct {
	ID           uint  `gorm:"primaryKey" json:"id"`
	PatientID    uint  `gorm:"not null;index" json:"patient_id"`
	MyMedicineID *uint `json:"my_medicine_id,omitempty"`
	GroupID      *uint `json:"group_id,omitempty"` // ทำให้เป็น nullable
	NotiItemID   uint  `gorm:"not null;index" json:"noti_item_id"`

	SymptomNote string `gorm:"type:text;not null" json:"symptom_note"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Patient    Patient    `gorm:"foreignKey:PatientID"`
	MyMedicine MyMedicine `gorm:"foreignKey:MyMedicineID"`
	Group      Group      `gorm:"foreignKey:GroupID"`
	NotiItem   NotiItem   `gorm:"foreignKey:NotiItemID"`
}
