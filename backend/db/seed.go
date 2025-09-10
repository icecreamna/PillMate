package db

import (
	"log"

	"github.com/fouradithep/pillmate/models"
	"gorm.io/gorm"
)



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
		{FormatName: "วันเว้นวัน / ทุกกี่วัน"},
		{FormatName: "รายสัปดาห์ (Weekly)"},
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
        MedName: "TYLENOL 500",
        GenericName: "Paracetamol",
        Properties: "บรรเทาอาการปวดลดไข้",
        Strength: "500mg",
        FormID: 1,
        UnitID: 1, 
        InstructionID: 2, 
		// MedStatus: "active",
    },

	{
        MedName: "PROBUFEN 400",
        GenericName: "Ibuprofen",
        Properties: "บรรเทาอาการปวดและลดไข้ หรือลดการอักเสบ",
        Strength: "400 mg",
        FormID: 2,
        UnitID: 2, 
        InstructionID: 2, 
    },

	{
        MedName: "ไบโซลวอน สำหรับเด็ก",
        GenericName: "bromhexine",
        Properties: "ละลายเสมหะและบรรเทาอาการไอ",
        Strength: "4 mg/5 ml",
        FormID: 3,
        UnitID: 3, 
        InstructionID: 2, 
    },

	{
        MedName: "COUNTERPAIN COOL",
        GenericName: "menthol",
        Properties: "ใช้ทาบรรเทาอาการปวดกล้ามเนื้อ เนื่องจากการพลิกหรือเคล็ด",
        Strength: "4%",
        FormID: 5, 
    },

	{
        MedName: "ยาทาแก้ผดผื่นคัน คาลาไมน์",
        GenericName: "calamine+zinc oxide",
        Properties: "บรรเทาอาการระคายเคืองของผิวหนัง ผื่น ลมพิษในระดับเล็กน้อย",
        Strength: "(10 G+5 G)/100 ML",
        FormID: 5, 
    },

	{
        MedName: "วินซูลิน-30/70",
        GenericName: "insulin",
        Properties: "ใช้สำหรับรักษาโรคเบาหวาน โดยช่วยลดระดับน้ำตาลในเลือด",
        Strength: "100 iu/1ml",
        FormID: 4,
        UnitID: 7, 
        InstructionID: 1, 
    },

	{
        MedName: "TEARS NATURALE II",
        GenericName: "hypromellose(hydroxypropyl methylcellulose)+dextran 70",
        Properties: "รักษาภาวะตาแห้งที่ขาดเมือกและขาดน้ำ",
        Strength: "(0.3 G+0.1 G)/100 ML",
        FormID: 6,
        UnitID: 12, 
    },
    
    // ...ใส่ยาอีกตามต้องการ
	}
	for _, medicine := range medicines {
    if err := db.Create(&medicine).Error; err != nil {
        log.Println("Seed medicineinfo error:", err)
    }
	}
	log.Println("Seed ข้อมูลสำเร็จ")
}