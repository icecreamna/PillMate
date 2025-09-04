package models

import (
	
	"time"
)

// ใบสั่งยา
type Prescription struct {
	ID           	uint `gorm:"primaryKey" json:"id"`
	IDCardNumber    string `gorm:"type:char(13)" json:"id_card_number"`
	MedicineInfoID  uint `gorm:"not null" json:"medicine_info_id"`
	Quantity		int `gorm:"not null;default:1" json:"quantity"` // ปริมาณ
	HospitalID      uint `gorm:"not null" json:"hospital_id"`
	CreatedAt 		time.Time `json:"created_at"` //วันเวลาที่ออกใบสั่งยา
	AppSyncStatus 	bool `gorm:"default:false" json:"app_sync_status"` // false=ยังไม่ซิงค์

	// Patient  		Patient `gorm:"foreignKey:IDCardNumber;references:IDCardNumber"`
	MedicineInfo 	MedicineInfo `gorm:"foreignKey:MedicineInfoID"`
	Hospital 		Hospital `gorm:"foreignKey:HospitalID"`
}