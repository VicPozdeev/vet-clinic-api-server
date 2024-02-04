package dto

import (
	"vet-clinic/models"
)

// CategoryDto defines a data transfer object for category.
type CategoryDto struct {
	Name string `json:"name" validate:"required,ruprintascii"` // Allowed characters: printable ASCII (Russian and English).
}

// ToModel creates models.Category from this DTO.
func (d *CategoryDto) ToModel() *models.Category {
	return models.NewCategory(d.Name)
}
