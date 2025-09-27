package models

import(
	"gorm.io/gorm"
	"time"
)

// กลุ่ม 
type Group struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	PatientID uint           `gorm:"not null;index" json:"patient_id"`
	GroupName string         `gorm:"type:varchar(255);not null" json:"group_name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
// ยาของฉัน
type MyMedicine struct {
	ID           	uint `gorm:"primaryKey" json:"id"`
	PatientID 		uint `gorm:"not null" json:"patient_id"`
	MedName 	 	string `gorm:"type:varchar(255);not null" json:"med_name"` // ชื่อยา
	Properties 	 	string `gorm:"not null" json:"properties"` //รายละเอียด สรรพคุณ
	FormID 			uint `gorm:"not null" json:"form_id"`//รูปแบบยา
    UnitID 			uint `gorm:"default:null" json:"unit_id"` // หน่วยยา
	InstructionID 	uint `gorm:"default:null" json:"instruction_id"` // ช่วงเวลาใช้ยา
	AmountPerTime 	string `gorm:"not null" json:"amount_per_time"` //ครั้งละ กี่ หน่วย
	TimesPerDay 	string `gorm:"not null" json:"times_per_day"` //วันละกี่ครั้ง
	Source 		 	string `gorm:"check:source IN ('manual','hospital')" json:"source"`
	PrescriptionID  *uint `gorm:"index" json:"prescription_id,omitempty"`
	GroupID 		*uint `gorm:"index" json:"group_id,omitempty"`

	Group   		Group `gorm:"foreignKey:GroupID"`
	Patient   		Patient `gorm:"foreignKey:PatientID"`
	Form 			Form `gorm:"foreignKey:FormID"`
    Unit 			Unit `gorm:"foreignKey:UnitID"`
	Instruction 	Instruction `gorm:"foreignKey:InstructionID"`
	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `gorm:"index" json:"-"`

}
