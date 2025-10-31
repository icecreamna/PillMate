package routes

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/dto"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
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
	// api.Post("/doctors", func(c *fiber.Ctx) error {
	// 	// บังคับให้เป็น JSON
	// 	if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
	// 		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": "Content-Type must be application/json"})
	// 	}

	// 	var in dto.CreateDoctorDTO
	// 	if err := c.BodyParser(&in); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	// 	}

	// 	// ดึง actorID จาก token ถ้ามี (ไม่บังคับ)
	// 	var actorID uint
	// 	if v, ok := c.Locals("admin_id").(uint); ok {
	// 		actorID = v
	// 	}

	// 	created, err := handlers.CreateDoctor(db.DB, &in, actorID)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// 	}

	// 	// ✅ ตอบเป็น DTO (ไม่คืน password)
	// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
	// 		"message": "created",
	// 		"data":    dto.ToWebAdminDTO(created),
	// 	})
	// })

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
	// api.Put("/doctors/:id", func(c *fiber.Ctx) error {
	// 	// บังคับให้เป็น JSON
	// 	if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
	// 		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{"error": "Content-Type must be application/json"})
	// 	}

	// 	idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	// 	if err != nil || idU64 == 0 {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	// 	}

	// 	var in dto.UpdateDoctorDTO
	// 	if err := c.BodyParser(&in); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	// 	}

	// 	updated, err := handlers.UpdateDoctor(db.DB, uint(idU64), &in)
	// 	if err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	// 		}
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	// ✅ ตอบเป็น DTO
	// 	return c.JSON(fiber.Map{"message": "updated", "data": dto.ToWebAdminDTO(updated)})
	// })

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
	// api.Patch("/doctors/:id/password", func(c *fiber.Ctx) error {
	// 	// บังคับเป็น JSON
	// 	if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
	// 		return c.Status(fiber.StatusUnsupportedMediaType).
	// 			JSON(fiber.Map{"error": "Content-Type must be application/json"})
	// 	}

	// 	idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	// 	if err != nil || idU64 == 0 {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	// 	}

	// 	var in dto.AdminResetPasswordDTO
	// 	if err := c.BodyParser(&in); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	// 	}

	// 	updated, err := handlers.ResetDoctorPassword(db.DB, uint(idU64), in.NewPassword)
	// 	if err != nil {
	// 		// เช่น รหัสสั้นเกิน / ซ้ำรหัสเดิม / ไม่พบผู้ใช้
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"message": "password reset success",
	// 		"data":    dto.ToWebAdminDTO(updated),
	// 	})
	// })
}

// ต้องเรียกจากกลุ่มที่มี middleware แล้ว เช่น
// doctor := app.Group("/doctor/",
//     handlers.AuthAny,
//     handlers.RequireRole("doctor", "admin-app"),
// )
// routes.SetupDoctorPublicRoutes(doctor)
//
// เส้นทางที่ได้:
//   GET /doctor/me
//   GET /doctor/doctors/:id

