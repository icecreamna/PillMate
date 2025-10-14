package routes

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"gorm.io/gorm"
)

// การใช้งาน:
//
// doctor := app.Group("/doctor/",
//     handlers.AuthAny,
//     handlers.RequireRole("doctor", "admin-app"),
// )
// routes.SetupPrescriptionRoutes(doctor)
//
// เส้นทางที่ได้:
//   POST   /doctor/prescriptions
//   GET    /doctor/prescriptions?q=&doctor_id=
//   GET    /doctor/prescriptions/:id
//   PUT    /doctor/prescriptions/:id
//   DELETE /doctor/prescriptions/:id
//
// หมายเหตุ: อย่าลืมส่ง Header: Content-Type: application/json เมื่อมี Body

func SetupPrescriptionRoutes(api fiber.Router) {
	// =========================
	// CREATE
	// POST /doctor/prescriptions
	//
	// ตัวอย่าง Body (CreatePrescriptionDTO):
	// {
	//   "id_card_number": "1101700234567",
	//   "doctor_id": 7, // (optional) ถ้าไม่ส่งจะดึงจาก token ฝั่งเซิร์ฟเวอร์
	//   "items": [
	//     { "medicine_info_id": 12, "amount_per_time": "1 เม็ด", "times_per_day": "2 ครั้ง" },
	//     { "medicine_info_id": 45, "amount_per_time": "5 ml",  "times_per_day": "3 ครั้ง" }
	//   ],
	//   "sync_until": "2025-12-31T00:00:00+07:00", // (optional) ไม่ส่ง = +60 วัน (กำหนดใน hook)
	//   "app_sync_status": false                    // (optional) default=false
	// }
	// =========================
	api.Post("/prescriptions", func(c *fiber.Ctx) error {
		// บังคับ JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		var in dto.CreatePrescriptionDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		// ถ้าไม่ได้ส่ง doctor_id มา และ token มี admin_id อยู่ ให้ใส่แทน
		if in.DoctorID == 0 {
			if v, ok := c.Locals("admin_id").(uint); ok && v > 0 {
				in.DoctorID = v
			}
		}

		created, err := handlers.CreatePrescription(db.DB, &in)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "created",
			"data":    created, // จะมี items กลับมาด้วย
		})
	})

	// =========================
	// LIST (no pagination) + ค้นหา q + ตัวเลือกกรอง doctor_id
	// GET /doctor/prescriptions?q=&doctor_id=
	//
	// ตัวอย่างการเรียก:
	//   GET /doctor/prescriptions
	//   GET /doctor/prescriptions?q=1101700
	//   GET /doctor/prescriptions?doctor_id=7
	//   GET /doctor/prescriptions?q=4567&doctor_id=7
	// =========================
	api.Get("/prescriptions", func(c *fiber.Ctx) error {
		q := c.Query("q", "")

		var docIDPtr *uint
		// อ่านจาก query ก่อน (ถ้าส่งมา)
		if s := strings.TrimSpace(c.Query("doctor_id", "")); s != "" {
			if v, err := strconv.ParseUint(s, 10, 64); err == nil && v > 0 {
				vv := uint(v)
				docIDPtr = &vv
			}
		}
		// ถ้าไม่ส่งมา แต่ token เป็น doctor ให้กรองด้วยตัวเอง (optional – ปรับตามนโยบาย)
		if docIDPtr == nil {
			if v, ok := c.Locals("admin_id").(uint); ok && v > 0 {
				vv := v
				docIDPtr = &vv
			}
		}

		items, err := handlers.ListPrescriptions(db.DB, q, docIDPtr)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// จะได้แต่ละใบพร้อม items
		return c.JSON(fiber.Map{"data": items})
	})

	// =========================
	// GET ONE
	// GET /doctor/prescriptions/:id
	// =========================
	api.Get("/prescriptions/:id", func(c *fiber.Ctx) error {
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		rec, err := handlers.GetPrescriptionByID(db.DB, uint(idU64))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": rec}) // preload items แล้ว
	})

	// =========================
	// UPDATE (partial เฉพาะหัวเอกสาร)
	// PUT /doctor/prescriptions/:id
	//
	// ตัวอย่าง Body (UpdatePrescriptionDTO):
	// {
	//   "id_card_number": "1101700234567",        // optional
	//   "doctor_id": 7,                            // optional
	//   "sync_until": "2026-01-31T00:00:00+07:00", // optional
	//   "app_sync_status": true                    // optional
	// }
	//
	// หมายเหตุ: การแก้ไข items ให้ทำเป็น endpoint แยก (replace/add/update/delete) ภายหลัง
	// =========================
	api.Put("/prescriptions/:id", func(c *fiber.Ctx) error {
		// บังคับ JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var in dto.UpdatePrescriptionDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		updated, err := handlers.UpdatePrescription(db.DB, uint(idU64), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{
			"message": "updated",
			"data":    updated,
		})
	})

	// =========================
	// DELETE (soft delete ทั้งหัวและ items)
	// DELETE /doctor/prescriptions/:id
	// =========================
	api.Delete("/prescriptions/:id", func(c *fiber.Ctx) error {
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if err := handlers.DeletePrescription(db.DB, uint(idU64)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}
