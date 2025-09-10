package models

import(
	"gorm.io/gorm"
	"time"
)

// บันทึกอาการ
type Symptom struct {
    ID   			uint `gorm:"primaryKey" json:"id"`
	PatientID       uint `gorm:"not null" json:"patient_id"`
	MyMedicineID    uint `gorm:"not null" json:"my_medicine_id"`
	GroupMedicineID uint `gorm:"default:null" json:"group_medicine_id"`
	NotiItemID 		uint `gorm:"not null" json:"noti_item_id"`
	SymptomName		string `gorm:"not null" json:"symptom_name"`
	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `gorm:"index" json:"-"`

	Patient 		Patient `gorm:"foreignKey:PatientID"`
	MyMedicine 		MyMedicine `gorm:"foreignKey:MyMedicineID"`
	GroupMedicine 	GroupMedicine `gorm:"foreignKey:GroupMedicineID"`
	NotiItem 		NotiItem `gorm:"foreignKey:NotiItemID"`
}