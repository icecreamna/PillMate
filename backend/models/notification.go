package models

import(

	"time"
	"github.com/lib/pq"
	
)

// รูปแบบการแจ้งเตือน
type NotiFormat struct {
    ID   	     uint `gorm:"primaryKey" json:"id"`
	FormatName   string `gorm:"type:varchar(255);not null;unique" json:"format_name"` //ชื่อรูปแบบ
	
}

// ข้อมูลแจ้งเตือน
type NotiInfo struct {
    ID   				uint `gorm:"primaryKey" json:"id"`
	MyMedicineID    	uint `gorm:"default:null" json:"my_medicine_id"`
	GroupMedicineID 	uint `gorm:"default:null" json:"group_medicine_id"`
	StartDate 			time.Time `gorm:"type:date" json:"start_date"`
	EndDate 			time.Time `gorm:"type:date" json:"end_date"`
	NotiFormatID		uint `gorm:"not null" json:"noti_format_id"`
	Times        		pq.StringArray  `gorm:"type:text[]" json:"times"` // ["08:00","12:00","20:00"]	
	IntervalHours		int `gorm:"default:null" json:"interval_hours"` // แจ้งเตือนทุกกี่ชั่วโมง
	TimesPerDay 	    int `gorm:"default:null" json:"times_per_day"` //กี่ครั้งต่อวัน
	IntervalDay			int `gorm:"default:null" json:"interval_day"` // แจ้งเตือนทุกกี่วัน
	DaysOfWeek 			pq.Int64Array `gorm:"type:int[]" json:"days_of_week"`// วันที่จะแจ้งเตือน อาทิตย์=0 … เสาร์=6
	CyclePattern  		pq.Int64Array `gorm:"type:int[];default:null" json:"cycle_pattern"` // ตัวอย่าง: [21,7]

	MyMedicine 			MyMedicine `gorm:"foreignKey:MyMedicineID"`
	GroupMedicine 		GroupMedicine `gorm:"foreignKey:GroupMedicineID"`
	NotiFormat			NotiFormat `gorm:"foreignKey:NotiFormatID"`
}

// รายการแจ้งเตือน
type NotiItem struct {
    ID   			uint `gorm:"primaryKey" json:"id"`
	PatientID 		uint `gorm:"not null" json:"patient_id"`
	MyMedicineID    uint  `gorm:"not null" json:"my_medicine_id"`
	GroupMedicineID uint `gorm:"default:null" json:"group_medicine_id"`
	NotiInfoID 		uint `gorm:"default:null" json:"noti_info_id"`
	MedName 	 	string `gorm:"type:varchar(255);not null" json:"med_name"` // ชื่อยา
	GroupName 	 	string `gorm:"type:varchar(255);not null" json:"group_name"` // ชื่อกลุ่ม
	Quantity		int `gorm:"not null;default:1" json:"quantity"` // ปริมาณ
	FormID 			uint `gorm:"not null" json:"form_id"`//รูปแบบยา
    UnitID 			uint `gorm:"default:null" json:"unit_id"` // หน่ยยา
	InstructionID 	uint `gorm:"default:null" json:"instruction_id"` // ช่วงเวลาใช้ยา
	NotifyTime      time.Time `gorm:"type:time" json:"notify_time"` //เวลาที่แจ้งเตือน
	NotifyDate      time.Time `gorm:"type:date" json:"notify_date"` //วันที่แจ้งเตือน
	TakenStatus		bool `gorm:"default:false" json:"taken_status"` //สถานะการกินยา
	TakenTimeAt		time.Time `gorm:"autoCreateTime" json:"taken_time_at"` //วันเวลาที่เปลี่ยนสถานะ
	NotifyStatus	bool `gorm:"default:false" json:"notify_status"` //สถานะการแจ้งเตือน default:false = ยังไม่แจ้งเตือน
	CreatedAt 		time.Time `json:"created_at"` //วันเวลาที่สร้างรายการแจ้งเตือน

	Patient   		Patient `gorm:"foreignKey:PatientID"`
	MyMedicine 		MyMedicine `gorm:"foreignKey:MyMedicineID"`
	GroupMedicine 	GroupMedicine `gorm:"foreignKey:GroupMedicineID"`
	NotiInfo 		NotiInfo `gorm:"foreignKey:NotiInfoID"`
	Form 			Form `gorm:"foreignKey:FormID"`
    Unit 			Unit `gorm:"foreignKey:UnitID"`
	Instruction 	Instruction `gorm:"foreignKey:InstructionID"`
}