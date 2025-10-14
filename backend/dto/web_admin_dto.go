package dto

import (
	"time"

	"github.com/fouradithep/pillmate/models"
)

type SuperAdminLoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type DoctorLoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateDoctorDTO struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UpdateDoctorDTO struct {
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

/* ====== Response DTO (ป้องกันการรั่ว password) ====== */

type WebAdminDTO struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

/* ====== Helpers: model -> DTO ====== */

func ToWebAdminDTO(m *models.WebAdmin) WebAdminDTO {
	return WebAdminDTO{
		ID:        m.ID,
		Username:  m.Username,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Role:      m.Role,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToWebAdminDTOs(list []models.WebAdmin) []WebAdminDTO {
	out := make([]WebAdminDTO, 0, len(list))
	for i := range list {
		out = append(out, ToWebAdminDTO(&list[i]))
	}
	return out
}

type AdminResetPasswordDTO struct {
	NewPassword string `json:"new_password" validate:"required"` // เอา min=6 ออก
}