package handlers

import (
	"time"
	"gorm.io/gorm"
	"github.com/fouradithep/pillmate/models"
)

var appLoc = func() *time.Location {
	l, _ := time.LoadLocation("Asia/Bangkok")
	return l
}()

func now() time.Time { return time.Now().In(appLoc) }

// เลือก NotiItem ที่ถึงเวลาในหน้าต่าง N นาทีจาก "ตอนนี้"
func GetDueNow(db *gorm.DB, patientID uint, windowMinutes int) ([]models.NotiItem, error) {
    if windowMinutes <= 0 { windowMinutes = 1 }

    // ใช้เวลา Local = Asia/Bangkok (คุณตั้งไว้แล้ว)
    nowLocal := time.Now().In(appLoc)
    today := time.Date(nowLocal.Year(), nowLocal.Month(), nowLocal.Day(), 0, 0, 0, 0, appLoc)

    // สร้าง "time-only" ด้วยปี 0001, โซน UTC ให้ตรงกับตอนบันทึก NotifyTime
    startT := time.Date(1, 1, 1, nowLocal.Hour(), nowLocal.Minute(), 0, 0, time.UTC)

    endLocal := nowLocal.Add(time.Duration(windowMinutes) * time.Minute)
    endT := time.Date(1, 1, 1, endLocal.Hour(), endLocal.Minute(), 59, 0, time.UTC)
    // จะใช้วินาที 59 หรือใช้วินาที 0 แล้ว WHERE <= ก็ได้ ทั้งสองทางเท่ากันถ้า NotifyTime เป็นนาทีเป๊ะ

    var out []models.NotiItem
    err := db.Model(&models.NotiItem{}).
        Where("patient_id = ? AND notify_status = FALSE", patientID).
        Where("notify_date = ?", today).
        Where("notify_time >= ? AND notify_time <= ?", startT, endT). // รวมปลายช่วง
        Order("notify_date, notify_time, id").
        Find(&out).Error
    return out, err
}

