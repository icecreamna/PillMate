package routes

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
)

// เส้นทางข้อมูลอ้างอิงเริ่มต้น (Forms, Units, Instructions)
func SetupInitialDataRoutes(app *fiber.App) {

	// ===== Helpers =====
	parseIDParam := func(c *fiber.Ctx, paramName string) (uint, error) {
		raw := c.Params(paramName)
		id, err := strconv.Atoi(raw)
		if err != nil || id <= 0 {
			return 0, fiber.ErrBadRequest
		}
		return uint(id), nil
	}
	withRelations := func(c *fiber.Ctx) bool {
		// ?with_relations=true
		return c.Query("with_relations") == "true"
	}

	// ===== Forms =====

	// GET /forms
	app.Get("/forms", func(c *fiber.Ctx) error {
		includeRelations := withRelations(c)
		forms, err := handlers.GetForms(db.DB, includeRelations)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(forms)
	})

	// GET /form/:id
	app.Get("/form/:id", func(c *fiber.Ctx) error {
		formID, err := parseIDParam(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "invalid form id"})
		}
		includeRelations := withRelations(c)
		form, err := handlers.GetForm(db.DB, formID, includeRelations)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).
					JSON(fiber.Map{"error": "form not found"})
			}
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(form)
	})

	// GET /forms/:id/units  (Lite: id + unit_name)
	type unitLite struct {
		ID       uint   `json:"id"`
		UnitName string `json:"unit_name"`
	}
	app.Get("/forms/:id/units", func(c *fiber.Ctx) error {
		formID, err := parseIDParam(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "invalid form id"})
		}
		units, err := handlers.GetUnitsByFormID(db.DB, formID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).
					JSON(fiber.Map{"error": "form not found"})
			}
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		result := make([]unitLite, 0, len(units))
		for _, u := range units {
			result = append(result, unitLite{ID: u.ID, UnitName: u.UnitName})
		}
		return c.JSON(fiber.Map{"form_id": formID, "units": result})
	})

	// ===== Units =====

	// GET /units
	app.Get("/units", func(c *fiber.Ctx) error {
		includeRelations := withRelations(c)
		units, err := handlers.GetUnits(db.DB, includeRelations)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(units)
	})

	// GET /unit/:id
	app.Get("/unit/:id", func(c *fiber.Ctx) error {
		unitID, err := parseIDParam(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "invalid unit id"})
		}
		includeRelations := withRelations(c)
		unit, err := handlers.GetUnit(db.DB, unitID, includeRelations)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).
					JSON(fiber.Map{"error": "unit not found"})
			}
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(unit)
	})

	// GET /units/:id/forms (Lite: id + form_name)
	type formLite struct {
		ID       uint   `json:"id"`
		FormName string `json:"form_name"`
	}
	app.Get("/units/:id/forms", func(c *fiber.Ctx) error {
		unitID, err := parseIDParam(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "invalid unit id"})
		}
		forms, err := handlers.GetFormsByUnitID(db.DB, unitID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).
					JSON(fiber.Map{"error": "unit not found"})
			}
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		result := make([]formLite, 0, len(forms))
		for _, f := range forms {
			result = append(result, formLite{ID: f.ID, FormName: f.FormName})
		}
		return c.JSON(fiber.Map{"unit_id": unitID, "forms": result})
	})

	// ===== Instructions =====

	// GET /instructions
	app.Get("/instructions", func(c *fiber.Ctx) error {
		instructions, err := handlers.GetInstructions(db.DB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(instructions)
	})

	// GET /instruction/:id
	app.Get("/instruction/:id", func(c *fiber.Ctx) error {
		instructionID, err := parseIDParam(c, "id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": "invalid instruction id"})
		}
		instruction, err := handlers.GetInstruction(db.DB, instructionID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).
					JSON(fiber.Map{"error": "instruction not found"})
			}
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(instruction)
	})
}
