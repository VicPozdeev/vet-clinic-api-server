package dto

import (
	"vet-clinic/models"
)

// ServiceDto defines a data transfer object for service.
type ServiceDto struct {
	Name       string  `json:"name" validate:"required,ruprintascii,max=255"` // Allowed characters: printable ASCII (Russian and English).
	Price      float64 `json:"price" format:"float64"`
	CategoryID uint    `json:"categoryId"`
}

// ToModel creates models.Service from this DTO.
func (d *ServiceDto) ToModel() *models.Service {
	return &models.Service{
		Name:       d.Name,
		Price:      d.Price,
		CategoryID: d.CategoryID,
	}
}
