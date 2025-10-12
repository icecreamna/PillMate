package handlers

import (
	"fmt"
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

const dateLayout = "2006-01-02"

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
	MyMedicineID  *uint    `json:"my_medicine_id"`
	GroupID       *uint    `json:"group_id"`
	StartDate     string   `json:"start_date"`
	EndDate       string   `json:"end_date"`
	IntervalHours int      `json:"interval_hours"` // > 0
	TimesPerDay   *int     `json:"times_per_day"`
	Times         []string `json:"times"` // optional
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
	MyMedicineID *uint    `json:"my_medicine_id"`
	GroupID      *uint    `json:"group_id"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	CyclePattern []int64  `json:"cycle_pattern"` // เช่น [21,7]
	Times        []string `json:"times"`
}

// ===================================================================
//                             DTOs (RESP)
// ===================================================================

type NotiInfoResp struct {
	ID             uint     `json:"id"`
	MyMedicineID   *uint    `json:"my_medicine_id,omitempty"`
	GroupID        *uint    `json:"group_id,omitempty"`
	NotiFormatID   uint     `json:"noti_format_id"`
	NotiFormatName string   `gorm:"column:noti_format_name"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	IntervalHours  *int     `json:"interval_hours,omitempty"`
	IntervalDay    *int     `json:"interval_day,omitempty"`
	TimesPerDay    *int     `json:"times_per_day,omitempty"` // optional
	Times          []string `json:"times,omitempty"`
	CyclePattern   []int    `json:"cycle_pattern,omitempty"`
}

// ===================== Helpers (ทั่วไป) =====================

// ตรวจช่วงวันที่
func validateRange(startDate, endDate time.Time) error {
	if !startDate.IsZero() && !endDate.IsZero() && endDate.Before(startDate) {
		return gorm.ErrInvalidData
	}
	return nil
}

// ต้องเลือกอย่างใดอย่างหนึ่ง (XOR) ระหว่าง my_medicine_id กับ group_id
func validateExclusiveTarget(myMedicineID, groupID *uint) bool {
	return (myMedicineID != nil && groupID == nil) || (myMedicineID == nil && groupID != nil)
}

