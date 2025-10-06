package handlers

import (
	"errors"
	"sort"
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ===================================================================
//                     Expand: แตกเวลาจาก NotiInfo
// ===================================================================

// แยกชั่วโมง/นาทีจาก "HH:MM"
func splitHourMinute(hhmm string) (int, int, error) {
	if len(hhmm) != 5 || hhmm[2] != ':' {
		return 0, 0, errors.New("invalid HH:MM")
	}
	hour := int(hhmm[0]-'0')*10 + int(hhmm[1]-'0')
	minute := int(hhmm[3]-'0')*10 + int(hhmm[4]-'0')
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, errors.New("invalid HH:MM range")
	}
	return hour, minute, nil
}

// สร้างเวลาจากวันที่ + "HH:MM" โดยยึด time.Local
func makeLocalDateTime(baseDate time.Time, hhmm string) (time.Time, error) {
	year, month, day := baseDate.In(time.Local).Date()
	hour, minute, err := splitHourMinute(hhmm)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, month, day, hour, minute, 0, 0, time.Local), nil
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

	ruleStartDate := time.Date(notiInfo.StartDate.Year(), notiInfo.StartDate.Month(), notiInfo.StartDate.Day(), 0, 0, 0, 0, time.Local)
	ruleEndDate := time.Date(notiInfo.EndDate.Year(), notiInfo.EndDate.Month(), notiInfo.EndDate.Day(), 23, 59, 59, 0, time.Local)

	windowStartDate := time.Date(winStart.Year(), winStart.Month(), winStart.Day(), 0, 0, 0, 0, time.Local)
	windowEndDate := time.Date(winEnd.Year(), winEnd.Month(), winEnd.Day(), 23, 59, 59, 0, time.Local)

	intersectStart, intersectEnd, hasOverlap := intersectRange(ruleStartDate, ruleEndDate, windowStartDate, windowEndDate)
	if !hasOverlap {
		return nil, nil
	}

	var occurTimes []time.Time
	appendOccur := func(t time.Time) { occurTimes = append(occurTimes, t) }

	switch notiInfo.NotiFormatID {
	case NotiFormatFixedTimes:
		if notiInfo.Times == nil || len(*notiInfo.Times) == 0 {
			return nil, errors.New("times required")
		}
		for cursorDate := intersectStart; !cursorDate.After(intersectEnd); cursorDate = cursorDate.AddDate(0, 0, 1) {
			for _, hhmm := range *notiInfo.Times {
				ts, err := makeLocalDateTime(cursorDate, hhmm)
				if err != nil {
					return nil, err
				}
				appendOccur(ts)
			}
		}

	case NotiFormatEveryNDays:
		if notiInfo.IntervalDay == nil || *notiInfo.IntervalDay <= 0 {
			return nil, errors.New("interval_day required")
		}
		if notiInfo.Times == nil || len(*notiInfo.Times) == 0 {
			return nil, errors.New("times required")
		}
		for ruleCursor := ruleStartDate; !ruleCursor.After(intersectEnd); ruleCursor = ruleCursor.AddDate(0, 0, *notiInfo.IntervalDay) {
			if ruleCursor.Before(intersectStart) {
				continue
			}
			for _, hhmm := range *notiInfo.Times {
				ts, err := makeLocalDateTime(ruleCursor, hhmm)
				if err != nil {
					return nil, err
				}
				appendOccur(ts)
			}
		}

	case NotiFormatInterval:
		if notiInfo.IntervalHours == nil || *notiInfo.IntervalHours <= 0 {
			return nil, errors.New("interval_hours required")
		}
		stepDuration := time.Duration(*notiInfo.IntervalHours) * time.Hour
		cursor := intersectStart

		var perDayLimit int
		if notiInfo.TimesPerDay != nil && *notiInfo.TimesPerDay > 0 {
			perDayLimit = *notiInfo.TimesPerDay
		}
		todayCount := 0
		currentDay := time.Date(cursor.Year(), cursor.Month(), cursor.Day(), 0, 0, 0, 0, time.Local)

		for !cursor.After(intersectEnd) {
			cursorDay := time.Date(cursor.Year(), cursor.Month(), cursor.Day(), 0, 0, 0, 0, time.Local)
			if cursorDay.After(currentDay) {
				currentDay = cursorDay
				todayCount = 0
			}
			if cursor.Before(intersectStart) {
				cursor = cursor.Add(stepDuration)
				continue
			}
			if perDayLimit == 0 || todayCount < perDayLimit {
				appendOccur(cursor)
				todayCount++
			}
			cursor = cursor.Add(stepDuration)
		}

	case NotiFormatCycle:
		if notiInfo.CyclePattern == nil || len(*notiInfo.CyclePattern) == 0 {
			return nil, errors.New("cycle_pattern required")
		}
		if notiInfo.Times == nil || len(*notiInfo.Times) == 0 {
			return nil, errors.New("times required")
		}
		pattern := *notiInfo.CyclePattern // ตัวอย่าง [21,7] = on/off
		ruleCursor := ruleStartDate
		patternIndex := 0
		for !ruleCursor.After(intersectEnd) {
			// on days
			onDays := int(pattern[patternIndex%len(pattern)])
			for i := 0; i < onDays && !ruleCursor.After(intersectEnd); i++ {
				if !ruleCursor.Before(intersectStart) {
					for _, hhmm := range *notiInfo.Times {
						ts, err := makeLocalDateTime(ruleCursor, hhmm)
						if err != nil {
							return nil, err
						}
						appendOccur(ts)
					}
				}
				ruleCursor = ruleCursor.AddDate(0, 0, 1)
			}
			patternIndex++
			// off days
			if patternIndex < len(pattern) {
				offDays := int(pattern[patternIndex%len(pattern)])
				ruleCursor = ruleCursor.AddDate(0, 0, offDays)
				patternIndex++
			} else if len(pattern) >= 2 {
				// วนซ้ำค่า off ล่าสุด
				ruleCursor = ruleCursor.AddDate(0, 0, int(pattern[1]))
			}
		}

	default:
		return nil, errors.New("unknown noti_format_id")
	}

	sort.Slice(occurTimes, func(i, j int) bool { return occurTimes[i].Before(occurTimes[j]) })
	return occurTimes, nil
}

