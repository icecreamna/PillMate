package handlers

import (
	"errors"
	"sort"
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const ymdLayout = "2006-01-02"

// ===================================================================
//                     Expand: แตกเวลาจาก NotiInfo
// ===================================================================

// แยกชั่วโมง/นาทีจาก "HH:MM"
func splitHourMinute(hhmm string) (int, int, error) {
	if len(hhmm) != 5 || hhmm[2] != ':' {
		return 0, 0, errors.New("invalid HH:MM")
	}
	h := int(hhmm[0]-'0')*10 + int(hhmm[1]-'0')
	m := int(hhmm[3]-'0')*10 + int(hhmm[4]-'0')
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, errors.New("invalid HH:MM range")
	}
	return h, m, nil
}

// สร้างเวลาจากวันที่ + "HH:MM" โดยยึด time.Local
func makeLocalDateTime(day time.Time, hhmm string) (time.Time, error) {
	year, month, date := day.In(time.Local).Date()
	hour, minute, err := splitHourMinute(hhmm)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, month, date, hour, minute, 0, 0, time.Local), nil
}

// ตัดช่วงซ้อนระหว่าง [aStart,aEnd] กับ [bStart,bEnd]
func intersectRange(aStart, aEnd, bStart, bEnd time.Time) (time.Time, time.Time, bool) {
	if aEnd.Before(bStart) || bEnd.Before(aStart) {
		return time.Time{}, time.Time{}, false
	}
	if aStart.Before(bStart) {
		aStart = bStart
	}
	if aEnd.After(bEnd) {
		aEnd = bEnd
	}
	return aStart, aEnd, true
}

// ExpandNotiInfoOccurrences: คืน array เวลาที่ต้องแจ้งเตือนของ notiInfo ภายใน [winStart, winEnd]
func ExpandNotiInfoOccurrences(notiInfo *models.NotiInfo, winStart, winEnd time.Time) ([]time.Time, error) {
	if notiInfo == nil {
		return nil, errors.New("nil notiInfo")
	}

	ruleStart := time.Date(notiInfo.StartDate.Year(), notiInfo.StartDate.Month(), notiInfo.StartDate.Day(), 0, 0, 0, 0, time.Local)
	ruleEnd := time.Date(notiInfo.EndDate.Year(), notiInfo.EndDate.Month(), notiInfo.EndDate.Day(), 23, 59, 59, 0, time.Local)

	windowStart := time.Date(winStart.Year(), winStart.Month(), winStart.Day(), 0, 0, 0, 0, time.Local)
	windowEnd := time.Date(winEnd.Year(), winEnd.Month(), winEnd.Day(), 23, 59, 59, 0, time.Local)

	start, end, ok := intersectRange(ruleStart, ruleEnd, windowStart, windowEnd)
	if !ok {
		return nil, nil
	}

	var times []time.Time
	appendTime := func(t time.Time) { times = append(times, t) }

	switch notiInfo.NotiFormatID {
	case NotiFormatFixedTimes:
		if notiInfo.Times == nil || len(*notiInfo.Times) == 0 {
			return nil, errors.New("times required")
		}
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			for _, s := range *notiInfo.Times {
				ts, err := makeLocalDateTime(d, s)
				if err != nil {
					return nil, err
				}
				appendTime(ts)
			}
		}

	case NotiFormatEveryNDays:
		if notiInfo.IntervalDay == nil || *notiInfo.IntervalDay <= 0 {
			return nil, errors.New("interval_day required")
		}
		if notiInfo.Times == nil || len(*notiInfo.Times) == 0 {
			return nil, errors.New("times required")
		}
		for d := ruleStart; !d.After(end); d = d.AddDate(0, 0, *notiInfo.IntervalDay) {
			if d.Before(start) {
				continue
			}
			for _, s := range *notiInfo.Times {
				ts, err := makeLocalDateTime(d, s)
				if err != nil {
					return nil, err
				}
				appendTime(ts)
			}
		}

	case NotiFormatInterval:
		if notiInfo.IntervalHours == nil || *notiInfo.IntervalHours <= 0 {
			return nil, errors.New("interval_hours required")
		}
		step := time.Duration(*notiInfo.IntervalHours) * time.Hour
		cursor := start

		var perDayLimit int
		if notiInfo.TimesPerDay != nil && *notiInfo.TimesPerDay > 0 {
			perDayLimit = *notiInfo.TimesPerDay
		}
		countToday := 0
		curDay := time.Date(cursor.Year(), cursor.Month(), cursor.Day(), 0, 0, 0, 0, time.Local)

		for !cursor.After(end) {
			day := time.Date(cursor.Year(), cursor.Month(), cursor.Day(), 0, 0, 0, 0, time.Local)
			if day.After(curDay) {
				curDay = day
				countToday = 0
			}
			if cursor.Before(start) {
				cursor = cursor.Add(step)
				continue
			}
			if perDayLimit == 0 || countToday < perDayLimit {
				appendTime(cursor)
				countToday++
			}
			cursor = cursor.Add(step)
		}

	case NotiFormatCycle:
		if notiInfo.CyclePattern == nil || len(*notiInfo.CyclePattern) == 0 {
			return nil, errors.New("cycle_pattern required")
		}
		if notiInfo.Times == nil || len(*notiInfo.Times) == 0 {
			return nil, errors.New("times required")
		}
		pattern := *notiInfo.CyclePattern // ตัวอย่าง [21,7] = on/off
		cursor := ruleStart
		index := 0
		for !cursor.After(end) {
			// on days
			onDays := int(pattern[index%len(pattern)])
			for i := 0; i < onDays && !cursor.After(end); i++ {
				if !cursor.Before(start) {
					for _, s := range *notiInfo.Times {
						ts, err := makeLocalDateTime(cursor, s)
						if err != nil {
							return nil, err
						}
						appendTime(ts)
					}
				}
				cursor = cursor.AddDate(0, 0, 1)
			}
			index++
			// off days
			if index < len(pattern) {
				offDays := int(pattern[index%len(pattern)])
				cursor = cursor.AddDate(0, 0, offDays)
				index++
			} else if len(pattern) >= 2 {
				// วนซ้ำค่า off ล่าสุด
				cursor = cursor.AddDate(0, 0, int(pattern[1]))
			}
		}

	default:
		return nil, errors.New("unknown noti_format_id")
	}

	sort.Slice(times, func(i, j int) bool { return times[i].Before(times[j]) })
	return times, nil
}

