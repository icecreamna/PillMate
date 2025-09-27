package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"strconv"
	"errors"
	"gorm.io/gorm"
	"time"
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
	// Body: { "id_card_number": "1101700203451" }
	api.Post("/my-medicine/sync-from-prescription", func(c *fiber.Ctx) error {
		type SyncFromPrescriptionRequest struct {
			IDCardNumber string `json:"id_card_number"`
		}
		var req SyncFromPrescriptionRequest
		if err := c.BodyParser(&req); err != nil || len(req.IDCardNumber) != 13 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id_card_number"})
		}

		// 1) หา patient จากเลขบัตร
		var patient models.Patient
		if err := db.DB.
			Where("id_card_number = ?", req.IDCardNumber).
			First(&patient).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "patient not found"})
		}

		// 2) ดึง prescriptions ที่ยังไม่ซิงก์ และยังไม่หมดเขตซิงก์ (ถ้ามีฟิลด์ SyncUntil)
		now := time.Now()
		var prescriptions []models.Prescription
		query := db.DB.
			Where("id_card_number = ? AND app_sync_status = ?", req.IDCardNumber, false).
			Where("sync_until >= ?", now)

		if err := query.Find(&prescriptions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if len(prescriptions) == 0 {
			return c.JSON(fiber.Map{"message": "no prescriptions to sync", "synced": 0})
		}

		createdMyMedicines := make([]fiber.Map, 0, len(prescriptions))

		// 3) ทำเป็น transaction: สร้าง MyMedicine จากแต่ละใบสั่งยา แล้วมาร์ก sync = true
		if err := db.DB.Transaction(func(tx *gorm.DB) error {
			for _, prescription := range prescriptions {
				// โหลดข้อมูลยาอ้างอิง
				var medicineInfo models.MedicineInfo
				if err := tx.First(&medicineInfo, prescription.MedicineInfoID).Error; err != nil {
					return err
				}

				// เตรียม MyMedicine จาก prescription + medicine info
				myMedicine := models.MyMedicine{
					PatientID:      patient.ID,
					MedName:        medicineInfo.MedName,
					Properties:     medicineInfo.Properties,
					FormID:         medicineInfo.FormID,
					UnitID:         medicineInfo.UnitID,
					InstructionID:  medicineInfo.InstructionID, // หรือ prescription.InstructionID ถ้าจะใช้ตามใบสั่งยา
					AmountPerTime:  prescription.AmountPerTime,
					TimesPerDay:    prescription.TimesPerDay,
					Source:         "hospital",
					PrescriptionID: &prescription.ID, // ผูกใบสั่งยาไว้
				}

				if err := tx.Create(&myMedicine).Error; err != nil {
					return err
				}

				// อัปเดตสถานะ prescription เป็นซิงก์แล้ว
				if err := tx.Model(&models.Prescription{}).
					Where("id = ?", prescription.ID).
					Update("app_sync_status", true).Error; err != nil {
					return err
				}

				createdMyMedicines = append(createdMyMedicines, fiber.Map{
					"mymedicine_id":   myMedicine.ID,
					"prescription_id": prescription.ID,
					"med_name":        myMedicine.MedName,
					"form_id":         myMedicine.FormID,
					"unit_id":         myMedicine.UnitID,
					"instruction_id":  myMedicine.InstructionID,
					"amount_per_time": myMedicine.AmountPerTime,
					"times_per_day":   myMedicine.TimesPerDay,
					"source":          myMedicine.Source,
				})
			}
			return nil
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message":    "sync from prescriptions successful",
			"patient_id": patient.ID,
			"synced":     len(createdMyMedicines),
			"data":       createdMyMedicines,
		})
	})


	// GET /api/forms/:id/units  (ดึงด้วย id)
	api.Get("/forms/:id/units", func(c *fiber.Ctx) error {
	formID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form id"})
	}

	var form models.Form
	if err := db.DB.
		Preload("Units", func(tx *gorm.DB) *gorm.DB {
			return tx.Select("units.id, units.unit_name").Order("units.unit_name ASC")
		}).
		First(&form, formID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// คืนหน่วยที่เชื่อมกับรูปแบบยานั้นๆ
	type UnitLite struct {
		ID       uint   `json:"id"`
		UnitName string `json:"unit_name"`
	}
	units := make([]UnitLite, 0, len(form.Units))
	for _, u := range form.Units {
		units = append(units, UnitLite{ID: u.ID, UnitName: u.UnitName})
	}
	return c.JSON(fiber.Map{"form_id": formID, "units": units})
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
