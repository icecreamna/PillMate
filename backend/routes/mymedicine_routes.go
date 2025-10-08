package routes

import (
	"errors"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"

	"strings"
)

func SetupMyMedicineRoutes(api fiber.Router) {

	// add MyMedicine
	// POST /api/my-medicine
	// 	 Body
	// {
	//   "med_name": "Paracetamol",
	//   "properties": "แก้ปวด ลดไข้",
	//   "form_id": 1,
	//   "unit_id": 1,
	//   "instruction_id": 2,
	//   "amount_per_time": "1",
	//   "times_per_day": "3",
	//   "source": "manual"
	// }
	api.Post("/my-medicine", func(c *fiber.Ctx) error {
		// ดึง patient_id ที่ middleware ใส่ไว้
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var mymedicine models.MyMedicine
		if err := c.BodyParser(&mymedicine); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		// กัน client ส่ง id/patient_id มาเอง
		mymedicine.ID = 0
		mymedicine.PatientID = patientID
		if mymedicine.Source == "" {
			mymedicine.Source = "manual"
		}

		created, err := handlers.AddMyMedicine(db.DB, &mymedicine)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Add MyMedicine Successful",
			"data": fiber.Map{
				"id":              created.ID,
				"patient_id":      created.PatientID,
				"med_name":        created.MedName,
				"properties":      created.Properties,
				"form_id":         created.FormID,
				"unit_id":         created.UnitID,
				"instruction_id":  created.InstructionID,
				"amount_per_time": created.AmountPerTime,
				"times_per_day":   created.TimesPerDay,
				"source":          created.Source,
			},
		})
	})

	// POST /api/my-medicine/sync-from-prescription
	// ดึง prescriptions ของผู้ใช้ที่ล็อกอิน (อิงจาก id_card_number ในตาราง patients) แล้วซิงก์เข้า my_medicines
	api.Post("/my-medicine/sync-from-prescription", func(c *fiber.Ctx) error {
		// auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// 1) โหลดเลขบัตรของผู้ใช้ที่ล็อกอิน
		var me models.Patient
		if err := db.DB.Select("id_card_number").Where("id = ?", patientID).First(&me).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		var idCard string
		if me.IDCardNumber != nil {
			idCard = strings.TrimSpace(*me.IDCardNumber)
		} else {
			idCard = ""
		}

		if idCard == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "missing id_card_number for this account",
			})
		}

		// 2) ดึง prescriptions ที่ยังไม่ซิงก์
		var prescriptions []models.Prescription
		if err := db.DB.
			Where("id_card_number = ? AND app_sync_status = ?", idCard, false).
			// ถ้ามีฟิลด์ SyncUntil และอยากเช็คหมดอายุ ให้ปลดคอมเมนต์และ import time
			// Where("sync_until >= ?", time.Now()).
			Find(&prescriptions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if len(prescriptions) == 0 {
			return c.JSON(fiber.Map{"message": "no prescriptions to sync", "synced": 0})
		}

		createdMyMedicines := make([]fiber.Map, 0, len(prescriptions))

		// 3) ทำเป็น transaction: สร้าง MyMedicine ทีละใบสั่งยา แล้วมาร์ก prescription.app_sync_status = true
		if err := db.DB.Transaction(func(tx *gorm.DB) error {
			for _, p := range prescriptions {
				// โหลดข้อมูลยาอ้างอิง
				var mi models.MedicineInfo
				if err := tx.First(&mi, p.MedicineInfoID).Error; err != nil {
					return err
				}

				// เตรียม MyMedicine จาก prescription + medicine info
				mm := models.MyMedicine{
					PatientID:      patientID,
					MedName:        mi.MedName,
					Properties:     mi.Properties,
					FormID:         mi.FormID,
					UnitID:         mi.UnitID,
					InstructionID:  mi.InstructionID, // หรือใช้จาก p ถ้ามีใน prescription
					AmountPerTime:  p.AmountPerTime,
					TimesPerDay:    p.TimesPerDay,
					Source:         "hospital",
					PrescriptionID: &p.ID,
				}

				if err := tx.Create(&mm).Error; err != nil {
					return err
				}

				// มาร์ก prescription ว่าถูกซิงก์แล้ว
				if err := tx.Model(&models.Prescription{}).
					Where("id = ?", p.ID).
					Update("app_sync_status", true).Error; err != nil {
					return err
				}

				createdMyMedicines = append(createdMyMedicines, fiber.Map{
					"mymedicine_id":   mm.ID,
					"prescription_id": p.ID,
					"med_name":        mm.MedName,
					"form_id":         mm.FormID,
					"unit_id":         mm.UnitID,
					"instruction_id":  mm.InstructionID,
					"amount_per_time": mm.AmountPerTime,
					"times_per_day":   mm.TimesPerDay,
					"source":          mm.Source,
				})
			}
			return nil
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message":    "sync from prescriptions successful",
			"patient_id": patientID,
			"synced":     len(createdMyMedicines),
			"data":       createdMyMedicines,
		})
	})

	// Read All
	// GET /api/my-medicines
	api.Get("/my-medicines", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		mymedicines, err := handlers.GetMyMedicines(db.DB, patientID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Get MyMedicines Successful",
			"data":    mymedicines, // << ส่งทั้ง slice กลับไปเลย
		})
	})

	// Read One
	// GET /api/my-medicine/:id
	api.Get("/my-medicine/:id", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// อ่าน mymedicine id จาก path
		mymedicineID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || mymedicineID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid medicine id"})
		}

		// เรียก handler
		mymedicine, err := handlers.GetMyMedicine(db.DB, patientID, uint(mymedicineID))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "medicine not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// สำเร็จ
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Get MyMedicine Successful",
			"data":    mymedicine,
		})
	})

	// PUT /api/my-medicine/:id
	api.Put("/my-medicine/:id", func(c *fiber.Ctx) error {
		// เอา patient_id จาก middleware
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// อ่าน mymedicine id จาก path
		mymedicineID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || mymedicineID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid medicine id"})
		}

		// รับ body ที่จะอัปเดต
		var in models.MyMedicine
		if err := c.BodyParser(&in); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		// เรียก handler
		updated, err := handlers.UpdateMyMedicine(db.DB, patientID, uint(mymedicineID), &in)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "medicine not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "Updated mymedicine successful",
			"data":    updated,
		})
	})

	// DELETE /api/my-medicine/:id
	api.Delete("/my-medicine/:id", func(c *fiber.Ctx) error {
		// ดึง patient_id จาก middleware
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// อ่าน mymedicine id จาก path
		mymedicineID, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || mymedicineID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid medicine id"})
		}

		// ลบ
		if err := handlers.DeleteMyMedicine(db.DB, patientID, uint(mymedicineID)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "medicine not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "Mymedicine deleted"})
	})
}