// ===================================================================
//            Generate: สร้าง NotiItem จาก NotiInfo ในช่วงวัน
// ===================================================================

func GenerateNotiItemsForNotiInfoRange(db *gorm.DB, notiInfoID uint, fromDate, toDate time.Time) ([]models.NotiItem, error) {
	// โหลด NotiInfo + relation พื้นฐาน
	var notiInfo models.NotiInfo
	if err := db.Preload("MyMedicine").
		Preload("Group").
		First(&notiInfo, notiInfoID).Error; err != nil {
		return nil, err
	}

	occurrences, err := ExpandNotiInfoOccurrences(&notiInfo, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	if len(occurrences) == 0 {
		return []models.NotiItem{}, nil
	}

	// เตรียม target sources
	type medSource struct {
		PatientID     uint
		MyMedicineID  uint
		GroupID       *uint
		MedName       string
		GroupName     string
		AmountPerTime string
		FormID        uint
		UnitID        *uint
		InstructionID *uint
	}
	var sources []medSource

	if notiInfo.MyMedicineID != nil {
		m := notiInfo.MyMedicine
		sources = append(sources, medSource{
			PatientID:     m.PatientID,
			MyMedicineID:  m.ID,
			GroupID:       nil,
			MedName:       m.MedName,
			GroupName:     "",
			AmountPerTime: m.AmountPerTime,
			FormID:        m.FormID,
			UnitID:        m.UnitID,
			InstructionID: m.InstructionID,
		})
	} else if notiInfo.GroupID != nil {
		var members []models.MyMedicine
		if err := db.Where("patient_id = ? AND group_id = ?", notiInfo.Group.PatientID, *notiInfo.GroupID).
			Find(&members).Error; err != nil {
			return nil, err
		}
		for _, m := range members {
			sources = append(sources, medSource{
				PatientID:     m.PatientID,
				MyMedicineID:  m.ID,
				GroupID:       notiInfo.GroupID,
				MedName:       m.MedName,
				GroupName:     notiInfo.Group.GroupName,
				AmountPerTime: m.AmountPerTime,
				FormID:        m.FormID,
				UnitID:        m.UnitID,
				InstructionID: m.InstructionID,
			})
		}
	} else {
		return nil, gorm.ErrInvalidData
	}

	// สร้าง NotiItem (upsert กันซ้ำ)
	var newItems []models.NotiItem
	for _, src := range sources {
		for _, ts := range occurrences {
			notifyDate := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.Local)
			notifyTime := time.Date(1, 1, 1, ts.Hour(), ts.Minute(), ts.Second(), 0, time.Local) // type: time-only

			item := models.NotiItem{
				PatientID:    src.PatientID,
				MyMedicineID: src.MyMedicineID,
				GroupID:      nil, // set ถ้ากรุ๊ป
				NotiInfoID:   notiInfo.ID,

				MedName:       src.MedName,
				GroupName:     src.GroupName,
				AmountPerTime: src.AmountPerTime,

				FormID:        src.FormID,
				UnitID:        nil,
				InstructionID: nil,

				NotifyTime:   notifyTime,
				NotifyDate:   notifyDate,
				TakenStatus:  false,
				TakenTimeAt:  nil,
				NotifyStatus: false,
			}
			if src.GroupID != nil {
				item.GroupID = src.GroupID
			}
			if src.UnitID != nil {
				item.UnitID = src.UnitID
			}
			if src.InstructionID != nil {
				item.InstructionID = src.InstructionID
			}
			newItems = append(newItems, item)
		}
	}

	if len(newItems) == 0 {
		return []models.NotiItem{}, nil
	}
	// ต้องมี unique index ที่ noti_items: (my_medicine_id, notify_date, notify_time, noti_info_id)
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "my_medicine_id"}, {Name: "notify_date"}, {Name: "notify_time"}, {Name: "noti_info_id"}},
		DoNothing: true,
	}).Create(&newItems).Error; err != nil {
		return nil, err
	}

	return newItems, nil
}

