package dto

import (
	"vet-clinic/models"
)

// PetDto defines a data transfer object for pet.
type PetDto struct {
	Name     string `json:"name" validate:"omitempty,rualpha,max=255"` // Alphabetic characters only (Russian and English).
	Type     string `json:"type" validate:"ruprintascii"`              // Allowed characters: printable ASCII (Russian and English).
	Breed    string `json:"breed" validate:"ruprintascii"`             // Allowed characters: printable ASCII (Russian and English).
	Colour   string `json:"colour" validate:"ruprintascii"`            // Allowed characters: printable ASCII (Russian and English).
	Sex      string `json:"sex" validate:"omitempty,rualpha"`          // Alphabetic characters only (Russian and English).
	ClientID uint   `json:"clientId"`
}

// ToModel creates models.Pet from this DTO.
func (d *PetDto) ToModel() *models.Pet {
	return &models.Pet{
		Name:     d.Name,
		Type:     d.Type,
		Breed:    d.Breed,
		Colour:   d.Colour,
		Sex:      d.Sex,
		ClientID: d.ClientID,
	}
}
