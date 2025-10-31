package handlers

import (
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

/** ========= Views / DTOs สำหรับ Scan & JSON ========= **/

// ใช้ struct แบบ "มีชื่อ" เพื่อหลีกเลี่ยงชนิดไม่ตรงกันเวลาคืนค่า
type MyMedicineView struct {
	ID              uint    `json:"id"`
	MedName         string  `json:"med_name"`
	Properties      string  `json:"properties"`
	FormID          uint    `json:"form_id"`
	UnitID          *uint   `json:"unit_id,omitempty"`
	InstructionID   *uint   `json:"instruction_id,omitempty"`
	GroupID         *uint   `json:"group_id,omitempty"`
	FormName        string  `json:"form_name"`
	UnitName        *string `json:"unit_name,omitempty"`
	InstructionName *string `json:"instruction_name,omitempty"`
	AmountPerTime   string  `json:"amount_per_time"`
	TimesPerDay     string  `json:"times_per_day"`
	Source          string  `json:"source"`

	// จาก prescription_items (อาจเป็น nil หาก source=manual)
	StartDate  *string `json:"start_date,omitempty"`
	EndDate    *string `json:"end_date,omitempty"`
	ExpireDate *string `json:"expire_date,omitempty"`
	Note       *string `json:"note,omitempty"`
}

/** ========= Helpers ========= **/

// เที่ยงคืนของ "วันนี้" ตาม Asia/Bangkok (ใช้กรองยังไม่หมดอายุให้เสมอ)
func todayBangkokMidnight() time.Time {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	now := time.Now().In(loc)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}

/** ===================== CREATE ===================== **/

// เพิ่มรายการยา
func AddMyMedicine(db *gorm.DB, mymedicine *models.MyMedicine) (*models.MyMedicine, error) {
	if err := db.Create(mymedicine).Error; err != nil {
		return nil, err
	}
	return mymedicine, nil
}

/** ====================== READ ====================== **/

// ดึง "โมเดลดิบ" รายการเดียว (ไม่ JOIN) — แต่ยังคงกรอง "ยังไม่หมดอายุ"
func GetMyMedicine(db *gorm.DB, patientID, mymedicineID uint) (*models.MyMedicine, error) {
	var m models.MyMedicine
	today := todayBangkokMidnight()

	if err := db.
		Where("id = ? AND patient_id = ?", mymedicineID, patientID).
		Where("deleted_at IS NULL").
		Where("(expire_at IS NULL OR expire_at >= ?)", today).
		First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// ดึง "รายการเดียว" แบบ JOIN เอาวันที่/โน้ตจาก prescription_items (และกรองยังไม่หมดอายุ)
func GetMyMedicineWithDates(db *gorm.DB, patientID, mymedicineID uint) (*MyMedicineView, error) {
	var row MyMedicineView
	today := todayBangkokMidnight()

	if err := db.Table("my_medicines AS m").
		Select(`
			m.id,
			m.med_name,
			m.properties,
			m.form_id,
			m.unit_id,
			m.instruction_id,
			m.group_id,
			f.form_name,
			u.unit_name,
			i.instruction_name,
			m.amount_per_time,
			m.times_per_day,
			m.source,
			pi.start_date   AS start_date,
			pi.end_date     AS end_date,
			pi.expire_date  AS expire_date,
			pi.note         AS note
		`).
		Joins("LEFT JOIN forms f ON f.id = m.form_id").
		Joins("LEFT JOIN units u ON u.id = m.unit_id").
		Joins("LEFT JOIN instructions i ON i.id = m.instruction_id").
		Joins("LEFT JOIN prescription_items pi ON pi.id = m.prescription_item_id").
		Where("m.id = ? AND m.patient_id = ? AND m.deleted_at IS NULL", mymedicineID, patientID).
		Where("(m.expire_at IS NULL OR m.expire_at >= ?)", today).
		Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ดึง "ทั้งหมด" ของผู้ป่วย (JOIN prescription_items) — กรองยังไม่หมดอายุ
func GetMyMedicines(db *gorm.DB, patientID uint) ([]MyMedicineView, error) {
	var result []MyMedicineView
	today := todayBangkokMidnight()

	err := db.Table("my_medicines AS m").
		Select(`
			m.id,
			m.med_name,
			m.properties,
			m.form_id,
			m.unit_id,
			m.instruction_id,
			m.group_id,
			f.form_name,
			u.unit_name,
			i.instruction_name,
			m.amount_per_time,
			m.times_per_day,
			m.source,
			pi.start_date   AS start_date,
			pi.end_date     AS end_date,
			pi.expire_date  AS expire_date,
			pi.note         AS note
		`).
		Joins("LEFT JOIN forms f ON f.id = m.form_id").
		Joins("LEFT JOIN units u ON u.id = m.unit_id").
		Joins("LEFT JOIN instructions i ON i.id = m.instruction_id").
		Joins("LEFT JOIN prescription_items pi ON pi.id = m.prescription_item_id").
		Where("m.patient_id = ? AND m.deleted_at IS NULL", patientID).
		Where("(m.expire_at IS NULL OR m.expire_at >= ?)", today).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

/** ===================== UPDATE ===================== **/

// ป้องกันการแก้ owner/PK และ "ย้ายลิงก์ใบสั่งยา" โดยไม่ได้ตั้งใจ
func UpdateMyMedicine(db *gorm.DB, patientID, mymedicineID uint, in *models.MyMedicine) (*models.MyMedicine, error) {
	// กันเปลี่ยน owner/PK
	in.ID = 0
	in.PatientID = 0

	// กันย้ายลิงก์ prescription โดยไม่ตั้งใจ (ลบสองบรรทัดนี้หากต้องการอนุญาต)
	in.PrescriptionID = nil
	in.PrescriptionItemID = nil

	tx := db.Model(&models.MyMedicine{}).
		Where("id = ? AND patient_id = ?", mymedicineID, patientID).
		Omit("prescription_id", "prescription_item_id").
		Updates(in)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var out models.MyMedicine
	today := todayBangkokMidnight()
	if err := db.Where("id = ? AND patient_id = ?", mymedicineID, patientID).
		Where("deleted_at IS NULL").
		Where("(expire_at IS NULL OR expire_at >= ?)", today).
		First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

/** ===================== DELETE ===================== **/

// ถ้า Source = "hospital" และมี PrescriptionID -> เซ็ต prescriptions.app_sync_status = false
func DeleteMyMedicine(db *gorm.DB, patientID, mymedicineID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var myMedicine models.MyMedicine

		// 1) ตรวจว่าเป็นของ patient นี้จริง (จะลบได้แม้หมดอายุแล้ว)
		if err := tx.Where("id = ? AND patient_id = ?", mymedicineID, patientID).
			First(&myMedicine).Error; err != nil {
			return err // รวม gorm.ErrRecordNotFound
		}

		// 2) ลบ (soft delete)
		if err := tx.Delete(&myMedicine).Error; err != nil {
			return err
		}

		// 3) โรลแบ็กสถานะซิงก์ของใบสั่งยา ถ้ามาจากโรงพยาบาล
		if myMedicine.Source == "hospital" && myMedicine.PrescriptionID != nil {
			if err := tx.Model(&models.Prescription{}).
				Where("id = ?", *myMedicine.PrescriptionID).
				Update("app_sync_status", false).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
