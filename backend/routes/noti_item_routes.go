package routes

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/dto" // ใช้ DTO ตัดข้อมูลส่วนเกิน
)

func SetupNotiItemsRoutes(api fiber.Router) {

	// LIST — ดึงรายการ NotiItem (รองรับฟิลเตอร์)
	// GET /api/noti-items?patient_id=&my_medicine_id=&group_id=&noti_info_id=&date_from=YYYY-MM-DD&date_to=YYYY-MM-DD&taken_status=true|false&notify_status=true|false
	// ตัวอย่าง:
	//   /api/noti-items?patient_id=7&date_from=2025-10-01&date_to=2025-10-07
	//   /api/noti-items?my_medicine_id=12&taken_status=true
	//   /api/noti-items?group_id=3&noti_info_id=44&notify_status=false
	api.Get("/noti-items", func(ctx *fiber.Ctx) error {
		// ⛑️ auth: ดึง patient_id จาก Locals เสมอ
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

		// ✅ บังคับใช้ patientID จาก Locals ที่ระดับ handlers
		notiItems, err := handlers.ListNotiItems(db.DB, patientID, filter)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// ส่งเฉพาะฟิลด์ที่ต้องการด้วย DTO (ตัด relations/ฟิลด์ว่าง)
		return ctx.JSON(fiber.Map{"data": dto.NotiItemsToDTO(notiItems)})
	})

	// UPDATE — Mark Taken (เซ็ต/ยกเลิก “ทานแล้ว”)
	// PATCH /api/noti-items/:id/taken
	// Body (ตัวอย่าง):
	// {
	//   "taken": true
	// }
	api.Patch("/noti-items/:id/taken", func(ctx *fiber.Ctx) error {
		// ⛑️ auth
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
		// ⛑️ auth
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
		// ⛑️ auth
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
		// ⛑️ auth
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
