package handlers

import (
	"time"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)

var appLoc = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Bangkok")
	return l
}()

func now() time.Time { return time.Now().In(appLoc) }

// เลือก NotiItem ที่ถึงเวลาในหน้าต่าง N นาทีจาก "ตอนนี้" (อิงวันตาม Asia/Bangkok)
// - NotifyDate: ชนิด DATE (ไม่มี TZ) → เทียบด้วย "YYYY-MM-DD"
// - NotifyTime: ชนิด TIME (ไม่มี TZ) → ใช้ time-only (anchor 0001-01-01, UTC) ในเงื่อนไข
// - รองรับหน้าต่างเวลาที่ "ข้ามเที่ยงคืน" โดยแบ่งคิวรีเป็น 2 ช่วง (วันนี้ + พรุ่งนี้)
func GetDueNow(db *gorm.DB, patientID uint, windowMinutes int) ([]models.NotiItem, error) {
	if windowMinutes <= 0 {
		windowMinutes = 1
	}

	// เวลาท้องถิ่น (ไทย) เพื่อให้นิยาม "วันนี้" ตรงกับการบันทึก NotifyDate
	nowLocal := now()
	todayLocal := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, appLoc)
	todayStr := todayLocal.Format("2006-01-02") // สำหรับเทียบกับคอลัมน์ DATE

	// time-only ในโซน UTC (เพราะคอลัมน์ TIME ไม่สน TZ) แต่ชั่วโมง/นาทีตาม local
	startT := time.Date(1, 1, 1, nowLocal.Hour(), nowLocal.Minute(), 0, 0, time.UTC)
	endLocal := nowLocal.Add(time.Duration(windowMinutes) * time.Minute)
	endT := time.Date(1, 1, 1, endLocal.Hour(), endLocal.Minute(), 59, 0, time.UTC)

	var out []models.NotiItem

	// เบสคิวรีร่วม
	base := db.Model(&models.NotiItem{}).
		Where("patient_id = ? AND notify_status = FALSE", patientID).
		Order("notify_date, notify_time, id")

	// เคส "ไม่ข้ามเที่ยงคืน" (เช่น 02:07..02:12)
	if !endLocal.Before(nowLocal) && (endT.Equal(startT) || endT.After(startT)) {
		err := base.
			Where("notify_date = ?", todayStr).
			Where("notify_time BETWEEN ? AND ?", startT, endT).
			Find(&out).Error
		return out, err
	}

	// เคส "ข้ามเที่ยงคืน" → แบ่ง 2 ช่วง: วันนี้ [start..23:59:59] และ พรุ่งนี้ [00:00..end]
	var part1, part2 []models.NotiItem

	// วันนี้
	err1 := base.
		Where("notify_date = ?", todayStr).
		Where(
			"notify_time BETWEEN ? AND ?",
			time.Date(1, 1, 1, startT.Hour(), startT.Minute(), 0, 0, time.UTC),
			time.Date(1, 1, 1, 23, 59, 59, 0, time.UTC),
		).
		Find(&part1).Error
	if err1 != nil {
		return nil, err1
	}

	// พรุ่งนี้
	tomorrowStr := todayLocal.Add(24 * time.Hour).Format("2006-01-02")
	err2 := base.
		Where("notify_date = ?", tomorrowStr).
		Where(
			"notify_time BETWEEN ? AND ?",
			time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(1, 1, 1, endT.Hour(), endT.Minute(), endT.Second(), 0, time.UTC),
		).
		Find(&part2).Error
	if err2 != nil {
		return nil, err2
	}

	out = append(out, part1...)
	out = append(out, part2...)
	return out, nil
}
