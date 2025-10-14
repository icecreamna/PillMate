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

// เลือก NotiItem ที่ถึงเวลาในหน้าต่าง N นาทีจาก "ตอนนี้"
func GetDueNow(db *gorm.DB, patientID uint, windowMinutes int) ([]models.NotiItem, error) {
	if windowMinutes <= 0 {
		windowMinutes = 1
	}

	// ✅ ใช้เวลา UTC เสมอ เพื่อเทียบกับ notify_time ที่เก็บเป็น UTC
	nowUTC := time.Now().UTC()
	today := nowUTC.Truncate(24 * time.Hour) // วันที่ UTC

	startT := time.Date(1, 1, 1, nowUTC.Hour(), nowUTC.Minute(), 0, 0, time.UTC)
	endUTC := nowUTC.Add(time.Duration(windowMinutes) * time.Minute)
	endT := time.Date(1, 1, 1, endUTC.Hour(), endUTC.Minute(), 59, 0, time.UTC)

	var out []models.NotiItem
	err := db.Model(&models.NotiItem{}).
		Where("patient_id = ? AND notify_status = FALSE", patientID).
		Where("notify_date = ?", today).
		Where("notify_time >= ? AND notify_time <= ?", startT, endT).
		Order("notify_date, notify_time, id").
		Find(&out).Error

	return out, err
}
