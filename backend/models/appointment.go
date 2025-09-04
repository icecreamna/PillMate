package models

import (
	// "gorm.io/gorm"
	"time"
)

type Hospital struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	HospitalName string `gorm:"type:varchar(100);not null;unique" json:"hospital_name"`
}

// ข้อมูลนัดพบแพทย์
type Appointment struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	IDCardNumber    string    `gorm:"type:char(13)" json:"id_card_number"`
	AppointmentDate time.Time `gorm:"type:date" json:"appointment_date"`
	AppointmentTime time.Time `gorm:"type:time" json:"appointment_time"`
	HospitalID      uint      `gorm:"not null" json:"hospital_id"`

	// Patient  Patient  `gorm:"foreignKey:IDCardNumber;references:IDCardNumber"`
	Hospital Hospital `gorm:"foreignKey:HospitalID"`

	AppSyncStatus bool `gorm:"default:false" json:"app_sync_status"` // false=ยังไม่ซิงค์
}

// ข้อมูลแจ้งเตือนนัดพบแพทย์
type AppointmentNoti struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	PatientID          uint      `gorm:"not null" json:"patient_id"`
	NotifyTime         time.Time `gorm:"type:time" json:"notify_time"`
	ReminderDateBefore time.Time `gorm:"type:date" json:"reminder_date_before"`
	ReminderDateOn     time.Time `gorm:"type:date" json:"reminder_date_on"`
	StatusBefore       bool      `gorm:"default:false" json:"status_before"`
	StatusOn           bool      `gorm:"default:false" json:"status_on"`

	Patient 			Patient `gorm:"foreignKey:PatientID"`
}

// บันทึกการแจ้งเตือน
type NotiLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PatientID  uint      `gorm:"not null" json:"patient_id"`
	Message    string    `gorm:"not null" json:"message"`
	NotifiedAt time.Time `gorm:"autoCreateTime" json:"notified_at"`

	Patient    Patient `gorm:"foreignKey:PatientID"`
}