// ===================================================================
//            Generate: สร้าง NotiItem จาก NotiInfo ในช่วงวัน
// ===================================================================

// จำกัดสิทธิ์ด้วย patientID (ต้องเป็น NotiInfo ของผู้ป่วยคนนั้นเท่านั้น)
func GenerateNotiItemsForNotiInfoRange(db *gorm.DB, patientID, notiInfoID uint, fromDate, toDate time.Time) ([]models.NotiItem, error) {
	// โหลด NotiInfo + relation พื้นฐาน (จำกัดด้วย EXISTS ว่าเป็นของ patient นี้จริง)
	var notiInfo models.NotiInfo
	if err := db.
		Preload("MyMedicine").
		Preload("Group").
		Where("id = ?", notiInfoID).
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
		First(&notiInfo).Error; err != nil {
		return nil, err
	}

	occurTimes, err := ExpandNotiInfoOccurrences(&notiInfo, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	if len(occurTimes) == 0 {
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
	var medSources []medSource

	if notiInfo.MyMedicineID != nil {
		mm := notiInfo.MyMedicine
		medSources = append(medSources, medSource{
			PatientID:     mm.PatientID,
			MyMedicineID:  mm.ID,
			GroupID:       nil,
			MedName:       mm.MedName,
			GroupName:     "",
			AmountPerTime: mm.AmountPerTime,
			FormID:        mm.FormID,
			UnitID:        mm.UnitID,
			InstructionID: mm.InstructionID,
		})
	} else if notiInfo.GroupID != nil {
		// ดึงสมาชิกกลุ่มของผู้ป่วยคนนี้เท่านั้น
		var groupMembers []models.MyMedicine
		if err := db.
			Where("patient_id = ? AND group_id = ?", patientID, *notiInfo.GroupID).
			Find(&groupMembers).Error; err != nil {
			return nil, err
		}
		for _, gm := range groupMembers {
			medSources = append(medSources, medSource{
				PatientID:     gm.PatientID,
				MyMedicineID:  gm.ID,
				GroupID:       notiInfo.GroupID,
				MedName:       gm.MedName,
				GroupName:     notiInfo.Group.GroupName,
				AmountPerTime: gm.AmountPerTime,
				FormID:        gm.FormID,
				UnitID:        gm.UnitID,
				InstructionID: gm.InstructionID,
			})
		}
	} else {
		return nil, gorm.ErrInvalidData
	}

	// สร้าง NotiItem (upsert กันซ้ำ)
	var newItems []models.NotiItem
	for _, src := range medSources {
		for _, ts := range occurTimes {
			notifyDateOnly := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, time.Local)
			notifyTimeOnly := time.Date(1, 1, 1, ts.Hour(), ts.Minute(), ts.Second(), 0, time.UTC) // type: time-only

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

				NotifyTime:   notifyTimeOnly,
				NotifyDate:   notifyDateOnly,
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
	if err := db.Clauses(clause.OnConflict{
		DoNothing: true, // ใช้ unique index noti_items_uniq_active ที่ DB สร้างไว้
	}).Create(&newItems).Error; err != nil {
		return nil, err
	}

	return newItems, nil
}

// เติมล่วงหน้า N วันจากวันนี้ (ใช้ใน cron/job ได้)
func GenerateNotiItemsDaysAhead(db *gorm.DB, patientID, notiInfoID uint, days int) ([]models.NotiItem, error) {
	nowLocal := time.Now().In(time.Local)
	todayStart := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, time.Local)
	dateTo := todayStart.AddDate(0, 0, days)
	return GenerateNotiItemsForNotiInfoRange(db, patientID, notiInfoID, todayStart, dateTo)
}

// GenerateNotiItemsDaysAheadForPatient: เติมล่วงหน้า N วันสำหรับ noti_infos ทั้งหมดของผู้ป่วย
func GenerateNotiItemsDaysAheadForPatient(db *gorm.DB, patientID uint, days int) ([]models.NotiItem, error) {
	if days <= 0 {
		return []models.NotiItem{}, nil
	}

	// ดึง noti_info.id ทั้งหมดที่เป็นของผู้ป่วยคนนี้ (ผ่าน EXISTS ไปยัง my_medicines / groups)
	var notiInfoIDs []uint
	if err := db.Model(&models.NotiInfo{}).
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
		Pluck("id", &notiInfoIDs).Error; err != nil {
		return nil, err
	}

	if len(notiInfoIDs) == 0 {
		return []models.NotiItem{}, nil
	}

	// loop เรียกของเดิมทีละ noti_info (ตัวเดิมมี on-conflict do nothing กันซ้ำแล้ว)
	allCreated := make([]models.NotiItem, 0, 1024)
	for _, id := range notiInfoIDs {
		created, err := GenerateNotiItemsDaysAhead(db, patientID, id, days)
		if err != nil {
			return nil, err
		}
		allCreated = append(allCreated, created...)
	}
	return allCreated, nil
}

// GenerateNotiItemsForPatientRange: สร้าง NotiItem ให้ noti_infos "ทั้งหมด" ของผู้ป่วยภายในช่วงวัน [fromDate, toDate]
func GenerateNotiItemsForPatientRange(db *gorm.DB, patientID uint, fromDate, toDate time.Time) ([]models.NotiItem, error) {
	if toDate.Before(fromDate) {
		return []models.NotiItem{}, nil
	}

	// ดึง noti_info.id ทั้งหมดของผู้ป่วย (ผ่าน EXISTS ไปยัง my_medicines / groups)
	var notiInfoIDs []uint
	if err := db.Model(&models.NotiInfo{}).
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
		Pluck("id", &notiInfoIDs).Error; err != nil {
		return nil, err
	}

	if len(notiInfoIDs) == 0 {
		return []models.NotiItem{}, nil
	}

	allCreated := make([]models.NotiItem, 0, 1024)
	for _, nid := range notiInfoIDs {
		created, err := GenerateNotiItemsForNotiInfoRange(db, patientID, nid, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		allCreated = append(allCreated, created...)
	}
	return allCreated, nil
}
