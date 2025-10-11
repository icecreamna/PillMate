package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

func AddMyMedicine(db *gorm.DB, mymedicine *models.MyMedicine) (*models.MyMedicine, error) {
	if err := db.Create(mymedicine).Error; err != nil {
		return nil, err
	}
	return mymedicine, nil
}

// READ: ดึงรายการยา "รายการเดียว" ของผู้ป่วย
func GetMyMedicine(db *gorm.DB, patientID, mymedicineID uint) (*models.MyMedicine, error) {
	var m models.MyMedicine
	if err := db.
		Where("id = ? AND patient_id = ?", mymedicineID, patientID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// READ: ดึง "ทั้งหมด" ของผู้ป่วยคนหนึ่ง
func GetMyMedicines(db *gorm.DB, patientID uint) ([]struct {
	ID              uint   `json:"id"`
	MedName         string `json:"med_name"`
	Properties      string `json:"properties"`
	FormID          uint   `json:"form_id"`
	UnitID          uint   `json:"unit_id"`
	InstructionID   uint   `json:"instruction_id"`
	FormName        string `json:"form_name"`
	UnitName        string `json:"unit_name"`
	InstructionName string `json:"instruction_name"`
	AmountPerTime   string `json:"amount_per_time"`
	TimesPerDay     string `json:"times_per_day"`
	Source          string `json:"source"`
}, error) {

	var result []struct {
		ID              uint   `json:"id"`
		MedName         string `json:"med_name"`
		Properties      string `json:"properties"`
		FormID          uint   `json:"form_id"`
		UnitID          uint   `json:"unit_id"`
		InstructionID   uint   `json:"instruction_id"`
		FormName        string `json:"form_name"`
		UnitName        string `json:"unit_name"`
		InstructionName string `json:"instruction_name"`
		AmountPerTime   string `json:"amount_per_time"`
		TimesPerDay     string `json:"times_per_day"`
		Source          string `json:"source"`
	}

	err := db.Table("my_medicines AS m").
		Select(`
			m.id,
			m.med_name,
			m.properties,
			m.form_id,
			m.unit_id,
			m.instruction_id,
			f.form_name,
			u.unit_name,
			i.instruction_name,
			m.amount_per_time,
			m.times_per_day,
			m.source
		`).
		Joins("LEFT JOIN forms f ON m.form_id = f.id").
		Joins("LEFT JOIN units u ON m.unit_id = u.id").
		Joins("LEFT JOIN instructions i ON m.instruction_id = i.id").
		Where("m.patient_id = ? AND m.deleted_at IS NULL", patientID).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateMyMedicine(db *gorm.DB, patientID, mymedicineID uint, in *models.MyMedicine) (*models.MyMedicine, error) {
	// กันไม่ให้ client แอบเปลี่ยน owner/primary key
	in.ID = 0
	in.PatientID = 0

	currentMymedicine := db.Model(&models.MyMedicine{}).
		Where("id = ? AND patient_id = ?", mymedicineID, patientID).
		Updates(in)
	if currentMymedicine.Error != nil {
		return nil, currentMymedicine.Error
	}
	if currentMymedicine.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var out models.MyMedicine
	if err := db.Where("id = ? AND patient_id = ?", mymedicineID, patientID).First(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

// DELETE: ลบรายการยา
// ถ้า Source = "hospital" และมี PrescriptionID -> เซ็ต prescriptions.app_sync_status = false
func DeleteMyMedicine(db *gorm.DB, patientID, mymedicineID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var myMedicine models.MyMedicine

		// 1) โหลดรายการยาที่เป็นของ patient นี้
		if err := tx.Where("id = ? AND patient_id = ?", mymedicineID, patientID).
			First(&myMedicine).Error; err != nil {
			return err // รวมถึง gorm.ErrRecordNotFound
		}

		// 2) ลบ (soft delete ถ้ามี DeletedAt)
		if err := tx.Delete(&myMedicine).Error; err != nil {
			return err
		}

		// 3) ถ้ามาจากโรงพยาบาลและมี PrescriptionID -> rollback สถานะซิงก์
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
