package handlers

import (
	"time"

	"github.com/fouradithep/pillmate/models"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// ===== ผูก ID ให้ตรงกับ seed noti_formats =====
const (
	NotiFormatFixedTimes uint = 1 // เวลาเฉพาะ (Fixed Times)
	NotiFormatInterval   uint = 2 // ทุกกี่ชั่วโมง (Interval)
	NotiFormatEveryNDays uint = 3 // วันเว้นวัน / ทุกกี่วัน (EveryNDays)
	NotiFormatCycle      uint = 4 // ทานต่อเนื่อง/พักยา (Cycle)
)

// ===== Helper: ตรวจช่วงวันที่ =====
func validateRange(startDate, endDate time.Time) error {
	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		return gorm.ErrInvalidData
	}
	return nil
}

const dateLayout = "2006-01-02"

// ===================== Helpers: target & uniqueness =====================

// ต้องเลือกอย่างใดอย่างหนึ่ง (XOR) ระหว่าง my_medicine_id กับ group_id
func validateExclusiveTarget(myMedicineID, groupID *uint) bool {
	return (myMedicineID != nil && groupID == nil) || (myMedicineID == nil && groupID != nil)
}

// ห้ามมี NotiInfo ซ้ำสำหรับ target เดียวกัน (ยา 1 ตัวหรือกลุ่ม 1 กลุ่ม)
func ensureNoExistingNotiForTarget(tx *gorm.DB, myMedicineID, groupID *uint) error {
	var count int64
	query := tx.Model(&models.NotiInfo{})
	if myMedicineID != nil {
		query = query.Where("my_medicine_id = ?", *myMedicineID)
	} else {
		query = query.Where("group_id = ?", *groupID)
	}
	if err := query.Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrInvalidData
	}
	return nil
}

// ===================================================================
//                             DTOs (CREATE)
// ===================================================================

// Fixed Times
type CreateNotiFixedTimesReq struct {
	MyMedicineID *uint    `json:"my_medicine_id"`
	GroupID      *uint    `json:"group_id"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Times        []string `json:"times"` // ["08:00","12:00"]
}

// Interval (every N hours)
type CreateNotiIntervalReq struct {
	MyMedicineID  *uint   `json:"my_medicine_id"`
	GroupID       *uint   `json:"group_id"`
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	IntervalHours int     `json:"interval_hours"` // > 0
	TimesPerDay   *int    `json:"times_per_day"`  // optional
}

// Every N days
type CreateNotiEveryNDaysReq struct {
	MyMedicineID *uint    `json:"my_medicine_id"`
	GroupID      *uint    `json:"group_id"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	IntervalDay  int      `json:"interval_day"` // > 0
	Times        []string `json:"times"`
}

// Cycle (e.g., 21 on / 7 off)
type CreateNotiCycleReq struct {
	MyMedicineID *uint   `json:"my_medicine_id"`
	GroupID      *uint   `json:"group_id"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	CyclePattern []int64 `json:"cycle_pattern"` // เช่น [21,7]
	Times        []string `json:"times"`
}

// ===================================================================
//                              CREATE
// ===================================================================

func CreateNotiFixedTimes(db *gorm.DB, req CreateNotiFixedTimesReq) (*models.NotiInfo, error) {
	// ต้องเลือก target อย่างใดอย่างหนึ่ง
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	// ห้ามซ้ำต่อ target
	if err := ensureNoExistingNotiForTarget(db, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
	}

	if len(req.Times) == 0 {
		return nil, gorm.ErrInvalidData
	}

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}

	if err := validateRange(startDate, endDate); err != nil {
		return nil, err
	}

	timesArr := pq.StringArray(req.Times)
	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     startDate,
		EndDate:       endDate,
		NotiFormatID:  NotiFormatFixedTimes,
		Times:         &timesArr,
		IntervalHours: nil,
		TimesPerDay:   nil,
		IntervalDay:   nil,
		CyclePattern:  nil,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateNotiInterval(db *gorm.DB, req CreateNotiIntervalReq) (*models.NotiInfo, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureNoExistingNotiForTarget(db, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
	}

	if req.IntervalHours <= 0 {
		return nil, gorm.ErrInvalidData
	}

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}

	if err := validateRange(startDate, endDate); err != nil {
		return nil, err
	}

	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     startDate,
		EndDate:       endDate,
		NotiFormatID:  NotiFormatInterval,
		IntervalHours: &req.IntervalHours,
		TimesPerDay:   req.TimesPerDay,
		Times:         nil,
		IntervalDay:   nil,
		CyclePattern:  nil,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateNotiEveryNDays(db *gorm.DB, req CreateNotiEveryNDaysReq) (*models.NotiInfo, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureNoExistingNotiForTarget(db, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
	}

	if req.IntervalDay <= 0 || len(req.Times) == 0 {
		return nil, gorm.ErrInvalidData
	}

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}

	if err := validateRange(startDate, endDate); err != nil {
		return nil, err
	}

	intervalDay := req.IntervalDay
	timesArr := pq.StringArray(req.Times)
	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     startDate,
		EndDate:       endDate,
		NotiFormatID:  NotiFormatEveryNDays,
		IntervalDay:   &intervalDay,
		Times:         &timesArr,
		IntervalHours: nil,
		TimesPerDay:   nil,
		CyclePattern:  nil,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateNotiCycle(db *gorm.DB, req CreateNotiCycleReq) (*models.NotiInfo, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureNoExistingNotiForTarget(db, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
	}

	if len(req.CyclePattern) == 0 || len(req.Times) == 0 {
		return nil, gorm.ErrInvalidData
	}

	startDate, err := time.Parse(dateLayout, req.StartDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}
	endDate, err := time.Parse(dateLayout, req.EndDate)
	if err != nil {
		return nil, gorm.ErrInvalidData
	}

	if err := validateRange(startDate, endDate); err != nil {
		return nil, err
	}

	cyclePatternArr := pq.Int64Array(req.CyclePattern)
	timesArr := pq.StringArray(req.Times)
	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     startDate,
		EndDate:       endDate,
		NotiFormatID:  NotiFormatCycle,
		CyclePattern:  &cyclePatternArr,
		Times:         &timesArr,
		IntervalHours: nil,
		TimesPerDay:   nil,
		IntervalDay:   nil,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// ===================================================================
//                           LIST / GET
// ===================================================================

func ListNotiInfos(db *gorm.DB, filter map[string]any) ([]models.NotiInfo, error) {
	query := db.Model(&models.NotiInfo{}).
		Preload("NotiFormat").
		Preload("MyMedicine").
		Preload("Group")

	for col, val := range filter {
		query = query.Where(col+" = ?", val)
	}

	var items []models.NotiInfo
	if err := query.Order("id").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func GetNotiInfo(db *gorm.DB, id uint) (*models.NotiInfo, error) {
	var item models.NotiInfo
	if err := db.Preload("NotiFormat").
		Preload("MyMedicine").
		Preload("Group").
		First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// ===================================================================
//                             DELETE
// ===================================================================

func DeleteNotiInfo(db *gorm.DB, id uint) error {
	return db.Delete(&models.NotiInfo{}, id).Error
}
