package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

// =========================
// Forms
// =========================

// GetForm — อ่าน Form หนึ่งรายการ เลือก preload Units ได้
func GetForm(db *gorm.DB, formID uint, includeRelations bool) (*models.Form, error) {
	var form models.Form

	query := db
	if includeRelations {
		query = query.Preload("Units")
	}

	if err := query.First(&form, formID).Error; err != nil {
		return nil, err
	}
	return &form, nil
}

// GetForms — อ่าน Form ทั้งหมด เลือก preload Units ได้
func GetForms(db *gorm.DB, includeRelations bool) ([]models.Form, error) {
	var forms []models.Form

	query := db
	if includeRelations {
		query = query.Preload("Units")
	}

	if err := query.Find(&forms).Error; err != nil {
		return nil, err
	}
	return forms, nil
}

// GetUnitsByFormID — ดึง Units ที่ผูกกับ Form หนึ่ง ๆ
func GetUnitsByFormID(db *gorm.DB, formID uint) ([]models.Unit, error) {
	// เบาและตรงจุดด้วย Association (ถ้าต้องการเฉพาะบางคอลัมน์/จัดเรียง แนะนำใช้ JOIN + Scan)
	var form models.Form
	if err := db.Select("id").First(&form, formID).Error; err != nil {
		return nil, err
	}

	var units []models.Unit
	if err := db.Model(&form).Association("Units").Find(&units); err != nil {
		return nil, err
	}
	return units, nil
}

// =========================
// Units
// =========================

// GetUnit — อ่าน Unit หนึ่งรายการ เลือก preload Forms ได้
func GetUnit(db *gorm.DB, unitID uint, includeRelations bool) (*models.Unit, error) {
	var unit models.Unit

	query := db
	if includeRelations {
		query = query.Preload("Forms")
	}

	if err := query.First(&unit, unitID).Error; err != nil {
		return nil, err
	}
	return &unit, nil
}

// GetUnits — อ่าน Unit ทั้งหมด เลือก preload Forms ได้
func GetUnits(db *gorm.DB, includeRelations bool) ([]models.Unit, error) {
	var units []models.Unit

	query := db
	if includeRelations {
		query = query.Preload("Forms")
	}

	if err := query.Find(&units).Error; err != nil {
		return nil, err
	}
	return units, nil
}

// GetFormsByUnitID — ดึง Forms ที่ผูกกับ Unit หนึ่ง ๆ
func GetFormsByUnitID(db *gorm.DB, unitID uint) ([]models.Form, error) {
	var unit models.Unit
	if err := db.Select("id").First(&unit, unitID).Error; err != nil {
		return nil, err
	}

	var forms []models.Form
	if err := db.Model(&unit).Association("Forms").Find(&forms); err != nil {
		return nil, err
	}
	return forms, nil
}

// =========================
// FormUnit (pivot)
// =========================

// GetFormUnits — ดึง mapping ทั้งหมดของตาราง form_units
func GetFormUnits(db *gorm.DB) ([]models.FormUnit, error) {
	var formUnits []models.FormUnit
	if err := db.Find(&formUnits).Error; err != nil {
		return nil, err
	}
	return formUnits, nil
}

// GetFormUnitsByFormID — ดึง mapping ตาม FormID
func GetFormUnitsByFormID(db *gorm.DB, formID uint) ([]models.FormUnit, error) {
	var formUnits []models.FormUnit
	if err := db.Where("form_id = ?", formID).Find(&formUnits).Error; err != nil {
		return nil, err
	}
	return formUnits, nil
}

// GetFormUnitsByUnitID — ดึง mapping ตาม UnitID
func GetFormUnitsByUnitID(db *gorm.DB, unitID uint) ([]models.FormUnit, error) {
	var formUnits []models.FormUnit
	if err := db.Where("unit_id = ?", unitID).Find(&formUnits).Error; err != nil {
		return nil, err
	}
	return formUnits, nil
}

// =========================
// Instructions
// =========================

// GetInstruction — อ่าน Instruction หนึ่งรายการ
func GetInstruction(db *gorm.DB, instructionID uint) (*models.Instruction, error) {
	var instruction models.Instruction
	if err := db.First(&instruction, instructionID).Error; err != nil {
		return nil, err
	}
	return &instruction, nil
}

// GetInstructions — อ่าน Instruction ทั้งหมด
func GetInstructions(db *gorm.DB) ([]models.Instruction, error) {
	var instructions []models.Instruction
	if err := db.Find(&instructions).Error; err != nil {
		return nil, err
	}
	return instructions, nil
}
