package routes

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/fouradithep/pillmate/db"
	"github.com/fouradithep/pillmate/handlers"
)

func SetupGroupMedicineRoutes(api fiber.Router) {
	// CREATE
	// POST /api/groups
	// Body:
	// {
	//   "group_name": "กินทุกวัน",
	//   "my_medicine_ids": [12, 15, 19]
	// }
	type createGroupBody struct {
		GroupName     string `json:"group_name"`
		MyMedicineIDs []uint `json:"my_medicine_ids"`
	}
	api.Post("/groups", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		var body createGroupBody
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		group, members, err := handlers.CreateGroup(db.DB, patientID, handlers.CreateGroupRequest{
			GroupName:     body.GroupName,
			MyMedicineIDs: body.MyMedicineIDs,
		})
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrInvalidData) || errors.Is(err, gorm.ErrRecordNotFound) {
				status = fiber.StatusBadRequest
			}
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "group created/ensured",
			"group":   group,
			"members": members,
		})
	})

	// LIST GROUPS (with member_count)
	// GET /api/groups
	api.Get("/groups", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		items, err := handlers.GetGroups(db.DB, patientID) // ตอนนี้คืน []GroupWithCount แล้ว
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": items})
	})

	// AVAILABLE (ungrouped medicines)
	// GET /api/groups/available-medicines
	api.Get("/groups/available-medicines", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		items, err := handlers.GetUngroupedMyMedicines(db.DB, patientID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": items})
	})

	// GET ONE (group detail + members)
	// GET /api/groups/:id
	api.Get("/groups/:id", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		item, err := handlers.GetGroup(db.DB, patientID, uint(id))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": item})
	})

	// UPDATE (rename + set full members)
	// PUT /api/groups/:id
	// Body:
	// {
	//   "new_group_name": "เช้าใหม่",      // optional
	//   "my_medicine_ids": [12, 15, 19]    // required: ชุดสมาชิก "สุดท้าย"
	// }
	type updateGroupBody struct {
		NewGroupName  *string `json:"new_group_name"`
		MyMedicineIDs []uint  `json:"my_medicine_ids"`
	}
	api.Put("/groups/:id", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		var body updateGroupBody
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		group, members, err := handlers.UpdateGroup(db.DB, patientID, uint(id), handlers.UpdateGroupRequest{
			NewGroupName:  body.NewGroupName,
			MyMedicineIDs: body.MyMedicineIDs,
		})
		if err != nil {
			status := fiber.StatusInternalServerError
			if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, gorm.ErrInvalidData) {
				status = fiber.StatusBadRequest
			}
			return c.Status(status).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"message": "group updated",
			"group":   group,
			"members": members,
		})
	})

	// DELETE
	// DELETE /api/groups/:id
	api.Delete("/groups/:id", func(c *fiber.Ctx) error {
		patientID, ok := c.Locals("patient_id").(uint)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		id, err := strconv.ParseUint(c.Params("id"), 10, 64)
		if err != nil || id == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		if err := handlers.DeleteGroup(db.DB, patientID, uint(id)); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message": "deleted"})
	})
}