// เติมล่วงหน้า N วันจากวันนี้ (ใช้ใน cron/job ได้)
func GenerateNotiItemsDaysAhead(db *gorm.DB, notiInfoID uint, days int) ([]models.NotiItem, error) {
	now := time.Now().In(time.Local)
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, days)
	return GenerateNotiItemsForNotiInfoRange(db, notiInfoID, from, to)
}

// ===================================================================
//                   Query/List: ดึงรายการ NotiItem
// ===================================================================

type ListNotiItemsFilter struct {
	PatientID    *uint
	MyMedicineID *uint
	GroupID      *uint
	NotiInfoID   *uint
	DateFrom     *string // "YYYY-MM-DD"
	DateTo       *string // "YYYY-MM-DD"
	TakenStatus  *bool
	NotifyStatus *bool
}

func ListNotiItems(db *gorm.DB, f ListNotiItemsFilter) ([]models.NotiItem, error) {
	q := db.Model(&models.NotiItem{}).
		Preload("Patient").
		Preload("MyMedicine").
		Preload("Group").
		Preload("NotiInfo").
		Preload("Form").
		Preload("Unit").
		Preload("Instruction")

	if f.PatientID != nil {
		q = q.Where("patient_id = ?", *f.PatientID)
	}
	if f.MyMedicineID != nil {
		q = q.Where("my_medicine_id = ?", *f.MyMedicineID)
	}
	if f.GroupID != nil {
		q = q.Where("group_id = ?", *f.GroupID)
	}
	if f.NotiInfoID != nil {
		q = q.Where("noti_info_id = ?", *f.NotiInfoID)
	}
	if f.TakenStatus != nil {
		q = q.Where("taken_status = ?", *f.TakenStatus)
	}
	if f.NotifyStatus != nil {
		q = q.Where("notify_status = ?", *f.NotifyStatus)
	}
	if f.DateFrom != nil && *f.DateFrom != "" {
		if d, err := time.ParseInLocation(ymdLayout, *f.DateFrom, time.Local); err == nil {
			q = q.Where("notify_date >= ?", d)
		}
	}
	if f.DateTo != nil && *f.DateTo != "" {
		if d, err := time.ParseInLocation(ymdLayout, *f.DateTo, time.Local); err == nil {
			q = q.Where("notify_date <= ?", d)
		}
	}

	var items []models.NotiItem
	if err := q.Order("notify_date, notify_time").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ===================================================================
//                Update: เปลี่ยนสถานะ taken/notify
// ===================================================================

// MarkNotiItemTaken: เซ็ต/ยกเลิกสถานะ “ทานแล้ว”
func MarkNotiItemTaken(db *gorm.DB, notiItemID uint, taken bool) (*models.NotiItem, error) {
	var item models.NotiItem
	if err := db.First(&item, notiItemID).Error; err != nil {
		return nil, err
	}

	updateFields := map[string]any{
		"taken_status": taken,
	}
	if taken {
		now := time.Now().In(time.Local)
		updateFields["taken_time_at"] = &now
	} else {
		// ถ้าอยากล้างเวลา เมื่อยกเลิก ให้ตั้งค่าเป็น NULL ได้เพราะเป็น *time.Time
		updateFields["taken_time_at"] = nil
	}

	if err := db.Model(&item).Updates(updateFields).Error; err != nil {
		return nil, err
	}
	if err := db.First(&item, notiItemID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// MarkNotiItemNotified: เซ็ต/ยกเลิก “แจ้งเตือนแล้ว”
func MarkNotiItemNotified(db *gorm.DB, notiItemID uint, notified bool) (*models.NotiItem, error) {
	var item models.NotiItem
	if err := db.First(&item, notiItemID).Error; err != nil {
		return nil, err
	}
	if err := db.Model(&item).Update("notify_status", notified).Error; err != nil {
		return nil, err
	}
	if err := db.First(&item, notiItemID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