// target (ยา/กลุ่ม) เป็นของ patientID จริงหรือไม่
func ensureTargetBelongsToPatient(db *gorm.DB, patientID uint, myMedicineID, groupID *uint) error {
	if myMedicineID != nil {
		var cnt int64
		if err := db.Model(&models.MyMedicine{}).
			Where("id = ? AND patient_id = ?", *myMedicineID, patientID).
			Count(&cnt).Error; err != nil {
			return err
		}
		if cnt == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	}

	// group
	var cnt int64
	if err := db.Model(&models.Group{}).
		Where("id = ? AND patient_id = ?", *groupID, patientID).
		Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ห้ามมี NotiInfo ซ้ำสำหรับ target เดียวกัน (ไม่ต้องมี patient_id ใน noti_infos)
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

// map model -> RESP DTO
func toNotiInfoResp(n *models.NotiInfo) *NotiInfoResp {
	var times []string
	if n.Times != nil {
		times = []string(*n.Times)
	}

	var cycle []int
	if n.CyclePattern != nil {
		cycle = make([]int, len(*n.CyclePattern))
		for i, v := range *n.CyclePattern {
			cycle[i] = int(v)
		}
	}

	return &NotiInfoResp{
		ID:            n.ID,
		MyMedicineID:  n.MyMedicineID,
		GroupID:       n.GroupID,
		NotiFormatID:  n.NotiFormatID,
		StartDate:     n.StartDate.Format(dateLayout),
		EndDate:       n.EndDate.Format(dateLayout),
		IntervalHours: n.IntervalHours,
		IntervalDay:   n.IntervalDay,
		Times:         times,
		CyclePattern:  cycle,
	}
}

// ===================================================================
//                              CREATE
// ===================================================================

func CreateNotiFixedTimes(db *gorm.DB, patientID uint, req CreateNotiFixedTimesReq) (*NotiInfoResp, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureTargetBelongsToPatient(db, patientID, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
	}
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
		PatientID:     patientID,
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
	return toNotiInfoResp(&item), nil
}

func CreateNotiInterval(db *gorm.DB, patientID uint, req CreateNotiIntervalReq) (*NotiInfoResp, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureTargetBelongsToPatient(db, patientID, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
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

	timesArr := pq.StringArray(req.Times)

	item := models.NotiInfo{
		PatientID:     patientID,
		MyMedicineID:  req.MyMedicineID,
		GroupID:       req.GroupID,
		StartDate:     startDate,
		EndDate:       endDate,
		NotiFormatID:  NotiFormatInterval,
		IntervalHours: &req.IntervalHours,
		TimesPerDay:   req.TimesPerDay,
		Times:         &timesArr,
		IntervalDay:   nil,
		CyclePattern:  nil,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return toNotiInfoResp(&item), nil
}

func CreateNotiEveryNDays(db *gorm.DB, patientID uint, req CreateNotiEveryNDaysReq) (*NotiInfoResp, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureTargetBelongsToPatient(db, patientID, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
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
		PatientID:     patientID,
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
	return toNotiInfoResp(&item), nil
}

func CreateNotiCycle(db *gorm.DB, patientID uint, req CreateNotiCycleReq) (*NotiInfoResp, error) {
	if !validateExclusiveTarget(req.MyMedicineID, req.GroupID) {
		return nil, gorm.ErrInvalidData
	}
	if err := ensureTargetBelongsToPatient(db, patientID, req.MyMedicineID, req.GroupID); err != nil {
		return nil, err
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
		PatientID:     patientID,
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
	return toNotiInfoResp(&item), nil
}

// ===================================================================
//                           LIST / GET
// ===================================================================

// ดึงเฉพาะของผู้ป่วยคนนั้นๆ ด้วย EXISTS (ไม่ preload)
func ListNotiInfos(db *gorm.DB, patientID uint, filter map[string]any) ([]NotiInfoResp, error) {
	q := db.Model(&models.NotiInfo{}).
		Where(`
			(
				my_medicine_id IS NOT NULL AND EXISTS (
					SELECT 1 FROM my_medicines m
					WHERE m.id = noti_infos.my_medicine_id AND m.patient_id = ?
				)
			) OR
			(
				group_id IS NOT NULL AND EXISTS (
					SELECT 1 FROM groups g
					WHERE g.id = noti_infos.group_id AND g.patient_id = ?
				)
			)
		`, patientID, patientID)

	// filter เพิ่ม: block การส่ง patient_id มาทับ
	for col, val := range filter {
		if col == "patient_id" {
			continue
		}
		q = q.Where(col+" = ?", val)
	}

	type row struct {
		ID            uint           `gorm:"column:id"`
		MyMedicineID  *uint          `gorm:"column:my_medicine_id"`
		GroupID       *uint          `gorm:"column:group_id"`
		NotiFormatID  uint           `gorm:"column:noti_format_id"`
		StartDate     string         `gorm:"column:start_date"`
		EndDate       string         `gorm:"column:end_date"`
		IntervalHours *int           `gorm:"column:interval_hours"`
		IntervalDay   *int           `gorm:"column:interval_day"`
		TimesPerDay   *int           `gorm:"column:times_per_day"`
		Times         pq.StringArray `gorm:"column:times"`
		CyclePattern  pq.Int64Array  `gorm:"column:cycle_pattern"`
	}

	q = q.Select(`
		id,
		my_medicine_id,
		group_id,
		noti_format_id,
		to_char(start_date,'YYYY-MM-DD') AS start_date,
		to_char(end_date,'YYYY-MM-DD')   AS end_date,
		interval_hours,
		interval_day,
		times_per_day,
		times,
		cycle_pattern
	`).Order("id")

	var rows []row
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]NotiInfoResp, 0, len(rows))
	for _, r := range rows {
		var times []string
		if r.Times != nil {
			times = []string(r.Times)
		}

		var cycle []int
		if r.CyclePattern != nil {
			cycle = make([]int, len(r.CyclePattern))
			for i, v := range r.CyclePattern {
				cycle[i] = int(v)
			}
		}

		out = append(out, NotiInfoResp{
			ID:            r.ID,
			MyMedicineID:  r.MyMedicineID,
			GroupID:       r.GroupID,
			NotiFormatID:  r.NotiFormatID,
			StartDate:     r.StartDate,
			EndDate:       r.EndDate,
			IntervalHours: r.IntervalHours,
			IntervalDay:   r.IntervalDay,
			TimesPerDay:   r.TimesPerDay,
			Times:         times,
			CyclePattern:  cycle,
		})
	}
	return out, nil
}

// Get รายการเดียวของผู้ป่วยคนนั้น (where id + EXISTS)
func GetNotiInfoByType(db *gorm.DB, patientID uint, typ string, id uint) (*NotiInfoResp, error) {
	type row struct {
		ID             uint           `gorm:"column:id"`
		MyMedicineID   *uint          `gorm:"column:my_medicine_id"`
		GroupID        *uint          `gorm:"column:group_id"`
		NotiFormatID   uint           `gorm:"column:noti_format_id"`
		NotiFormatName string         `gorm:"column:noti_format_name"`
		StartDate      string         `gorm:"column:start_date"`
		EndDate        string         `gorm:"column:end_date"`
		IntervalHours  *int           `gorm:"column:interval_hours"`
		IntervalDay    *int           `gorm:"column:interval_day"`
		TimesPerDay    *int           `gorm:"column:times_per_day"`
		Times          pq.StringArray `gorm:"column:times"`
		CyclePattern   pq.Int64Array  `gorm:"column:cycle_pattern"`
	}

	var q *gorm.DB

	switch typ {
	case "medicine":
		q = db.Table("noti_infos ni").
			Joins("LEFT JOIN noti_formats nf ON nf.id = ni.noti_format_id").
			Where("ni.my_medicine_id = ?", id).
			Where("ni.deleted_at IS NULL").
			Where("EXISTS (SELECT 1 FROM my_medicines m WHERE m.id = ni.my_medicine_id AND m.patient_id = ?)", patientID)

	case "group":
		q = db.Table("noti_infos ni").
			Joins("LEFT JOIN noti_formats nf ON nf.id = ni.noti_format_id").
			Where("ni.group_id = ?", id).
			Where("ni.deleted_at IS NULL").
			Where("EXISTS (SELECT 1 FROM groups g WHERE g.id = ni.group_id AND g.patient_id = ?)", patientID)

	default:
		return nil, fmt.Errorf("invalid type (must be 'medicine' or 'group')")
	}

	q = q.Select(`
		ni.id,
		ni.my_medicine_id,
		ni.group_id,
		ni.noti_format_id,
		nf.format_name AS noti_format_name,
		to_char(ni.start_date,'YYYY-MM-DD') AS start_date,
		to_char(ni.end_date,'YYYY-MM-DD')   AS end_date,
		ni.interval_hours,
		ni.interval_day,
		ni.times_per_day,
		ni.times,
		ni.cycle_pattern
	`)

	var r row
	if err := q.Take(&r).Error; err != nil {
		return nil, err
	}

	// 🧩 แปลง array จาก SQL → Go
	var times []string
	if r.Times != nil {
		times = []string(r.Times)
	}

	var cycle []int
	if r.CyclePattern != nil {
		cycle = make([]int, len(r.CyclePattern))
		for i, v := range r.CyclePattern {
			cycle[i] = int(v)
		}
	}

	out := &NotiInfoResp{
		ID:             r.ID,
		MyMedicineID:   r.MyMedicineID,
		GroupID:        r.GroupID,
		NotiFormatID:   r.NotiFormatID,
		NotiFormatName: r.NotiFormatName,
		StartDate:      r.StartDate,
		EndDate:        r.EndDate,
		IntervalHours:  r.IntervalHours,
		IntervalDay:    r.IntervalDay,
		TimesPerDay:    r.TimesPerDay,
		Times:          times,
		CyclePattern:   cycle,
	}

	return out, nil
}

// ===================================================================
//                             DELETE
// ===================================================================

// ลบเฉพาะเรคคอร์ดที่เป็นของผู้ป่วยคนนั้น (where id + EXISTS)
// และลบ NotiItem ที่ผูกกับ NotiInfo นั้น ๆ ไปด้วย
func DeleteNotiInfo(db *gorm.DB, patientID, id uint) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// หา noti_info.id ที่ "ผู้ป่วยคนนี้" มีสิทธิ์ลบ (จำกัดด้วย EXISTS)
	var allowedID uint
	if err := tx.Model(&models.NotiInfo{}).
		Select("id").
		Where("id = ?", id).
		Where(`
			(
				my_medicine_id IS NOT NULL AND EXISTS (
					SELECT 1 FROM my_medicines m
					WHERE m.id = noti_infos.my_medicine_id AND m.patient_id = ?
				)
			) OR
			(
				group_id IS NOT NULL AND EXISTS (
					SELECT 1 FROM groups g
					WHERE g.id = noti_infos.group_id AND g.patient_id = ?
				)
			)
		`, patientID, patientID).
		Take(&allowedID).Error; err != nil {
		// ไม่พบ หรือไม่มีสิทธิ์
		tx.Rollback()
		return err
	}

	// ลบลูก (noti_items) ที่อ้างถึง noti_info_id นี้ก่อน
	if err := tx.Where("noti_info_id = ?", allowedID).
		Delete(&models.NotiItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// ลบแม่ (noti_infos)
	if err := tx.Where("id = ?", allowedID).
		Delete(&models.NotiInfo{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
