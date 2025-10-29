package routes

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"gorm.io/gorm"
)

// จะได้เส้นทาง:
//   GET /api/appointments/latest   -> ใบนัดล่าสุดของผู้ใช้ที่ล็อกอิน
//   GET /api/appointments/:id      -> ใบนัดตาม id (ต้องเป็นของผู้ใช้เอง)
func SetupMobileAppointmentRoutes(api fiber.Router) {
	// GET /api/appointments/latest
	// api.Get("/appointments/latest", func(c *fiber.Ctx) error {
	// 	// ปิด cache responses (ข้อมูลส่วนบุคคล)
	// 	c.Set("Cache-Control", "no-store")

	// 	patientID, _ := c.Locals("patient_id").(uint) // ต้องมีจาก middleware auth ฝั่ง mobile
	// 	if patientID == 0 {
	// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	// 	}

	// 	latest, err := handlers.MobileGetLatestAppointment(db.DB, patientID)
	// 	if err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no appointment"})
	// 		}
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	return c.JSON(fiber.Map{"data": latest})
	// })

	// // GET /api/appointments/:id
	// api.Get("/appointments/:id", func(c *fiber.Ctx) error {
	// 	// ปิด cache responses (ข้อมูลส่วนบุคคล)
	// 	c.Set("Cache-Control", "no-store")

	// 	patientID, _ := c.Locals("patient_id").(uint)
	// 	if patientID == 0 {
	// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	// 	}

	// 	idU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	// 	if err != nil || idU64 == 0 {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	// 	}

	// 	rec, err := handlers.MobileGetAppointmentByID(db.DB, patientID, uint(idU64))
	// 	if err != nil {
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	// 		}
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	return c.JSON(fiber.Map{"data": rec})
	// })
	// =========================================================
	// GET ONE
	// GET /api/appointment/:id
	// อิงสิทธิ์จาก patient_id -> โหลด id_card_number แล้วคิวรีด้วย (id, id_card_number)
	// =========================================================
	// api.Get("/appointment/", func(c *fiber.Ctx) error {
	// 	// ⛑️ auth
	// 	patientID, ok := c.Locals("patient_id").(uint)
	// 	if !ok || patientID == 0 {
	// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	// 	}

	// 	// แปลงพารามิเตอร์ id
	// 	appointmentIDU64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	// 	if err != nil || appointmentIDU64 == 0 {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	// 	}
	// 	appointmentID := uint(appointmentIDU64)

	// 	// โหลดเลขบัตรประชาชนจากผู้ใช้ที่ล็อกอิน
	// 	var me models.Patient
	// 	if err := db.DB.Select("id_card_number").Where("id = ?", patientID).First(&me).Error; err != nil {
	// 		status := fiber.StatusInternalServerError
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			status = fiber.StatusUnauthorized
	// 		}
	// 		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	if me.IDCardNumber == nil || strings.TrimSpace(*me.IDCardNumber) == "" {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "missing id_card_number for this account",
	// 		})
	// 	}
	// 	idCard := strings.TrimSpace(*me.IDCardNumber)

	// 	// ดึงใบนัดหนึ่งรายการ (ยืนยันด้วย id + id_card_number)
	// 	appointment, getErr := handlers.GetAppointment(db.DB, appointmentID, idCard)
	// 	if getErr != nil {
	// 		status := fiber.StatusInternalServerError
	// 		if errors.Is(getErr, gorm.ErrRecordNotFound) {
	// 			status = fiber.StatusNotFound
	// 		}
	// 		return c.Status(status).JSON(fiber.Map{"error": getErr.Error()})
	// 	}
	// 	return c.JSON(fiber.Map{"data": dto.AppointmentToDTO(*appointment)})
	// })

	api.Get("/appointments/next", func(c *fiber.Ctx) error {
		// ปิด cache responses (ข้อมูลส่วนบุคคล)
		c.Set("Cache-Control", "no-store")

		// ⛑️ ตรวจสอบ auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// ✅ เรียก handler
		nextAppointment, err := handlers.GetNextAppointment(db.DB, patientID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "no upcoming appointment found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// ✅ แปลงเป็นเวลาไทยสำหรับการแสดงผล (ให้สอดคล้องกับ DTO เดิม)
		bkk, _ := time.LoadLocation("Asia/Bangkok")
		dateStr := nextAppointment.AppointmentDate.In(bkk).Format("2006-01-02")
		timeStr := nextAppointment.AppointmentTime.In(bkk).Format("15:04")

		// ✅ ส่งออกเฉพาะข้อมูลที่ต้องการ
		return c.JSON(fiber.Map{
			"appointment_date": dateStr,
			"appointment_time": timeStr,
			"note":             nextAppointment.Note, // *string -> JSON จะเป็นค่า string หรือ null
		})
	})

	// api.Post("/appointment", func(c *fiber.Ctx) error {
	// 	var req struct {
	// 		IDCardNumber    string `json:"id_card_number"`
	// 		AppointmentDate string `json:"appointment_date"`
	// 		AppointmentTime string `json:"appointment_time"`
	// 		Note            string `json:"note"`
	// 	}

	// 	// ✅ parse body
	// 	if err := c.BodyParser(&req); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "invalid request body",
	// 		})
	// 	}

	// 	// ✅ validate ข้อมูลเบื้องต้น
	// 	if strings.TrimSpace(req.IDCardNumber) == "" ||
	// 		strings.TrimSpace(req.AppointmentDate) == "" ||
	// 		strings.TrimSpace(req.AppointmentTime) == "" {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "id_card_number, appointment_date, and appointment_time are required",
	// 		})
	// 	}

	// 	// ✅ validate format วันที่ / เวลา
	// 	parsedDate, err := time.Parse("2006-01-02", req.AppointmentDate)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "invalid date format (use YYYY-MM-DD)",
	// 		})
	// 	}

	// 	parsedTime, err := time.Parse("15:04", req.AppointmentTime)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "invalid time format (use HH:mm)",
	// 		})
	// 	}

	// 	// ✅ สร้าง struct สำหรับ handler
	// 	appointment := models.Appointment{
	// 		IDCardNumber:    req.IDCardNumber,
	// 		AppointmentDate: parsedDate,
	// 		AppointmentTime: parsedTime,
	// 		DoctorID:        1, // mock doctor id
	// 		Note:            req.Note,
	// 	}

	// 	// ✅ ให้ handler ไปจัดการ DB
	// 	result, err := handlers.CreateAppointment(appointment)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": err.Error(),
	// 		})
	// 	}

	// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
	// 		"message": "appointment created successfully",
	// 		"data":    result,
	// 	})
	// })

	// =========================================================
	// LIST
	// GET /api/appointments?date_from=YYYY-MM-DD&date_to=YYYY-MM-DD&hospital_id=&doctor_id=
	// หมายเหตุ: ใช้ id_card_number จากผู้ใช้ที่ล็อกอินเพื่อจำกัดเรคคอร์ด
	// =========================================================
	// api.Get("/appointments", func(c *fiber.Ctx) error {
	// 	// ⛑️ auth
	// 	patientID, ok := c.Locals("patient_id").(uint)
	// 	if !ok || patientID == 0 {
	// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	// 	}

	// 	// โหลดเลขบัตรประชาชนจากผู้ใช้ที่ล็อกอิน
	// 	var me models.Patient
	// 	if err := db.DB.Select("id_card_number").Where("id = ?", patientID).First(&me).Error; err != nil {
	// 		status := fiber.StatusInternalServerError
	// 		if errors.Is(err, gorm.ErrRecordNotFound) {
	// 			status = fiber.StatusUnauthorized
	// 		}
	// 		return c.Status(status).JSON(fiber.Map{"error": err.Error()})
	// 	}
	// 	if me.IDCardNumber == nil || strings.TrimSpace(*me.IDCardNumber) == "" {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"error": "missing id_card_number for this account",
	// 		})
	// 	}

	// 	idCard := strings.TrimSpace(*me.IDCardNumber)

	// 	// อ่าน query params
	// 	dateFromStr := strings.TrimSpace(c.Query("date_from"))
	// 	dateToStr := strings.TrimSpace(c.Query("date_to"))
	// 	hospitalIDStr := strings.TrimSpace(c.Query("hospital_id"))
	// 	doctorIDStr := strings.TrimSpace(c.Query("doctor_id"))

	// 	const ymd = "2006-01-02"

	// 	q := db.DB.Model(&models.Appointment{}).
	// 		Preload("Hospital").
	// 		Preload("WebAdmin").
	// 		Where("id_card_number = ?", idCard)

	// 	// กรองช่วงวันที่ (ถ้าส่งมา)
	// 	if dateFromStr != "" {
	// 		if t, err := time.ParseInLocation(ymd, dateFromStr, time.Local); err == nil {
	// 			q = q.Where("appointment_date >= ?", t)
	// 		}
	// 	}
	// 	if dateToStr != "" {
	// 		if t, err := time.ParseInLocation(ymd, dateToStr, time.Local); err == nil {
	// 			q = q.Where("appointment_date <= ?", t)
	// 		}
	// 	}

	// 	// กรองด้วย hospital_id / doctor_id (ถ้าส่งมา)
	// 	if hospitalIDStr != "" {
	// 		if hid, err := strconv.ParseUint(hospitalIDStr, 10, 64); err == nil {
	// 			q = q.Where("hospital_id = ?", uint(hid))
	// 		}
	// 	}
	// 	if doctorIDStr != "" {
	// 		if did, err := strconv.ParseUint(doctorIDStr, 10, 64); err == nil {
	// 			q = q.Where("doctor_id = ?", uint(did))
	// 		}
	// 	}

	// 	var appointments []models.Appointment
	// 	if err := q.Order("appointment_date, appointment_time").Find(&appointments).Error; err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	// 	}

	// 	return c.JSON(fiber.Map{"data": dto.AppointmentsToDTO(appointments)})
	// })
}
