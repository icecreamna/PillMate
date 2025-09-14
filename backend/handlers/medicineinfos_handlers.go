package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// CREATE
func AddMedicineInfo(db *gorm.DB, in *models.MedicineInfo) (*models.MedicineInfo, error) {
	in.ID = 0
	if err := db.Create(in).Error; err != nil {
		return nil, err
	}
	return in, nil
}

// READ ONE
func GetMedicineInfo(db *gorm.DB, id uint) (*models.MedicineInfo, error) {
	var m models.MedicineInfo
	if err := db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// READ ALL
func GetMedicineInfos(db *gorm.DB) ([]models.MedicineInfo, error) {
	var list []models.MedicineInfo
	if err := db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UPDATE: อัปเดตเฉพาะฟิลด์ที่ส่งมา (zero-value จะไม่อัปเดต)
func UpdateMedicineInfo(db *gorm.DB, id uint, in *models.MedicineInfo) (*models.MedicineInfo, error) {
	in.ID = 0
	currentMedicine := db.Model(&models.MedicineInfo{}).Where("id = ?", id).Updates(in)
	if currentMedicine.Error != nil {
		return nil, currentMedicine.Error
	}
	if currentMedicine.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	var out models.MedicineInfo
	if err := db.First(&out, id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

// DELETE (มี DeletedAt → เป็น soft delete อัตโนมัติ)
func DeleteMedicineInfo(db *gorm.DB, id uint) error {
	currentMedicine := db.Delete(&models.MedicineInfo{}, id)
	if currentMedicine.Error != nil {
		return currentMedicine.Error
	}
	if currentMedicine.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}