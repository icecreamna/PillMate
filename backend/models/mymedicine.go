package models

import(
	// "gorm.io/gorm"
	// "time"
)

// ยาของฉัน
type MyMedicine struct {
	ID           	uint `gorm:"primaryKey" json:"id"`
	PatientID 		uint `gorm:"not null" json:"patient_id"`
	MedName 	 	string `gorm:"type:varchar(255);not null" json:"med_name"` // ชื่อยา
	Properties 	 	string `gorm:"not null" json:"properties"` //รายละเอียด สรรพคุณ
	Quantity		int `gorm:"not null;default:1" json:"quantity"` // ปริมาณ
	FormID 			uint `gorm:"not null" json:"form_id"`//รูปแบบยา
    UnitID 			uint `gorm:"default:null" json:"unit_id"` // หน่วยยา
	InstructionID 	uint `gorm:"default:null" json:"instruction_id"` // คำแนะนำการทานยา
	Source 		 	string `gorm:"check:source IN ('manual','hospital')" json:"source"`

	Patient   		Patient `gorm:"foreignKey:PatientID"`
	Form 			Form `gorm:"foreignKey:FormID"`
    Unit 			Unit `gorm:"foreignKey:UnitID"`
	Instruction 	Instruction `gorm:"foreignKey:InstructionID"`

	// Many-to-Many
	DosageTimes []DosageTime `gorm:"many2many:my_medicine_dosage_times;" json:"dosage_times"`
}