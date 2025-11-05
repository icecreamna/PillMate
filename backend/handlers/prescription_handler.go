package handlers

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/models"
)

// ===== CREATE (หัว + items) =====
func CreatePrescription(db *gorm.DB, in *dto.CreatePrescriptionDTO) (*dto.PrescriptionResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	idc := OnlyDigits(in.IDCardNumber)
	if len(idc) != 13 {
		return nil, errors.New("id_card_number must be 13 digits")
	}
	docID := in.DoctorID
	if docID == 0 {
		return nil, errors.New("doctor_id is required")
	}
	if len(in.Items) == 0 {
		return nil, errors.New("items must be non-empty")
	}

	// ตรวจ doctor
	{
		var doc models.WebAdmin
		if err := db.Where("role = ?", "doctor").First(&doc, docID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("doctor not found")
			}
			return nil, err
		}
	}

	var created models.Prescription

	if err := db.Transaction(func(tx *gorm.DB) error {
		// 1) สร้างหัว
		created = models.Prescription{
			IDCardNumber:  idc,
			DoctorID:      docID,
			AppSyncStatus: false,
		}
		if in.SyncUntil != nil && !in.SyncUntil.IsZero() {
			created.SyncUntil = *in.SyncUntil
		}
		if in.AppSyncStatus != nil {
			created.AppSyncStatus = *in.AppSyncStatus
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}

		// 2) สร้าง items
		items := make([]models.PrescriptionItem, 0, len(in.Items))
		for i, it := range in.Items {
			if it.MedicineInfoID == 0 {
				return fmt.Errorf("items[%d].medicine_info_id is required", i)
			}
			amt := Norm(it.AmountPerTime)
			tpd := Norm(it.TimesPerDay)
			if amt == "" {
				return fmt.Errorf("items[%d].amount_per_time required", i)
			}
			if tpd == "" {
				return fmt.Errorf("items[%d].times_per_day required", i)
			}

			// FK: medicine info
			var mi models.MedicineInfo
			if err := tx.First(&mi, it.MedicineInfoID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("items[%d] medicine_info not found", i)
				}
				return err
			}

			items = append(items, models.PrescriptionItem{
				PrescriptionID: created.ID,
				MedicineInfoID: it.MedicineInfoID,
				AmountPerTime:  amt,
				TimesPerDay:    tpd,
				StartDate:      it.StartDate, // NEW
  				EndDate:        it.EndDate,   // NEW (hook จะตั้ง ExpireDate ให้เอง)
  				Note:           it.Note,      // NEW
			})
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		// // 3) **สำคัญ**: ปิดซิงก์ของใบเก่าทั้งหมดให้เป็น "เมื่อวาน"
		// //    ตามนโยบาย: เมื่อมีใบใหม่เข้ามา ให้ใบเก่าหมดอายุการซิงก์ทันที
		// yesterday := time.Now().AddDate(0, 0, -1)
		// if err := tx.Model(&models.Prescription{}).
		// 	Where("id_card_number = ? AND id <> ?", idc, created.ID).
		// 	Update("sync_until", yesterday).Error; err != nil {
		// 	return err
		// }

		return nil
	}); err != nil {
		return nil, err
	}

	// โหลดพร้อม items
	var out models.Prescription
	if err := db.Preload("Items").First(&out, created.ID).Error; err != nil {
		return nil, err
	}
	res := dto.NewPrescriptionResponse(out)
	return &res, nil
}


// ===== LIST (no pagination) =====
// ค้นหา q ใน id_card_number และกรองตาม doctor_id (optional)
func ListPrescriptions(db *gorm.DB, q string, doctorID *uint) ([]dto.PrescriptionResponse, error) {
	var rows []models.Prescription

	tx := db.Model(&models.Prescription{}).Preload("Items")
	if s := Norm(q); s != "" {
		like := "%" + s + "%"
		tx = tx.Where("id_card_number ILIKE ?", like)
	}
	if doctorID != nil && *doctorID != 0 {
		tx = tx.Where("doctor_id = ?", *doctorID)
	}

	if err := tx.Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]dto.PrescriptionResponse, 0, len(rows))
	for _, m := range rows {
		out = append(out, dto.NewPrescriptionResponse(m))
	}
	return out, nil
}

// ===== GET ONE =====

func GetPrescriptionByID(db *gorm.DB, id uint) (*dto.PrescriptionResponse, error) {
	var rec models.Prescription
	if err := db.Preload("Items").First(&rec, id).Error; err != nil {
		return nil, err
	}
	res := dto.NewPrescriptionResponse(rec)
	return &res, nil
}

// ===== UPDATE (partial เฉพาะหัวเอกสาร) =====
// ถ้าจะอัปเดต items แนะนำทำ endpoint แยก (replace/add/update/delete)
func UpdatePrescription(db *gorm.DB, id uint, in *dto.UpdatePrescriptionDTO) (*dto.PrescriptionResponse, error) {
	if in == nil {
		return nil, errors.New("nil input")
	}

	var rec models.Prescription
	if err := db.First(&rec, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{
		"updated_at": time.Now(),
	}

	if in.IDCardNumber != nil {
		idc := OnlyDigits(*in.IDCardNumber)
		if len(idc) != 13 {
			return nil, errors.New("id_card_number must be 13 digits")
		}
		updates["id_card_number"] = idc
	}

	if in.DoctorID != nil {
		if *in.DoctorID == 0 {
			return nil, errors.New("doctor_id is required")
		}
		var doc models.WebAdmin
		if err := db.Where("role = ?", "doctor").First(&doc, *in.DoctorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("doctor not found")
			}
			return nil, err
		}
		updates["doctor_id"] = *in.DoctorID
	}

	if in.SyncUntil != nil {
		updates["sync_until"] = *in.SyncUntil
	}
	if in.AppSyncStatus != nil {
		updates["app_sync_status"] = *in.AppSyncStatus
	}

	if err := db.Model(&rec).Updates(updates).Error; err != nil {
		return nil, err
	}

	// reload + items
	if err := db.Preload("Items").First(&rec, id).Error; err != nil {
		return nil, err
	}
	res := dto.NewPrescriptionResponse(rec)
	return &res, nil
}

// ===== DELETE (soft delete พร้อมลูก) =====

func DeletePrescription(db *gorm.DB, id uint) error {
	// soft delete ทั้งหัว + associations (Items)
	return db.Select(clause.Associations).Delete(&models.Prescription{ID: id}).Error
}
