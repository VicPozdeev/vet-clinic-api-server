package dto

import (
	"vet-clinic/models"
)

// DepartmentDto defines a data transfer object for department.
type DepartmentDto struct {
	Name     string `json:"name" validate:"required,rualpha,max=255"` // Alphabetic characters only (Russian and English).
	Services []uint `json:"services"`
}

// ToModel creates models.Department from this DTO.
func (d *DepartmentDto) ToModel() *models.Department {
	var services []*models.Service
	for _, id := range d.Services {
		service := &models.Service{BaseModel: &models.BaseModel{ID: id}}
		services = append(services, service)
	}
	return &models.Department{
		Name:     d.Name,
		Services: services,
	}
}
