package dto

import (
	"vet-clinic/models"
)

// LeadDto defines a data transfer object for lead.
type LeadDto struct {
	Name            string `json:"name" validate:"omitempty,rualpha,max=255"`                // Alphabetic characters only (Russian and English).
	Phone           string `json:"phone" validate:"omitempty,e164" example:"+79876543210"`   // E.164 phone number string.
	Email           string `json:"email" validate:"omitempty,email" example:"mail@mail.com"` // E-mail string.
	Comment         string `json:"comment" validate:"ruprintascii"`                          // Allowed characters: printable ASCII (Russian and English).
	Type            string `json:"type" validate:"ruprintascii"`                             // Allowed characters: printable ASCII (Russian and English).
	Status          string `json:"status" validate:"ruprintascii"`                           // Allowed characters: printable ASCII (Russian and English).
	DoctorID        uint   `json:"doctorId"`
	LastUpdatedByID uint   `json:"-"`
}

// ToModel creates models.Lead from this DTO.
func (d *LeadDto) ToModel() *models.Lead {
	return &models.Lead{
		Name:            d.Name,
		Phone:           d.Phone,
		Email:           d.Email,
		Comment:         d.Comment,
		Type:            d.Type,
		Status:          d.Status,
		DoctorID:        d.DoctorID,
		LastUpdatedByID: d.LastUpdatedByID,
	}
}
