package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
	"github.com/fouradithep/pillmate/models"
	"strconv"
	"errors"
	"gorm.io/gorm"
	
	"strings"
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
	// ดึง "prescriptions ที่ยังไม่ซิงก์" ของผู้ใช้ (จาก patients.id_card_number)
	// แล้วซิงก์เข้า my_medicines ทีละ "รายการยา (item)" ตามโมเดลใหม่ (หัว + items)
	// เงื่อนไขสำคัญ: ซิงก์เฉพาะที่ sync_until >= TODAY(ตาม time.Local)
	api.Post("/my-medicine/sync-from-prescription", func(c *fiber.Ctx) error {
		// 0) auth
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok || patientID == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		// 1) โหลดเลขบัตรของผู้ใช้ที่ล็อกอิน
		var me models.Patient
		if err := db.DB.Select("id_card_number").Where("id = ?", patientID).First(&me).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		idCard := strings.TrimSpace(me.IDCardNumber)
		if idCard == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id_card_number for this account"})
		}

		// คำนวณ TODAY ตาม timezone ที่เซ็ตไว้ (เช่น Asia/Bangkok)
		today := time.Now().In(time.Local).Truncate(24 * time.Hour)

		// 2) ดึง prescriptions ที่ยังไม่ซิงก์ (พร้อม items) และยังไม่หมดอายุซิงก์
		var prescs []models.Prescription
		q := db.DB.
			Where("id_card_number = ? AND app_sync_status = ?", idCard, false).
			Where("sync_until >= ?", today).
			Preload("Items").
			Order("id ASC")

		if err := q.Find(&prescs).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if len(prescs) == 0 {
			return c.JSON(fiber.Map{
				"message": "no prescriptions to sync (none eligible by sync_until)",
				"synced":  0,
				"today":   today.Format("2006-01-02"),
			})
		}

		type syncedRow struct {
			MyMedicineID       uint   `json:"mymedicine_id"`
			PrescriptionID     uint   `json:"prescription_id"`
			PrescriptionItemID *uint  `json:"prescription_item_id,omitempty"`
			MedName            string `json:"med_name"`
			FormID             uint   `json:"form_id"`
			UnitID             *uint  `json:"unit_id,omitempty"`
			InstructionID      *uint  `json:"instruction_id,omitempty"`
			AmountPerTime      string `json:"amount_per_time"`
			TimesPerDay        string `json:"times_per_day"`
			Source             string `json:"source"`
			MedicineInfoID     *uint  `json:"medicine_info_id,omitempty"`
		}

		created := make([]syncedRow, 0, 16)

		// 3) ทำในทรานแซกชัน: สร้าง MyMedicine ต่อ "รายการยา" แล้วมาร์ก prescription เป็น synced
		if err := db.DB.Transaction(func(tx *gorm.DB) error {
			for _, p := range prescs {
				for _, it := range p.Items {
					// โหลดข้อมูลยาอ้างอิง
					var mi models.MedicineInfo
					if err := tx.First(&mi, it.MedicineInfoID).Error; err != nil {
						return err
					}

					// ป้องกันการสร้างซ้ำ (กรณีเรียกซ้ำ/ fail กลางทาง)
					var exists int64
					check := tx.Model(&models.MyMedicine{}).
						Where("patient_id = ? AND prescription_id = ?", patientID, p.ID)

					// ถ้ามีคอลัมน์ prescription_item_id ใน my_medicines ให้ใช้เป็น key กันซ้ำหลัก
					if tx.Migrator().HasColumn(&models.MyMedicine{}, "prescription_item_id") {
						check = check.Where("prescription_item_id = ?", it.ID)
					}
					if err := check.Count(&exists).Error; err != nil {
						return err
					}
					if exists > 0 {
						continue // ข้ามตัวที่เคยสร้างแล้ว
					}

					// เตรียม pointer ให้ปลอดภัย (รองรับฟิลด์ pointer *uint)
					var (
						unitPtr, instrPtr *uint
						medInfoPtr        *uint
						prescItemIDPtr    *uint
					)

					// medicine_info_id (อ้างอิงต้นทางยา)
					if tx.Migrator().HasColumn(&models.MyMedicine{}, "medicine_info_id") {
						v := mi.ID
						medInfoPtr = &v
					}
					// prescription_item_id (อ้างอิงรายการในใบสั่ง)
					if tx.Migrator().HasColumn(&models.MyMedicine{}, "prescription_item_id") {
						v := it.ID
						prescItemIDPtr = &v
					}
					// unit_id (*uint → เช็ค nil)
					if tx.Migrator().HasColumn(&models.MyMedicine{}, "unit_id") {
						if mi.UnitID != nil {
							unitPtr = mi.UnitID
						}
					}
					// instruction_id (*uint → เช็ค nil)
					if tx.Migrator().HasColumn(&models.MyMedicine{}, "instruction_id") {
						if mi.InstructionID != nil {
							instrPtr = mi.InstructionID
						}
					}

					mm := models.MyMedicine{
						PatientID:     patientID,
						MedName:       mi.MedName,
						Properties:    mi.Properties,
						FormID:        mi.FormID,
						UnitID:        unitPtr,
						InstructionID: instrPtr,
						AmountPerTime: it.AmountPerTime,
						TimesPerDay:   it.TimesPerDay,
						Source:        "hospital",

						PrescriptionID:     &p.ID,
						PrescriptionItemID: prescItemIDPtr,
						MedicineInfoID:     medInfoPtr,
					}

					if err := tx.Create(&mm).Error; err != nil {
						return err
					}

					created = append(created, syncedRow{
						MyMedicineID:       mm.ID,
						PrescriptionID:     p.ID,
						PrescriptionItemID: prescItemIDPtr,
						MedName:            mm.MedName,
						FormID:             mm.FormID,
						UnitID:             unitPtr,
						InstructionID:      instrPtr,
						AmountPerTime:      mm.AmountPerTime,
						TimesPerDay:        mm.TimesPerDay,
						Source:             mm.Source,
						MedicineInfoID:     medInfoPtr,
					})
				}

				// มาร์ก prescription ว่าซิงก์แล้ว
				if err := tx.Model(&models.Prescription{}).
					Where("id = ?", p.ID).
					Update("app_sync_status", true).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message":    "sync from prescriptions successful",
			"patient_id": patientID,
			"synced":     len(created),
			"today":      today.Format("2006-01-02"),
			"data":       created,
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
