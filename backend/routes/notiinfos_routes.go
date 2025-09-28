package routes

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
)

func SetupNotiInfosRoutes(api fiber.Router) {

	// CREATE — Fixed Times
	// POST /api/noti/fixed-times
	// Body (ตัวอย่าง):
	// {
	//   "my_medicine_id": 12,           
	//   "group_id": 3,                  
	//   "start_date": "2025-10-01",
	//   "end_date": "2025-10-31",
	//   "times": ["08:00","20:00"]      // >= 1 เวลา
	// }
	api.Post("/noti/fixed-times", func(c *fiber.Ctx) error {
		var body handlers.CreateNotiFixedTimesReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		item, err := handlers.CreateNotiFixedTimes(db.DB, body)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrInvalidData) { status = fiber.StatusBadRequest }
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(item)
	})

	// CREATE — Interval (every N hours)
	// POST /api/noti/interval
	// Body (ตัวอย่าง):
	// {
	//   "my_medicine_id": 12,           // optional
	//   "group_id": 3,                  // optional
	//   "start_date": "2025-10-01",
	//   "end_date": "2025-10-15",
	//   "interval_hours": 6,            // > 0 (บังคับ)
	//   "times_per_day": 3              // optional
	// }
	api.Post("/noti/interval", func(c *fiber.Ctx) error {
		var body handlers.CreateNotiIntervalReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		item, err := handlers.CreateNotiInterval(db.DB, body)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrInvalidData) { status = fiber.StatusBadRequest }
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(item)
	})

	// CREATE — Every N Days
	// POST /api/noti/every-n-days
	// Body (ตัวอย่าง):
	// {
	//   "my_medicine_id": 12,           // optional
	//   "group_id": 3,                  // optional
	//   "start_date": "2025-10-01",
	//   "end_date": "2025-12-31",
	//   "interval_day": 2,              // > 0 (บังคับ)
	//   "times": ["09:00"]              // >= 1 เวลา
	// }
	api.Post("/noti/every-n-days", func(c *fiber.Ctx) error {
		var body handlers.CreateNotiEveryNDaysReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		item, err := handlers.CreateNotiEveryNDays(db.DB, body)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrInvalidData) { status = fiber.StatusBadRequest }
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(item)
	})

	// CREATE — Cycle
	// POST /api/noti/cycle
	// Body (ตัวอย่าง):
	// {
	//   "my_medicine_id": 12,           // optional
	//   "group_id": 3,                  // optional
	//   "start_date": "2025-10-01",
	//   "end_date": "2026-01-31",
	//   "cycle_pattern": [21,7],         // >= 1 ค่า (เช่น [กิน, พัก])
	//   "times": ["08:00","20:00"]      // >= 1 เวลา
	// }
	api.Post("/noti/cycle", func(c *fiber.Ctx) error {
		var body handlers.CreateNotiCycleReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}
		item, err := handlers.CreateNotiCycle(db.DB, body)
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrInvalidData) { status = fiber.StatusBadRequest }
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(item)
	})

	// LIST (optional filters: my_medicine_id, group_id, format_id)
	// GET /api/noti-infos?my_medicine_id=..&group_id=..&format_id=..
	// ตัวอย่าง:
	//   /api/noti-infos?my_medicine_id=12
	//   /api/noti-infos?group_id=3
	//   /api/noti-infos?format_id=4
	//   /api/noti-infos?group_id=3&format_id=2
	api.Get("/noti-infos", func(c *fiber.Ctx) error {
		filter := map[string]any{}
		if v := c.Query("my_medicine_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil { filter["my_medicine_id"] = uint(id) }
		}
		if v := c.Query("group_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil { filter["group_id"] = uint(id) }
		}
		if v := c.Query("format_id"); v != "" {
			if id, err := strconv.ParseUint(v, 10, 64); err == nil { filter["noti_format_id"] = uint(id) }
		}
		items, err := handlers.ListNotiInfos(db.DB, filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": items})
	})

	// GET ONE
	// GET /api/noti-infos/:id
	api.Get("/noti-infos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		item, err := handlers.GetNotiInfo(db.DB, uint(id))
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrRecordNotFound) { status = fiber.StatusNotFound }
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": item})
	})


	// DELETE
	// DELETE /api/noti-infos/:id
	api.Delete("/noti-infos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}
		if err := handlers.DeleteNotiInfo(db.DB, uint(id)); err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrRecordNotFound) { status = fiber.StatusNotFound }
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}
