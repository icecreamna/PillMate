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
func GetMyMedicines(db *gorm.DB, patientID uint) ([]models.MyMedicine, error) {
	var list []models.MyMedicine
	if err := db.
		Where("patient_id = ?", patientID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
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