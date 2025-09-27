package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"strconv"
	"errors"
	"gorm.io/gorm"
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
