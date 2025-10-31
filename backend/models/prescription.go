package models

import (
	"time"

	"gorm.io/gorm"
)

// ใบสั่งยา (หัวเอกสาร) — 1 ใบสั่งยา สามารถมีหลายรายการยา (PrescriptionItem)
type Prescription struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	IDCardNumber string         `gorm:"type:varchar(13);not null;index" json:"id_card_number"`
	DoctorID     uint           `gorm:"not null;index" json:"doctor_id"`

	// timestamps + soft delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	AppSyncStatus bool      `gorm:"default:false" json:"app_sync_status"` // false=ยังไม่ซิงค์
	SyncUntil     time.Time `gorm:"not null;index" json:"sync_until"`     // = MAX(expire_date) ของรายการยา (ถ้าไม่มี ให้ +60 วัน)

	// รายการยาในใบสั่งนี้ (ลบหัว → ลบลูกด้วย; hard delete จะ CASCADE โดย DB)
	Items []PrescriptionItem `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items"`

	// ความสัมพันธ์อื่น (เลือก preload ตามต้องการ)
	WebAdmin WebAdmin `gorm:"foreignKey:DoctorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

// รายการยาในใบสั่ง (หนึ่งแถวต่อยาหนึ่งตัว/หนึ่งโดส)
type PrescriptionItem struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	PrescriptionID uint           `gorm:"not null;index" json:"prescription_id"`
	MedicineInfoID uint           `gorm:"not null;index" json:"medicine_info_id"`
	AmountPerTime  string         `gorm:"not null" json:"amount_per_time"` // เช่น "1 เม็ด"
	TimesPerDay    string         `gorm:"not null" json:"times_per_day"`   // เช่น "2 ครั้ง"

	// เก็บเป็น "YYYY-MM-DD" (10 ตัว)
	StartDate  *string `gorm:"type:varchar(10)" json:"start_date,omitempty"`
	EndDate    *string `gorm:"type:varchar(10)" json:"end_date,omitempty"`
	ExpireDate *string `gorm:"type:varchar(10)" json:"expire_date,omitempty"` // = EndDate + 1 วัน (คำนวณอัตโนมัติ)
	Note       *string `gorm:"type:text"        json:"note,omitempty"`

	// timestamps + soft delete
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// ความสัมพันธ์อ้างอิงข้อมูลยา
	MedicineInfo MedicineInfo `gorm:"foreignKey:MedicineInfoID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

/* ===================== Helpers (ภายในไฟล์) ===================== */

func addOneDay(dateStr string) (string, bool) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return "", false
	}
	next := t.AddDate(0, 0, 1)
	return next.Format("2006-01-02"), true
}

// แปลง "YYYY-MM-DD" → time.Time (00:00 น. ตาม Asia/Bangkok)
func dateStrToLocalMidnight(dateStr string) (time.Time, bool) {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return time.Time{}, false
	}
	// ชัดเจนว่าเป็น 00:00 ของวันนั้น (consistency กับ where sync_until >= today)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc), true
}

// คำนวณ MAX(expire_date) ของ prescription_id แล้วอัปเดต prescription.sync_until
func updateParentSyncUntil(tx *gorm.DB, prescriptionID uint) error {
	// อ่าน MAX(expire_date) (string) จาก items
	var maxExpire *string
	err := tx.Model(&PrescriptionItem{}).
		Where("prescription_id = ? AND deleted_at IS NULL", prescriptionID).
		Select("MAX(expire_date)").
		Scan(&maxExpire).Error
	if err != nil {
		return err
	}

	var target time.Time
	if maxExpire != nil && *maxExpire != "" {
		if t, ok := dateStrToLocalMidnight(*maxExpire); ok {
			target = t
		}
	}

	// ถ้ายังไม่ได้ค่า (เช่น ไม่มีรายการ/ไม่มี expire_date) ให้ default +60 วันจากวันนี้
	if target.IsZero() {
		target = time.Now().AddDate(0, 0, 60)
	}

	return tx.Model(&Prescription{}).
		Where("id = ?", prescriptionID).
		Update("sync_until", target).Error
}

/* ===================== Hooks ===================== */

// ตั้ง ExpireDate อัตโนมัติจาก EndDate ก่อนเซฟ
func (pi *PrescriptionItem) BeforeSave(tx *gorm.DB) error {
	if pi.EndDate != nil && *pi.EndDate != "" {
		if next, ok := addOneDay(*pi.EndDate); ok {
			pi.ExpireDate = &next
		}
	}
	return nil
}

// หลังจากสร้าง/อัปเดต PrescriptionItem → ดัน sync_until ของหัวให้เป็น MAX(expire_date)
func (pi *PrescriptionItem) AfterSave(tx *gorm.DB) error {
	if pi.PrescriptionID == 0 {
		return nil
	}
	return updateParentSyncUntil(tx, pi.PrescriptionID)
}

// หลังจากลบ (soft delete) PrescriptionItem → คำนวณใหม่เหมือนกัน
func (pi *PrescriptionItem) AfterDelete(tx *gorm.DB) error {
	if pi.PrescriptionID == 0 {
		return nil
	}
	return updateParentSyncUntil(tx, pi.PrescriptionID)
}

// ตั้งค่าเริ่มต้นตอนสร้างใบสั่ง (กรณีไม่มีรายการ/ไม่มี expire_date)
// *หมายเหตุ:* ถ้ามีการเพิ่ม items ตามมา Hooks ของ PrescriptionItem จะอัปเดต sync_until ให้อัตโนมัติ
func (p *Prescription) BeforeCreate(tx *gorm.DB) error {
	if p.SyncUntil.IsZero() {
		p.SyncUntil = time.Now().AddDate(0, 0, 60)
	}
	return nil
}
