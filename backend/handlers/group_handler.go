package handlers

import (
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

//
// CREATE: สร้าง Group + ใส่ยาหลายตัวให้กลุ่มนั้น (หนึ่งยา = หนึ่งกลุ่ม)
//
type CreateGroupRequest struct {
	GroupName     string `json:"group_name"`
	MyMedicineIDs []uint `json:"my_medicine_ids"`
}

func CreateGroup(db *gorm.DB, patientID uint, req CreateGroupRequest) (*models.Group, []models.MyMedicine, error) {
	if req.GroupName == "" || len(req.MyMedicineIDs) == 0 {
		return nil, nil, gorm.ErrInvalidData
	}

	var group models.Group
	var updated []models.MyMedicine

	err := db.Transaction(func(tx *gorm.DB) error {
		// สร้าง/ใช้กลุ่ม (unique ภายใน patient)
		group = models.Group{PatientID: patientID, GroupName: req.GroupName}
		if err := tx.Where("patient_id = ? AND group_name = ?", patientID, req.GroupName).
			FirstOrCreate(&group).Error; err != nil {
			return err
		}

		// ต้องเป็นยาของ patient และยังไม่ถูกจัดกลุ่ม (group_id IS NULL)
		var cnt int64
		if err := tx.Model(&models.MyMedicine{}).
			Where("patient_id = ? AND id IN ? AND group_id IS NULL", patientID, req.MyMedicineIDs).
			Count(&cnt).Error; err != nil {
			return err
		}
		if cnt != int64(len(req.MyMedicineIDs)) {
			return gorm.ErrInvalidData
		}

		// ใส่กลุ่มให้ยา
		if err := tx.Model(&models.MyMedicine{}).
			Where("patient_id = ? AND id IN ?", patientID, req.MyMedicineIDs).
			Update("group_id", group.ID).Error; err != nil {
			return err
		}

		// โหลดสมาชิกที่อัปเดตไว้ตอบกลับ
		return tx.Where("patient_id = ? AND group_id = ?", patientID, group.ID).
			Find(&updated).Error
	})
	if err != nil {
		return nil, nil, err
	}
	return &group, updated, nil
}

//
// READ: ดึงรายละเอียดกลุ่ม + สมาชิก (ยาของกลุ่มนั้น)
//
type GroupDetail struct {
	Group    models.Group        `json:"group"`
	Members  []models.MyMedicine `json:"members"`
}

func GetGroup(db *gorm.DB, patientID, groupID uint) (*GroupDetail, error) {
	var g models.Group
	if err := db.Where("id = ? AND patient_id = ?", groupID, patientID).First(&g).Error; err != nil {
		return nil, err
	}
	var members []models.MyMedicine
	if err := db.Where("patient_id = ? AND group_id = ?", patientID, groupID).Find(&members).Error; err != nil {
		return nil, err
	}
	return &GroupDetail{Group: g, Members: members}, nil
}

//
// READ: ดึงรายการกลุ่มทั้งหมดของผู้ป่วย (ไม่รวมสมาชิก)
//
func GetGroups(db *gorm.DB, patientID uint) ([]models.Group, error) {
	var groups []models.Group
	if err := db.Where("patient_id = ?", patientID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

//
// UPDATE: เปลี่ยนชื่อกลุ่ม (optional) + ตั้งสมาชิก “ชุดสุดท้าย”
//
type UpdateGroupRequest struct {
	NewGroupName  *string `json:"new_group_name"`  // optional
	MyMedicineIDs []uint  `json:"my_medicine_ids"` // required: ชุดสมาชิกสุดท้าย
}

func UpdateGroup(db *gorm.DB, patientID, groupID uint, req UpdateGroupRequest) (*models.Group, []models.MyMedicine, error) {
	if len(req.MyMedicineIDs) == 0 {
		return nil, nil, gorm.ErrInvalidData
	}

	var group models.Group
	var members []models.MyMedicine

	err := db.Transaction(func(tx *gorm.DB) error {
		// กลุ่มต้องเป็นของ patient
		if err := tx.Where("id = ? AND patient_id = ?", groupID, patientID).First(&group).Error; err != nil {
			return err
		}

		// เปลี่ยนชื่อ (กันชื่อชนภายใน patient)
		if req.NewGroupName != nil && *req.NewGroupName != "" && *req.NewGroupName != group.GroupName {
			var exists int64
			if err := tx.Model(&models.Group{}).
				Where("patient_id = ? AND group_name = ? AND id <> ?", patientID, *req.NewGroupName, group.ID).
				Count(&exists).Error; err != nil {
				return err
			}
			if exists > 0 {
				return gorm.ErrInvalidData
			}
			if err := tx.Model(&group).Update("group_name", *req.NewGroupName).Error; err != nil {
				return err
			}
		}

		// ปลดสมาชิกเดิมทั้งหมดของกลุ่มนี้ก่อน
		if err := tx.Model(&models.MyMedicine{}).
			Where("patient_id = ? AND group_id = ?", patientID, group.ID).
			Update("group_id", nil).Error; err != nil {
			return err
		}

		// ยืนยันว่า ids ที่จะใส่ทั้งหมดเป็นของ patient และยังไม่อยู่กลุ่มไหน
		var cnt int64
		if err := tx.Model(&models.MyMedicine{}).
			Where("patient_id = ? AND id IN ? AND group_id IS NULL", patientID, req.MyMedicineIDs).
			Count(&cnt).Error; err != nil {
			return err
		}
		if cnt != int64(len(req.MyMedicineIDs)) {
			return gorm.ErrInvalidData
		}

		// ใส่สมาชิกใหม่ทั้งหมดให้กลุ่มนี้
		if err := tx.Model(&models.MyMedicine{}).
			Where("patient_id = ? AND id IN ?", patientID, req.MyMedicineIDs).
			Update("group_id", group.ID).Error; err != nil {
			return err
		}

		// โหลดสมาชิกล่าสุด
		if err := tx.Where("patient_id = ? AND group_id = ?", patientID, group.ID).
			Find(&members).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return &group, members, nil
}

//
// DELETE: ลบกลุ่ม (ถอด group_id ออกจากยาทั้งหมดในกลุ่มก่อน)
//
func DeleteGroup(db *gorm.DB, patientID, groupID uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// กลุ่มต้องเป็นของ patient
		var g models.Group
		if err := tx.Where("id = ? AND patient_id = ?", groupID, patientID).First(&g).Error; err != nil {
			return err
		}

		// ถอดสมาชิก (set group_id = NULL)
		if err := tx.Model(&models.MyMedicine{}).
			Where("patient_id = ? AND group_id = ?", patientID, groupID).
			Update("group_id", nil).Error; err != nil {
			return err
		}

		// ลบกลุ่ม
		if err := tx.Delete(&g).Error; err != nil {
			return err
		}
		return nil
	})
}

//
// READ: ยาที่ยังไม่มีกลุ่ม (เลือกได้)
//
func GetUngroupedMyMedicines(db *gorm.DB, patientID uint) ([]models.MyMedicine, error) {
	var list []models.MyMedicine
	if err := db.
		Where("patient_id = ? AND group_id IS NULL", patientID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
