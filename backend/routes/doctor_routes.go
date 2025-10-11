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

// ต้องเรียกจากกลุ่มที่มี middleware แล้ว เช่น
// admin := app.Group("/admin/",
//     handlers.AuthAny,
//     handlers.RequireRole("superadmin", "admin-app"),
// )
// routes.SetupDoctorRoutes(admin)

func SetupDoctorRoutes(api fiber.Router) {
	// CREATE
	// POST /admin/doctors
	// Body (CreateDoctorDTO):
	// {
	//   "username": "doc1",
	//   "password": "secret",
	//   "first_name": "Somchai",
	//   "last_name": "Jaidee"
	// }
	api.Post("/doctors", func(c *fiber.Ctx) error {
		// บังคับให้เป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		var in dto.CreateDoctorDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		// ดึง actorID จาก token ถ้ามี (ไม่บังคับ)
		var actorID uint
		if v, ok := c.Locals("admin_id").(uint); ok {
			actorID = v
		}

		created, err := handlers.CreateDoctor(db.DB, &in, actorID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// ✅ ตอบเป็น DTO (ไม่คืน password)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "created",
			"data":    dto.ToWebAdminDTO(created),
		})
	})

	// LIST (ค้นหา + แบ่งหน้า)
	// GET /admin/doctors?q=&page=&page_size=
	api.Get("/doctors", func(c *fiber.Ctx) error {
    q := c.Query("q", "")
    page, _ := strconv.Atoi(c.Query("page", "1"))
    if page <= 0 { page = 1 }
    pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
    if pageSize <= 0 { pageSize = 20 }

    list, _, err := handlers.ListDoctors(db.DB, q, page, pageSize)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // คืนเฉพาะ data (แปลงเป็น DTO ด้วย)
    return c.JSON(fiber.Map{
        "data": dto.ToWebAdminDTOs(list),
    })
})

	// GET ONE
	// GET /admin/doctors/:id
	api.Get("/doctors/:id", func(c *fiber.Ctx) error {
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		doc, err := handlers.GetDoctorByID(db.DB, uint(idU64))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		// ✅ ตอบเป็น DTO
		return c.JSON(fiber.Map{"data": dto.ToWebAdminDTO(doc)})
	})

	// UPDATE (partial via DTO ที่เป็น pointer fields)
	// PUT /admin/doctors/:id
	// Body (UpdateDoctorDTO):
	// {
	//   "username": "newuser",       // optional
	//   "first_name": "NewName",     // optional
	//   "last_name": "Last",         // optional
	//   "password": "newpass"        // optional
	// }
	api.Put("/doctors/:id", func(c *fiber.Ctx) error {
		// บังคับให้เป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var in dto.UpdateDoctorDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		updated, err := handlers.UpdateDoctor(db.DB, uint(idU64), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		// ✅ ตอบเป็น DTO
		return c.JSON(fiber.Map{"message": "updated", "data": dto.ToWebAdminDTO(updated)})
	})

	// DELETE (soft/hard ตามโมเดล)
	// DELETE /admin/doctors/:id
	api.Delete("/doctors/:id", func(c *fiber.Ctx) error {
		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if err := handlers.DeleteDoctor(db.DB, uint(idU64)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})

	// PATCH /admin/doctors/:id/password  (admin reset password)
	// {
	//   "new_password": "yourNewSecret123"
    // }
	api.Patch("/doctors/:id/password", func(c *fiber.Ctx) error {
		// บังคับเป็น JSON
		if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
			return c.Status(fiber.StatusUnsupportedMediaType).
				JSON(fiber.Map{"error": "Content-Type must be application/json"})
		}

		idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || idU64 == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var in dto.AdminResetPasswordDTO
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		updated, err := handlers.ResetDoctorPassword(db.DB, uint(idU64), in.NewPassword)
		if err != nil {
			// เช่น รหัสสั้นเกิน / ซ้ำรหัสเดิม / ไม่พบผู้ใช้
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "password reset success",
			"data":    dto.ToWebAdminDTO(updated),
		})
	})
}
