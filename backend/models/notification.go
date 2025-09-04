package models

import(

	"time"
)

// เวลาแจ้งเตือน
type NotiTime struct {
    ID   			uint `gorm:"primaryKey" json:"id"`
	MyMedicineID    uint `gorm:"not null" json:"my_medicine_id"`
	DosageTimeID 	uint `gorm:"not null" json:"dosage_time_id"`
	NotifyTime      time.Time `gorm:"type:time" json:"notify_time"`

	MyMedicine 		MyMedicine `gorm:"foreignKey:MyMedicineID"`
	DosageTime 		DosageTime `gorm:"foreignKey:DosageTimeID"`
}

// ข้อมูลแจ้งเตือน
type NotiInfo struct {
    ID   			uint `gorm:"primaryKey" json:"id"`
	MyMedicineID    uint  `gorm:"not null" json:"my_medicine_id"`
	StartDate 		time.Time `gorm:"type:date" json:"start_date"`
	EndDate 		time.Time `gorm:"type:date" json:"end_date"`

	MyMedicine 		MyMedicine `gorm:"foreignKey:MyMedicineID"`
}

// รายการแจ้งเตือน
type NotiItem struct {
    ID   			uint `gorm:"primaryKey" json:"id"`
	MyMedicineID    uint  `gorm:"not null" json:"my_medicine_id"`
	MedName 	 	string `gorm:"type:varchar(255);not null" json:"med_name"` // ชื่อยา
	Quantity		int `gorm:"not null;default:1" json:"quantity"` // ปริมาณ
	FormID 			uint `gorm:"not null" json:"form_id"`//รูปแบบยา
    UnitID 			uint `gorm:"default:null" json:"unit_id"` // หน่ยยา
	InstructionID 	uint `gorm:"default:null" json:"instruction_id"` // คำแนะนำการทานยา
	NotifyTime      time.Time `gorm:"type:time" json:"notify_time"` //เวลาที่แจ้งเตือน
	NotifyDate      time.Time `gorm:"type:date" json:"notify_date"` //วันที่แจ้งเตือน
	Source 		 	string `gorm:"check:source IN ('manual','hospital')" json:"source"` //แหล่งที่มา
	TakenStatus		bool `gorm:"default:false" json:"taken_status"` //สถานะการกินยา
	TakenTimeAt		time.Time `gorm:"autoCreateTime" json:"taken_time_at"` //วันเวลาที่เปลี่ยนสถานะ
	NotifyStatus	bool `gorm:"default:false" json:"notify_status"` //สถานะการแจ้งเตือน
	CreatedAt 		time.Time `json:"created_at"` //วันเวลาที่สร้างรายการแจ้งเตือน

	MyMedicine 		MyMedicine `gorm:"foreignKey:MyMedicineID"`
	Form 			Form `gorm:"foreignKey:FormID"`
    Unit 			Unit `gorm:"foreignKey:UnitID"`
	Instruction 	Instruction `gorm:"foreignKey:InstructionID"`
}