package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"github.com/fouradithep/pillmate/db"
	"time"
	"strconv"
	"fmt"
	"strings"
)

func SetupPatientRoutes(app *fiber.App) {
	// ผู้ใช้กรอก อีเมล + รหัสผ่าน
	app.Post("/register", func(c *fiber.Ctx) error {
		patient := new(models.Patient)

		// แปลง JSON body -> struct
		if err := c.BodyParser(patient); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		// เรียกใช้ handler สร้าง Patient
		if err := handlers.CreatePatient(db.DB, patient); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Location(fmt.Sprintf("/patient/%d", patient.ID))
		return c.JSON(fiber.Map{
			"message": "Register Successful",
			"patient_id": patient.ID,
		})

	})

	// (body: { "email": "test2@example.com", "password": "1234" })
	// app.Post("/login", func(c *fiber.Ctx) error {
	// 	patient := new(models.Patient)

	// 	if err := c.BodyParser(patient); err != nil {
	// 		return c.SendStatus(fiber.StatusBadRequest)
	// 	}

	// 	token, err := handlers.LoginPatient(db.DB, patient)

	// 	if err != nil {
	// 		return c.SendStatus(fiber.StatusUnauthorized)
	// 	}

	// 	c.Cookie(&fiber.Cookie{
	// 		Name:     "jwt",
	// 		Value:    token,
	// 		Expires:  time.Now().Add(time.Hour * 72),
	// 		HTTPOnly: true,
	// 		Secure:   true,
	// 		SameSite: fiber.CookieSameSiteNoneMode,
	// 	})

	// 	return c.JSON(fiber.Map{
	// 		"message": "Login Successful",
	// 		"token":   token, //ios
	// 	})

	// })

	// app.Post("/logout", func(c *fiber.Ctx) error {
	// 	c.Cookie(&fiber.Cookie{
	// 		Name:     "jwt",  // ชื่อ cookie ที่เคยเก็บ JWT token
	// 		Value:    "",	// กำหนดให้เป็นค่าว่าง = ลบค่าเดิมออก
	// 		Expires:  time.Now().Add(-time.Hour), // ตั้งเวลาให้หมดอายุไปแล้ว = ลบ cookie นี้
	// 		HTTPOnly: true,
	// 		Secure:   true,
	// 		SameSite: fiber.CookieSameSiteNoneMode,
	// 	})
	// 	return c.JSON(fiber.Map{"message": "Logged out"})
	// })
}

func SetupOTPRoutes(app *fiber.App) {
	// POST /patient/:id/otp/request
	app.Post("/patient/:id/otp/request", func(c *fiber.Ctx) error {
		patientID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || patientID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid id"})
		}
		vc, err := handlers.IssueOTP(db.DB, uint(patientID), 3*time.Minute)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error":"failed to issue otp"})
		}
		return c.JSON(fiber.Map{"message":"OTP created and stored", "verification_id": vc.ID, "otp_expires_at": vc.ExpiresAt})
	})

	// POST /patient/:id/otp/resend
	app.Post("/patient/:id/otp/resend", func(c *fiber.Ctx) error {
		patientID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || patientID == 0 {
			return c.Status(400).JSON(fiber.Map{"error":"invalid id"})
		}
		if err := handlers.RevokeActiveOTP(db.DB, uint(patientID)); err != nil {
			return c.Status(500).JSON(fiber.Map{"error":"revoke otp error"})
		}
		vc, err := handlers.IssueOTP(db.DB, uint(patientID), 3*time.Minute)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error":"failed to issue otp"})
		}
		return c.JSON(fiber.Map{"message":"New OTP created and stored", "verification_id": vc.ID, "otp_expires_at": vc.ExpiresAt})
	})

	// POST /patient/:id/otp/verify  (body: { "otp_code": "123456" })
	app.Post("/patient/:id/otp/verify", func(c *fiber.Ctx) error {
		patientID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || patientID == 0 {
			return c.Status(400).JSON(fiber.Map{"error":"invalid id"})
		}
		var req struct{ OTP string `json:"otp_code"` } // ตรงกับคอลัมน์ otp_code
		if err := c.BodyParser(&req); err != nil || len(req.OTP) != 6 {
			return c.Status(400).JSON(fiber.Map{"error":"invalid request"})
		}
		if err := handlers.VerifyOTP(db.DB, uint(patientID), req.OTP); err != nil {
			switch err.Error() {
			case "no otp found, please request a new one":
				return c.Status(404).JSON(fiber.Map{"error": err.Error()})
			case "otp expired", "invalid otp":
				return c.Status(401).JSON(fiber.Map{"error": err.Error()})
			default:
				return c.Status(500).JSON(fiber.Map{"error":"verify error"})
			}
		}
		return c.JSON(fiber.Map{"message":"OTP verified successfully"})
	})

	// กรอกข้อมูลผู้ใช้ที่เหลือ
	// PUT /patient/:id/profile
	app.Put("/patient/:id/profile", func(c *fiber.Ctx) error {
		// ดึง id จาก path
		patientID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || patientID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid patient id"})
		}

		// parse body
		var req models.Patient
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		// update ตรง ๆ
		if err := db.DB.Model(&models.Patient{}).
			Where("id = ?", patientID).
			Updates(req).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "update failed"})
		}

		return c.JSON(fiber.Map{"message": "profile updated"})
	})
}

func SetupForgotPasswordRoutes(app *fiber.App) {
	// POST /patient/password/forgot  { "email": "user@mail.com" }
	app.Post("/patient/password/forgot", func(c *fiber.Ctx) error {
		var req struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.Email) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		// หา patient จากอีเมล  ดึงแค่ id ก็พอ
		var p models.Patient
		if err := db.DB.Select("id").Where("email = ?", req.Email).First(&p).Error; err != nil {
			
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "email not found"})
		}

		// ส่ง id กลับให้ client เอาไปเรียก /patient/:id/otp/request ต่อ
		return c.JSON(fiber.Map{
			"patient_id": p.ID,
		})
	})
}

func SetupPasswordRoutes(app *fiber.App) {
	// POST /patient/:id/reset-password
	// Body: { "new_password": "yourNewSecret" }
	app.Post("/patient/:id/reset-password", func(c *fiber.Ctx) error {
		patientID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || patientID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid patient id"})
		}

		var req struct {
			NewPassword string `json:"new_password"`
		}
		if err := c.BodyParser(&req); err != nil || strings.TrimSpace(req.NewPassword) == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
		}

		if err := handlers.UpdatePatientPassword(db.DB, uint(patientID), req.NewPassword); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to reset password"})
		}

		return c.JSON(fiber.Map{"message": "password reset successful"})
	})
}


