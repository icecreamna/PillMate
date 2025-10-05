package routes

import (
	"errors"
	"fmt"      
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/dto" // ใช้ DTO ตัดข้อมูลส่วนเกิน
	"github.com/fouradithep/pillmate/models"
)

func SetupNotiItemsRoutes(api fiber.Router) {

	// LIST — ดึงรายการ NotiItem (รองรับฟิลเตอร์)
	// GET /api/noti-items?patient_id=&my_medicine_id=&group_id=&noti_info_id=&date_from=YYYY-MM-DD&date_to=YYYY-MM-DD&taken_status=true|false&notify_status=true|false
	// ตัวอย่าง:
	//   /api/noti-items?patient_id=7&date_from=2025-10-01&date_to=2025-10-07
	//   /api/noti-items?my_medicine_id=12&taken_status=true
	//   /api/noti-items?group_id=3&noti_info_id=44&notify_status=false
	api.Get("/noti-items", func(ctx *fiber.Ctx) error {
		// auth: ดึง patient_id จาก Locals เสมอ
		patientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var filter handlers.ListNotiItemsFilter

		// หมายเหตุ: handlers.ListNotiItems จะบังคับใช้ patientID จาก Locals เสมอ
		// การส่ง ?patient_id=... มาจะไม่ override สิทธิ์ผู้ใช้
		if patientIDStr := ctx.Query("patient_id"); patientIDStr != "" {
			if parsed, err := strconv.ParseUint(patientIDStr, 10, 64); err == nil {
				// เก็บลงฟิลด์ได้ แต่ไม่ถูกใช้แทน patientID จาก Locals
				tmp := uint(parsed)
				filter.PatientID = &tmp
			}
		}
		if myMedIDStr := ctx.Query("my_medicine_id"); myMedIDStr != "" {
			if parsed, err := strconv.ParseUint(myMedIDStr, 10, 64); err == nil {
				myMedicineID := uint(parsed)
				filter.MyMedicineID = &myMedicineID
			}
		}
		if groupIDStr := ctx.Query("group_id"); groupIDStr != "" {
			if parsed, err := strconv.ParseUint(groupIDStr, 10, 64); err == nil {
				groupID := uint(parsed)
				filter.GroupID = &groupID
			}
		}
		if notiInfoIDStr := ctx.Query("noti_info_id"); notiInfoIDStr != "" {
			if parsed, err := strconv.ParseUint(notiInfoIDStr, 10, 64); err == nil {
				notiInfoID := uint(parsed)
				filter.NotiInfoID = &notiInfoID
			}
		}
		if fromStr := strings.TrimSpace(ctx.Query("date_from")); fromStr != "" {
			filter.DateFrom = &fromStr // รูปแบบ "YYYY-MM-DD"
		}
		if toStr := strings.TrimSpace(ctx.Query("date_to")); toStr != "" {
			filter.DateTo = &toStr // รูปแบบ "YYYY-MM-DD"
		}
		if takenStatusStr := ctx.Query("taken_status"); takenStatusStr != "" {
			taken := strings.EqualFold(takenStatusStr, "true") || takenStatusStr == "1"
			filter.TakenStatus = &taken
		}
		if notifyStatusStr := ctx.Query("notify_status"); notifyStatusStr != "" {
			notified := strings.EqualFold(notifyStatusStr, "true") || notifyStatusStr == "1"
			filter.NotifyStatus = &notified
		}

		// บังคับใช้ patientID จาก Locals ที่ระดับ handlers
		notiItems, err := handlers.ListNotiItems(db.DB, patientID, filter)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// ====== ดึง symptom ของ items เหล่านี้ครั้งเดียว แล้วทำ map: noti_item_id -> symptom_id
		ids := make([]uint, 0, len(notiItems))
		for _, it := range notiItems {
			ids = append(ids, it.ID)
		}
		symMap := map[uint]uint{}
		if len(ids) > 0 {
			var rows []struct {
				ID         uint
				NotiItemID uint
			}
			if err := db.DB.Model(&models.Symptom{}).
				Where("patient_id = ? AND noti_item_id IN ?", patientID, ids).
				Select("id, noti_item_id").
				Scan(&rows).Error; err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
			}
			for _, r := range rows {
				symMap[r.NotiItemID] = r.ID
			}
		}

		// ส่งเฉพาะฟิลด์ที่ต้องการด้วย DTO (ตัด relations/ฟิลด์ว่าง) + ผนวกสถานะ symptom (ระดับ item)
		dtoList := dto.NotiItemsToDTOWithSymptoms(notiItems, symMap)

		// ===== สร้าง 2 ก้อน: data (เดี่ยว) + group_cards (กลุ่ม) =====
		// ก้อน data: เอาเฉพาะที่ "ไม่มีกลุ่ม"
		flats := make([]dto.NotiItemDTO, 0, len(dtoList))

		// ก้อน group_cards: รวมด้วย (group_id, notify_date, notify_time)
		type GroupLine struct {
			NotiItemID    uint   `json:"noti_item_id"`
			MyMedicineID  uint   `json:"my_medicine_id"`
			MedName       string `json:"med_name"`
			AmountPerTime string `json:"amount_per_time"`
			FormID        uint   `json:"form_id"`
			UnitID        *uint  `json:"unit_id,omitempty"`
			InstructionID *uint  `json:"instruction_id,omitempty"`
		}
		type GroupCard struct {
			GroupID      uint        `json:"group_id"`
			GroupName    string      `json:"group_name,omitempty"`
			PatientID    uint        `json:"patient_id"`
			NotifyDate   string      `json:"notify_date"`   // "YYYY-MM-DD"
			NotifyTime   string      `json:"notify_time"`   // "HH:MM"
			TakenStatus  bool        `json:"taken_status"`  // true = ทุกชิ้น taken
			NotifyStatus bool        `json:"notify_status"` // true = ทุกชิ้น notified
			NotiInfoID   uint        `json:"noti_info_id"`

			// สรุปอาการ “ระดับการ์ดกลุ่ม” (ไม่อยู่ใน items)
			HasSymptom bool  `json:"has_symptom"`
			SymptomID  *uint `json:"symptom_id,omitempty"`

			Items []GroupLine `json:"items"`
		}
		groupMap := map[string]*GroupCard{}

		for i, raw := range notiItems {
			d := dtoList[i] // DTO ที่ map แล้ว (มี has_symptom / symptom_id สำหรับสรุปที่ระดับการ์ด)

			if raw.GroupID == nil {
				// ไม่มี group -> ใส่ก้อน data
				flats = append(flats, d)
				continue
			}

			key := fmt.Sprintf("%d|%s|%s", *raw.GroupID, d.NotifyDate, d.NotifyTime)
			card, ok := groupMap[key]
			if !ok {
				card = &GroupCard{
					GroupID:      *raw.GroupID,
					GroupName:    raw.GroupName,
					PatientID:    raw.PatientID,
					NotifyDate:   d.NotifyDate,
					NotifyTime:   d.NotifyTime,
					TakenStatus:  true, // เริ่มจาก true แล้ว AND ลงไป
					NotifyStatus: true, // เริ่มจาก true แล้ว AND ลงไป
					NotiInfoID:   d.NotiInfoID,
					HasSymptom:   false,
					SymptomID:    nil,
					Items:        []GroupLine{},
				}
				groupMap[key] = card
			}

			card.Items = append(card.Items, GroupLine{
				NotiItemID:    d.ID,
				MyMedicineID:  d.MyMedicineID,
				MedName:       d.MedName,
				AmountPerTime: d.AmountPerTime,
				FormID:        d.FormID,
				UnitID:        d.UnitID,
				InstructionID: d.InstructionID,
			})

			// สถานะรวมการ์ด: true ก็ต่อเมื่อ "ทุกชิ้น" เป็น true
			card.TakenStatus = card.TakenStatus && d.TakenStatus
			card.NotifyStatus = card.NotifyStatus && d.NotifyStatus

			// สรุปอาการไว้ที่ระดับการ์ด: ถ้ามี item ไหนมีอาการ ก็ถือว่าการ์ดมีอาการ
			if d.HasSymptom && !card.HasSymptom {
				card.HasSymptom = true
				card.SymptomID = d.SymptomID
			}
		}

		// map -> slice
		groupCards := make([]GroupCard, 0, len(groupMap))
		for _, c := range groupMap {
			groupCards = append(groupCards, *c)
		}

		// ส่งสองก้อนกลับไป
		return ctx.JSON(fiber.Map{
			"data":        flats,      // เฉพาะรายการที่ไม่มีกลุ่ม
			"group_cards": groupCards, // รวมกลุ่มตาม group_id + date + time (พร้อมสรุปอาการระดับการ์ด)
		})
	})


	// UPDATE — Mark Taken (เซ็ต/ยกเลิก “ทานแล้ว”)
	// PATCH /api/noti-items/:id/taken
	// Body (ตัวอย่าง):
	// {
	//   "taken": true
	// }
	api.Patch("/noti-items/:id/taken", func(ctx *fiber.Ctx) error {
		// auth
		patientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		parsedID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil || parsedID == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		notiItemID := uint(parsedID)

		type MarkTakenRequest struct{ Taken *bool `json:"taken"` }
		var req MarkTakenRequest
		if err := ctx.BodyParser(&req); err != nil || req.Taken == nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body (taken required)"})
		}

		updated, uerr := handlers.MarkNotiItemTaken(db.DB, patientID, notiItemID, *req.Taken)
		if uerr != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(uerr, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			return ctx.Status(status).JSON(fiber.Map{"error": uerr.Error()})
		}
		// ตอบกลับแบบ DTO ให้ payload สะอาด
		return ctx.JSON(fiber.Map{"data": dto.NotiItemToDTO(*updated)})
	})

	// UPDATE — Mark Notified (เซ็ต/ยกเลิก “แจ้งเตือนแล้ว”)
	// PATCH /api/noti-items/:id/notified
	// Body (ตัวอย่าง):
	// {
	//   "notified": true
	// }
	api.Patch("/noti-items/:id/notified", func(ctx *fiber.Ctx) error {
		// auth
		patientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		parsedID, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil || parsedID == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		notiItemID := uint(parsedID)

		type MarkNotifiedRequest struct{ Notified *bool `json:"notified"` }
		var req MarkNotifiedRequest
		if err := ctx.BodyParser(&req); err != nil || req.Notified == nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body (notified required)"})
		}

		updated, uerr := handlers.MarkNotiItemNotified(db.DB, patientID, notiItemID, *req.Notified)
		if uerr != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(uerr, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			return ctx.Status(status).JSON(fiber.Map{"error": uerr.Error()})
		}
		// ตอบกลับแบบ DTO
		return ctx.JSON(fiber.Map{"data": dto.NotiItemToDTO(*updated)})
	})

	// GENERATE — สร้างรายการล่วงหน้าตามช่วงวัน "ของผู้ป่วยทั้งคน" (กำหนด from/to เอง)
	// POST /api/noti-items/generate-range
	// Body (ตัวอย่าง):
	// {
	//   "from_date": "2025-10-01",
	//   "to_date": "2025-10-15"
	// }
	api.Post("/noti-items/generate-range", func(ctx *fiber.Ctx) error {
		// auth
		patientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		type GenerateRangeRequest struct {
			FromDate string `json:"from_date"` // "YYYY-MM-DD"
			ToDate   string `json:"to_date"`   // "YYYY-MM-DD"
		}
		var req GenerateRangeRequest
		if err := ctx.BodyParser(&req); err != nil || req.FromDate == "" || req.ToDate == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		from, errFrom := time.ParseInLocation("2006-01-02", req.FromDate, time.Local)
		to, errTo := time.ParseInLocation("2006-01-02", req.ToDate, time.Local)
		if errFrom != nil || errTo != nil || to.Before(from) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid date range"})
		}

		created, genErr := handlers.GenerateNotiItemsForPatientRange(db.DB, patientID, from, to)
		if genErr != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": genErr.Error()})
		}
		// ตอบกลับแบบ DTO
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": dto.NotiItemsToDTO(created)})
	})

	// GENERATE — เติมล่วงหน้า N วันนับจากวันนี้ (rolling window)
	// POST /api/noti-items/generate-days-ahead
	// Body (ตัวอย่าง):
	// {
	//   "days": 14
	// }
	api.Post("/noti-items/generate-days-ahead", func(ctx *fiber.Ctx) error {
		// auth
		patientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		type GenerateDaysAheadRequest struct {
			DaysAhead int `json:"days"`
		}
		var req GenerateDaysAheadRequest
		if err := ctx.BodyParser(&req); err != nil || req.DaysAhead <= 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		created, genErr := handlers.GenerateNotiItemsDaysAheadForPatient(db.DB, patientID, req.DaysAhead)
		if genErr != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": genErr.Error()})
		}
		// ตอบกลับแบบ DTO
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": dto.NotiItemsToDTO(created)})
	})

}
