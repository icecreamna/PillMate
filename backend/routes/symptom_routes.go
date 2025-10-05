package routes

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
)

func SetupSymptomRoutes(api fiber.Router) {

	// =========================================================
	// CREATE
	// POST /api/symptom
	// Body ตัวอย่าง:
	// {
	//   "noti_item_id": 123,
	//   "my_medicine_id": 45,    // optional (ถ้าไม่ส่ง handler จะเติมจาก noti_item ให้อัตโนมัติ)
	//   "group_id": 3,           // optional (ต้องสอดคล้องกับ noti_item ถ้ามี)
	//   "symptom_note": "เวียนหัว คลื่นไส้เล็กน้อย"
	// }
	// =========================================================
	api.Post("/symptom", func(ctx *fiber.Ctx) error {
		// auth
		authenticatedPatientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || authenticatedPatientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var requestBody models.Symptom
		if err := ctx.BodyParser(&requestBody); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		createdSymptom, err := handlers.CreateSymptom(db.DB, authenticatedPatientID, &requestBody)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrInvalidData) {
				status = fiber.StatusBadRequest
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			return ctx.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		// ส่งเฉพาะฟิลด์ของ Symptom (DTO)
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"data": dto.SymptomToDTO(*createdSymptom)})
	})

	// =========================================================
	// LIST (with filters)
	// GET /api/symptoms?my_medicine_id=&group_id=&noti_item_id=&created_from=YYYY-MM-DD&created_to=YYYY-MM-DD
	// ตัวอย่าง:
	//   /api/symptoms?created_from=2025-10-01&created_to=2025-10-07
	//   /api/symptoms?my_medicine_id=45
	//   /api/symptoms?group_id=3&noti_item_id=123
	// =========================================================
	api.Get("/symptoms", func(ctx *fiber.Ctx) error {
		// auth
		authenticatedPatientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || authenticatedPatientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var listFilter handlers.ListSymptomsFilter

		if myMedicineIDStr := ctx.Query("my_medicine_id"); myMedicineIDStr != "" {
			if parsed, err := strconv.ParseUint(myMedicineIDStr, 10, 64); err == nil {
				id := uint(parsed)
				listFilter.MyMedicineID = &id
			}
		}
		if groupIDStr := ctx.Query("group_id"); groupIDStr != "" {
			if parsed, err := strconv.ParseUint(groupIDStr, 10, 64); err == nil {
				id := uint(parsed)
				listFilter.GroupID = &id
			}
		}
		if notiItemIDStr := ctx.Query("noti_item_id"); notiItemIDStr != "" {
			if parsed, err := strconv.ParseUint(notiItemIDStr, 10, 64); err == nil {
				id := uint(parsed)
				listFilter.NotiItemID = &id
			}
		}
		if createdFromStr := strings.TrimSpace(ctx.Query("created_from")); createdFromStr != "" {
			listFilter.CreatedFrom = &createdFromStr
		}
		if createdToStr := strings.TrimSpace(ctx.Query("created_to")); createdToStr != "" {
			listFilter.CreatedTo = &createdToStr
		}

		symptomList, err := handlers.ListSymptoms(db.DB, authenticatedPatientID, listFilter)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// ส่งเฉพาะฟิลด์ของ Symptom (DTO)
		return ctx.JSON(fiber.Map{"data": dto.SymptomsToDTO(symptomList)})
	})

	// =========================================================
	// GET ONE
	// GET /api/symptom/:id
	// =========================================================
	api.Get("/symptom/:id", func(ctx *fiber.Ctx) error {
		// auth
		authenticatedPatientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || authenticatedPatientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		symptomIDUint64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil || symptomIDUint64 == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		symptomID := uint(symptomIDUint64)

		symptom, getErr := handlers.GetSymptom(db.DB, authenticatedPatientID, symptomID)
		if getErr != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(getErr, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			return ctx.Status(status).JSON(fiber.Map{"error": getErr.Error()})
		}
		// ส่งเฉพาะฟิลด์ของ Symptom (DTO)
		return ctx.JSON(fiber.Map{"data": dto.SymptomToDTO(*symptom)})
	})

	// =========================================================
	// UPDATE (note only)
	// PATCH /api/symptom/:id
	// Body ตัวอย่าง:
	// {
	//   "symptom_note": "เวียนหัวลดลงแล้ว แต่ยังง่วงนิดหน่อย"
	// }
	// =========================================================
	api.Patch("/symptom/:id", func(ctx *fiber.Ctx) error {
		// auth
		authenticatedPatientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || authenticatedPatientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		symptomIDUint64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil || symptomIDUint64 == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		symptomID := uint(symptomIDUint64)

		var requestBody struct {
			SymptomNote string `json:"symptom_note"`
		}
		if err := ctx.BodyParser(&requestBody); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		updatePayload := &models.Symptom{SymptomNote: strings.TrimSpace(requestBody.SymptomNote)}
		updatedSymptom, updateErr := handlers.UpdateSymptom(db.DB, authenticatedPatientID, symptomID, updatePayload)
		if updateErr != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(updateErr, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			if errors.Is(updateErr, gorm.ErrInvalidData) {
				status = fiber.StatusBadRequest
			}
			return ctx.Status(status).JSON(fiber.Map{"error": updateErr.Error()})
		}
		// ส่งเฉพาะฟิลด์ของ Symptom (DTO)
		return ctx.JSON(fiber.Map{"data": dto.SymptomToDTO(*updatedSymptom)})
	})

	// =========================================================
	// DELETE
	// DELETE /api/symptom/:id
	// =========================================================
	api.Delete("/symptom/:id", func(ctx *fiber.Ctx) error {
		// auth
		authenticatedPatientID, ok := ctx.Locals("patient_id").(uint)
		if !ok || authenticatedPatientID == 0 {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		symptomIDUint64, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil || symptomIDUint64 == 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		symptomID := uint(symptomIDUint64)

		if delErr := handlers.DeleteSymptom(db.DB, authenticatedPatientID, symptomID); delErr != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(delErr, gorm.ErrRecordNotFound) {
				status = fiber.StatusNotFound
			}
			return ctx.Status(status).JSON(fiber.Map{"error": delErr.Error()})
		}
		return ctx.JSON(fiber.Map{"message": "deleted"})
	})
}
