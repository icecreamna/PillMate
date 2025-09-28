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
func validateRange(s, e time.Time) error {
	if !s.IsZero() && !e.IsZero() && e.Before(s) {
		return gorm.ErrInvalidData
	}
	return nil
}

// ===================================================================
//                             DTOs (CREATE)
// ===================================================================

// Fixed Times
type CreateNotiFixedTimesReq struct {
	MyMedicineID *uint     `json:"my_medicine_id"`
	GroupID      *uint     `json:"group_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	Times        []string  `json:"times"` // ["08:00","12:00"]
}

// Interval (every N hours)
type CreateNotiIntervalReq struct {
	MyMedicineID  *uint     `json:"my_medicine_id"`
	GroupID       *uint     `json:"group_id"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	IntervalHours int       `json:"interval_hours"` // > 0
	TimesPerDay   *int      `json:"times_per_day"`  // optional
}

// Every N days
type CreateNotiEveryNDaysReq struct {
	MyMedicineID *uint     `json:"my_medicine_id"`
	GroupID      *uint     `json:"group_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	IntervalDay  int       `json:"interval_day"` // > 0
	Times        []string  `json:"times"`
}

// Cycle (e.g., 21 on / 7 off)
type CreateNotiCycleReq struct {
	MyMedicineID *uint     `json:"my_medicine_id"`
	GroupID      *uint     `json:"group_id"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	CyclePattern []int64   `json:"cycle_pattern"` // เช่น [21,7]
	Times        []string  `json:"times"`
}

// ===================================================================
//                              CREATE
// ===================================================================

func CreateNotiFixedTimes(db *gorm.DB, req CreateNotiFixedTimesReq) (*models.NotiInfo, error) {
	if len(req.Times) == 0 {
		return nil, gorm.ErrInvalidData
	}
	if err := validateRange(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	t := pq.StringArray(req.Times)
	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		NotiFormatID:  NotiFormatFixedTimes,
		Times:         &t,
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
	if req.IntervalHours <= 0 {
		return nil, gorm.ErrInvalidData
	}
	if err := validateRange(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
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
	if req.IntervalDay <= 0 || len(req.Times) == 0 {
		return nil, gorm.ErrInvalidData
	}
	if err := validateRange(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	i := req.IntervalDay
	t := pq.StringArray(req.Times)
	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		NotiFormatID:  NotiFormatEveryNDays,
		IntervalDay:   &i,
		Times:         &t,
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
	if len(req.CyclePattern) == 0 || len(req.Times) == 0 {
		return nil, gorm.ErrInvalidData
	}
	if err := validateRange(req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	cp := pq.Int64Array(req.CyclePattern)
	t := pq.StringArray(req.Times)
	item := models.NotiInfo{
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		NotiFormatID:  NotiFormatCycle,
		CyclePattern:  &cp,
		Times:         &t,
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
	q := db.Model(&models.NotiInfo{}).
		Preload("NotiFormat").
		Preload("MyMedicine").
		Preload("Group")

	for k, v := range filter {
		q = q.Where(k+" = ?", v)
	}

	var items []models.NotiInfo
	if err := q.Order("id").Find(&items).Error; err != nil {
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