func SetupDoctorPublicRoutes(api fiber.Router) {

    // =========================================================
    // GET /doctor/me — ดึงโปรไฟล์ตัวเองจาก token (admin_id)
    //
    // ✅ ตัวอย่าง Response 200:
    // {
    //   "data": {
    //     "id": 12,
    //     "username": "doc.somchai",
    //     "first_name": "สมชาย",
    //     "last_name": "ใจดี",
    //     "role": "doctor",
    //     "created_at": "2025-10-30T12:34:56+07:00",
    //     "updated_at": "2025-10-30T12:34:56+07:00"
    //   }
    // }
    //
    // ❌ ตัวอย่าง Response 401 (ไม่มี token / token ไม่ถูกต้อง):
    // { "error": "unauthorized" }
    //
    // ❌ ตัวอย่าง Response 404 (ไม่พบผู้ใช้):
    // { "error": "not found" }
    // =========================================================
    api.Get("/me", func(c *fiber.Ctx) error {
        v, ok := c.Locals("admin_id").(uint)
        if !ok || v == 0 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
        }

        doc, err := handlers.GetDoctorByID(db.DB, v)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
            }
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"data": dto.ToWebAdminDTO(doc)})
    })

    // =========================================================
    // PUT /doctor/me — หมอแก้โปรไฟล์ตัวเอง (re-use UpdateDoctor)
    //
    // ⚠️ Header: Content-Type: application/json
    //
    // ✅ ตัวอย่าง Request:
    // {
    //   "username": "doc.somchai",      // optional
    //   "first_name": "สมชาย",          // optional
    //   "last_name": "อัปเดต"           // optional
    // }
    //
    // ✅ ตัวอย่าง Response 200:
    // {
    //   "message": "updated",
    //   "data": {
    //     "id": 12,
    //     "username": "doc.somchai",
    //     "first_name": "สมชาย",
    //     "last_name": "อัปเดต",
    //     "role": "doctor",
    //     "created_at": "2025-10-30T12:34:56+07:00",
    //     "updated_at": "2025-10-31T09:00:00+07:00"
    //   }
    // }
    //
    // ❌ ตัวอย่าง Response 415 (Content-Type ไม่ใช่ JSON):
    // { "error": "Content-Type must be application/json" }
    //
    // ❌ ตัวอย่าง Response 400 (body ไม่ถูกต้อง / username ซ้ำ):
    // { "error": "invalid request body" }
    // { "error": "username already exists" }
    // =========================================================
    api.Put("/me", func(c *fiber.Ctx) error {
        v, ok := c.Locals("admin_id").(uint)
        if !ok || v == 0 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
        }
        if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
            return c.Status(fiber.StatusUnsupportedMediaType).
                JSON(fiber.Map{"error": "Content-Type must be application/json"})
        }

        var in dto.UpdateDoctorDTO
        if err := c.BodyParser(&in); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
        }

        // ไม่ให้เปลี่ยนรหัสผ่านผ่าน endpoint นี้ (ใช้ /doctor/me/password แทน)
        in.Password = nil

        updated, err := handlers.UpdateDoctor(db.DB, v, &in)
        if err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
            }
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"message": "updated", "data": dto.ToWebAdminDTO(updated)})
    })

    // =========================================================
    // PATCH /doctor/me/password — หมอเปลี่ยนรหัสผ่านตัวเอง
    //
    // ⚠️ Header: Content-Type: application/json
    //
    // ✅ ตัวอย่าง Request:
    // {
    //   "old_password": "OldPass@1234",
    //   "new_password": "NewPass@1234"
    // }
    //
    // ✅ ตัวอย่าง Response 200:
    // {
    //   "message": "password changed",
    //   "data": {
    //     "id": 12,
    //     "username": "doc.somchai",
    //     "first_name": "สมชาย",
    //     "last_name": "อัปเดต",
    //     "role": "doctor",
    //     "created_at": "2025-10-30T12:34:56+07:00",
    //     "updated_at": "2025-10-31T09:05:00+07:00"
    //   }
    // }
    //
    // ❌ ตัวอย่าง Response 400:
    // { "error": "old password is incorrect" }
    // { "error": "new password must be different from current password" }
    // { "error": "invalid request body" }
    //
    // ❌ ตัวอย่าง Response 415:
    // { "error": "Content-Type must be application/json" }
    // =========================================================
    type doctorChangePasswordDTO struct {
        OldPassword string `json:"old_password"`
        NewPassword string `json:"new_password"`
    }
    api.Patch("/me/password", func(c *fiber.Ctx) error {
        v, ok := c.Locals("admin_id").(uint)
        if !ok || v == 0 {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
        }
        if ct := c.Get("Content-Type"); !strings.Contains(ct, "application/json") {
            return c.Status(fiber.StatusUnsupportedMediaType).
                JSON(fiber.Map{"error": "Content-Type must be application/json"})
        }

        var in doctorChangePasswordDTO
        if err := c.BodyParser(&in); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
        }

        // อ่าน hash เดิมเพื่อตรวจ old_password
        var me models.WebAdmin
        if err := db.DB.Where("role = ?", "doctor").First(&me, v).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
            }
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }
        if err := bcrypt.CompareHashAndPassword([]byte(me.Password), []byte(strings.TrimSpace(in.OldPassword))); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "old password is incorrect"})
        }

        // ใช้ handler เดิมรีเซ็ต (ภายในกัน new==old อยู่แล้ว)
        updated, err := handlers.ResetDoctorPassword(db.DB, v, in.NewPassword)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(fiber.Map{"message": "password changed", "data": dto.ToWebAdminDTO(updated)})
    })

    // (ออปชันนัล) GET /doctor/doctors/:id — ถ้าต้องการให้หมอดูหมอคนอื่นได้
    //
    // ✅ ตัวอย่าง Response 200:
    // {
    //   "data": {
    //     "id": 99,
    //     "username": "doc.another",
    //     "first_name": "หมอ",
    //     "last_name": "อีกคน",
    //     "role": "doctor",
    //     "created_at": "...",
    //     "updated_at": "..."
    //   }
    // }
    //
    // ❌ ตัวอย่าง Response 400/404:
    // { "error": "invalid id" }
    // { "error": "not found" }
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
        return c.JSON(fiber.Map{"data": dto.ToWebAdminDTO(doc)})
    })
}



