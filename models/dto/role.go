package dto

import (
	"vet-clinic/models"
)

// RoleDto defines a data transfer object for role.
type RoleDto struct {
	Name string `json:"name" validate:"required,rualpha"` // Alphabetic characters only (Russian and English).
}

// ToModel creates models.Role from this DTO.
func (d *RoleDto) ToModel() *models.Role {
	return models.NewRole(d.Name)
}
