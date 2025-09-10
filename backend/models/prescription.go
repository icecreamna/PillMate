package models

import (
	
	"time"
)

// ใบสั่งยา
type Prescription struct {
	ID           	uint `gorm:"primaryKey" json:"id"`
	IDCardNumber    string `gorm:"type:char(13)" json:"id_card_number"`
	MedicineInfoID  uint `gorm:"not null" json:"medicine_info_id"`
	AmountPerTime 	string `gorm:"not null" json:"amount_per_time"` //ครั้งละ กี่ หน่วย
	TimesPerDay 	string `gorm:"not null" json:"times_per_day"` //วันละกี่ครั้ง
	HospitalID      uint `gorm:"not null" json:"hospital_id"`
	DoctorID  		uint `gorm:"not null" json:"doctor_id"`
	CreatedAt 		time.Time `json:"created_at"` //วันเวลาที่ออกใบสั่งยา
	AppSyncStatus 	bool `gorm:"default:false" json:"app_sync_status"` // false=ยังไม่ซิงค์

	MedicineInfo 	MedicineInfo `gorm:"foreignKey:MedicineInfoID"`
	Hospital 		Hospital `gorm:"foreignKey:HospitalID"`
	WebAdmin 		WebAdmin `gorm:"foreignKey:DoctorID"`
}