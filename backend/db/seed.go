package db

import (
	"log"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

// uint → *uint (ช่วยตอน seed/assign ค่าให้ฟิลด์ pointer)
func UintPtr(v uint) *uint { return &v }

func SeedInitialData(db *gorm.DB) {
	// seed ข้อมูล Form
	forms := []models.Form{
		{FormName: "ยาเม็ด"},
		{FormName: "แคปซูล"},
		{FormName: "ยาน้ำ"},
		{FormName: "ยาฉีด"},
		{FormName: "ยาใช้ทา"},
		{FormName: "ยาใช้หยด"},
	}
	for _, form := range forms {
		if err := db.FirstOrCreate(&form, models.Form{FormName: form.FormName}).Error; err != nil {
			log.Println("Seed form error:", err)
		}
	}

	// Seed ข้อมูล Unit
	units := []models.Unit{
		{UnitName: "เม็ด"},
		{UnitName: "แคปซูล"},
		{UnitName: "ช้อนชา"},
		{UnitName: "ช้อนโต๊ะ"},
		{UnitName: "มิลลิลิตร"},
		{UnitName: "cc"},
		{UnitName: "ยูนิต"},
		{UnitName: "มิลลิกรัม"},
		{UnitName: "ไมโครกรัม"},
		{UnitName: "กรัม"},
		{UnitName: "หลอด"},
		{UnitName: "หยด"},
	}

	for i := range units {
		if err := db.FirstOrCreate(&units[i], models.Unit{UnitName: units[i].UnitName}).Error; err != nil {
			log.Println("Seed unit error:", err)
		}
	}

	// ดึง Form "ยาเม็ด" เพื่อเอา ID
	var formTablet models.Form
	if err := db.Where("form_name = ?", "ยาเม็ด").First(&formTablet).Error; err != nil {
		log.Println("ไม่พบ Form ยาเม็ด:", err)
		return
	}

	// ดึง Form "แคปซูล"
	var formCapsule models.Form
	if err := db.Where("form_name = ?", "แคปซูล").First(&formCapsule).Error; err != nil {
		log.Println("ไม่พบ Form แคปซูล:", err)
		return
	}

	// ดึง Form "ยาน้ำ"
	var formLiquid models.Form
	if err := db.Where("form_name = ?", "ยาน้ำ").First(&formLiquid).Error; err != nil {
		log.Println("ไม่พบ Form ยาน้ำ:", err)
		return
	}

	// ดึง Form "ยาฉีด"
	var formInjection models.Form
	if err := db.Where("form_name = ?", "ยาฉีด").First(&formInjection).Error; err != nil {
		log.Println("ไม่พบ Form ยาฉีด:", err)
		return
	}

	// ดึง Form "ยาใช้ทา"
	var formTopical models.Form
	if err := db.Where("form_name = ?", "ยาใช้ทา").First(&formTopical).Error; err != nil {
		log.Println("ไม่พบ Form ยาใช้ทา:", err)
		return
	}

	// ดึง Form "ยาใช้หยด"
	var formDrop models.Form
	if err := db.Where("form_name = ?", "ยาใช้หยด").First(&formDrop).Error; err != nil {
		log.Println("ไม่พบ Form ยาใช้หยด:", err)
		return
	}

	// Mapping Form -> Units
	formUnits := map[string][]string{
		"ยาเม็ด":   {"เม็ด"},
		"แคปซูล":   {"แคปซูล"},
		"ยาน้ำ":    {"ช้อนชา", "ช้อนโต๊ะ", "มิลลิลิตร", "cc"},
		"ยาฉีด":    {"ยูนิต", "cc", "มิลลิลิตร", "มิลลิกรัม", "ไมโครกรัม"},
		"ยาใช้ทา":  {"กรัม", "มิลลิลิตร", "หลอด", "ช้อนชา"},
		"ยาใช้หยด": {"หยด", "มิลลิลิตร", "cc"},
	}

	for formName, unitNames := range formUnits {
		var form models.Form
		if err := db.Where("form_name = ?", formName).First(&form).Error; err != nil {
			log.Println("ไม่พบ Form:", formName, err)
			continue
		}

		var selectedUnits []models.Unit
		if err := db.Where("unit_name IN ?", unitNames).Find(&selectedUnits).Error; err != nil {
			log.Println("ไม่พบ Units ของ", formName, err)
			continue
		}

		// สร้าง Many-to-Many
		if err := db.Model(&form).Association("Units").Replace(&selectedUnits); err != nil {
			log.Println("เชื่อม Form กับ Units ไม่สำเร็จ:", formName, err)
		}
	}

	// seed ข้อมูล Instruction
	instructions := []models.Instruction{
		{InstructionName: "ก่อนอาหาร"},
		{InstructionName: "หลังอาหาร"},
		{InstructionName: "พร้อมอาหาร"},
		{InstructionName: "ก่อนนอน"},
	}
	for _, instruction := range instructions {
		if err := db.FirstOrCreate(&instruction, models.Instruction{InstructionName: instruction.InstructionName}).Error; err != nil {
			log.Println("Seed instruction error:", err)
		}
	}

	// seed ข้อมูล NotiFormat
	notiformats := []models.NotiFormat{
		{FormatName: "เวลาเฉพาะ (Fixed Times)"},
		{FormatName: "ทุกกี่ชั่วโมง (Interval)"},
		{FormatName: "วันเว้นวัน / ทุกกี่วัน (EveryNDays)"},
		{FormatName: "ทานต่อเนื่อง/พักยา (Cycle)"},
	}
	for _, notiformat := range notiformats {
		if err := db.FirstOrCreate(&notiformat, models.NotiFormat{FormatName: notiformat.FormatName}).Error; err != nil {
			log.Println("Seed notiformat error:", err)
		}
	}

	// seed ข้อมูล MedicineInfo
	medicines := []models.MedicineInfo{
		{
			MedName:       "Paracetamol",
			GenericName:   "Acetaminophen",
			Properties:    "บรรเทาอาการปวดลดไข้",
			Strength:      "500mg",
			FormID:        1,
			UnitID:        UintPtr(1),
			InstructionID: UintPtr(2),
			// MedStatus: "active",
		},

		{
			MedName:       "PROBUFEN 400",
			GenericName:   "Ibuprofen",
			Properties:    "บรรเทาอาการปวดและลดไข้ หรือลดการอักเสบ",
			Strength:      "400 mg",
			FormID:        2,
			UnitID:        UintPtr(2),
			InstructionID: UintPtr(2),
		},

		{
			MedName:       "ไบโซลวอน สำหรับเด็ก",
			GenericName:   "bromhexine",
			Properties:    "ละลายเสมหะและบรรเทาอาการไอ",
			Strength:      "4 mg/5 ml",
			FormID:        3,
			UnitID:        UintPtr(3),
			InstructionID: UintPtr(2),
		},

		{
			MedName:     "COUNTERPAIN COOL",
			GenericName: "menthol",
			Properties:  "ใช้ทาบรรเทาอาการปวดกล้ามเนื้อ เนื่องจากการพลิกหรือเคล็ด",
			Strength:    "4%",
			FormID:      5,
		},

		{
			MedName:     "ยาทาแก้ผดผื่นคัน คาลาไมน์",
			GenericName: "calamine+zinc oxide",
			Properties:  "บรรเทาอาการระคายเคืองของผิวหนัง ผื่น ลมพิษในระดับเล็กน้อย",
			Strength:    "(10 G+5 G)/100 ML",
			FormID:      5,
		},

		{
			MedName:       "วินซูลิน-30/70",
			GenericName:   "insulin",
			Properties:    "ใช้สำหรับรักษาโรคเบาหวาน โดยช่วยลดระดับน้ำตาลในเลือด",
			Strength:      "100 iu/1ml",
			FormID:        4,
			UnitID:        UintPtr(7),
			InstructionID: UintPtr(1),
		},

		{
			MedName:     "TEARS NATURALE II",
			GenericName: "hypromellose(hydroxypropyl methylcellulose)+dextran 70",
			Properties:  "รักษาภาวะตาแห้งที่ขาดเมือกและขาดน้ำ",
			Strength:    "(0.3 G+0.1 G)/100 ML",
			FormID:      6,
			UnitID:      UintPtr(12),
		},

		// ...ใส่ยาอีกตามต้องการ
	}
	for _, medicine := range medicines {
		if err := db.Create(&medicine).Error; err != nil {
			log.Println("Seed medicineinfo error:", err)
		}
	}

	// --- Seed Hospitals ---
	hospitals := []models.Hospital{
		{HospitalName: "โรงพยาบาลตัวอย่าง A"},
	}
	for i := range hospitals {
		if err := db.FirstOrCreate(&hospitals[i],
			models.Hospital{HospitalName: hospitals[i].HospitalName},
		).Error; err != nil {
			log.Println("Seed hospital error:", err)
		}
	}
	
	// --- Seed WebAdmins (หมอ / แอดมิน) ---
	admins := []models.WebAdmin{
		{
			Username:  "admin@pillmate.com",
			Password:  "admin1234", // default ครั้งแรก
			FirstName: "System",
			LastName:  "Admin",
			Role:      "superadmin",
		},
	}

	isBcrypt := func(s string) bool {
		return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
	}

	for _, a := range admins {
		// เตรียม hash สำหรับ "กรณีสร้างใหม่"
		hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("bcrypt error:", err)
			continue
		}

		var out models.WebAdmin
		// ใช้ UNIQUE(username) ที่มีอยู่แล้วได้เลย (case-sensitive)
		tx := db.
			Where("username = ?", a.Username).
			Attrs(models.WebAdmin{
				Username:  a.Username,
				Password:  string(hashed), // ถ้าสร้างใหม่จะเก็บเป็น hash ทันที
				FirstName: a.FirstName,
				LastName:  a.LastName,
				Role:      a.Role,
			}).
			FirstOrCreate(&out)
		if tx.Error != nil {
			log.Println("Seed webadmin FirstOrCreate error:", tx.Error)
			continue
		}

		// ถ้ามีอยู่แล้ว แต่ password ยังเป็น plaintext -> รีแฮชทับ (idempotent)
		if !isBcrypt(out.Password) {
			newHashed, err := bcrypt.GenerateFromPassword([]byte(out.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("bcrypt rehash error:", err)
				continue
			}
			if err := db.Model(&out).Update("password", string(newHashed)).Error; err != nil {
				log.Println("Seed webadmin rehash update error:", err)
			}
		}
	}

	// ----- ทำbackend web แล้วค่อยลบ ---------------------------------------------------------------------------

	
	// --- Seed Prescription --- ตอนทำbackend ให้วนยาแต่ละตัวเข้าตารางนะ เพราะแต่ละยาขนาดการกินต่างกัน
	// prescriptions := []models.Prescription{
	// { IDCardNumber: "1101700203452", MedicineInfoID: 1, AmountPerTime: "1", TimesPerDay: "3", HospitalID: 1, DoctorID: 1, AppSyncStatus: false },
	// { IDCardNumber: "1234567890123", MedicineInfoID: 2, AmountPerTime: "1", TimesPerDay: "3", HospitalID: 1, DoctorID: 1, AppSyncStatus: false },
	// }

	// for i := range prescriptions {
	// 	key := models.Prescription{
	// 		IDCardNumber:   prescriptions[i].IDCardNumber,
	// 		MedicineInfoID: prescriptions[i].MedicineInfoID,
	// 		HospitalID:     prescriptions[i].HospitalID,
	// 		DoctorID:       prescriptions[i].DoctorID,
	// 	}
	// 	attrs := models.Prescription{
	// 		AmountPerTime: prescriptions[i].AmountPerTime,
	// 		TimesPerDay:   prescriptions[i].TimesPerDay,
	// 		AppSyncStatus: false,
	// 	}
	// 	if err := db.Where(&key).Attrs(&attrs).FirstOrCreate(&prescriptions[i]).Error; err != nil {
	// 		log.Fatal("seed prescriptions failed: ", err)
	// 	}
	// }

	// --- Seed Appointments ---
	// appointments := []models.Appointment{
	// 	{
	// 		IDCardNumber:    "1101700203452",
	// 		AppointmentDate: time.Date(2025, 10, 10, 0, 0, 0, 0, time.Local),              // date-only
	// 		AppointmentTime: time.Date(1, 1, 1, 9, 30, 0, 0, time.UTC),                    // time-only
	// 		HospitalID:      1,
	// 		DoctorID:        1,
	// 		Note:            "งดอาหารก่อนตรวจ 8 ชั่วโมง",
	// 	},
	// }

	// for i := range appointments {
	// 	key := models.Appointment{
	// 		IDCardNumber:    appointments[i].IDCardNumber,
	// 		AppointmentDate: appointments[i].AppointmentDate,
	// 		AppointmentTime: appointments[i].AppointmentTime,
	// 		HospitalID:      appointments[i].HospitalID,
	// 		DoctorID:        appointments[i].DoctorID,
	// 	}
	// 	attrs := models.Appointment{
	// 		Note: appointments[i].Note,
	// 	}
	// 	if err := db.Where(&key).Attrs(&attrs).FirstOrCreate(&appointments[i]).Error; err != nil {
	// 		log.Println("seed appointments failed:", err)
	// 	}
	// }



	// ------ ลบถึงตรงนี้ ----------------------------------------------------------------------------------------------

	log.Println("Seed ข้อมูลสำเร็จ")
}
