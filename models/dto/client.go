package dto

import (
	"time"
	"vet-clinic/models"
)

// ClientDto defines a data transfer object for client.
type ClientDto struct {
	Surname    string    `json:"surname" validate:"omitempty,rualpha,max=255"`    // Alphabetic characters only (Russian and English).
	Name       string    `json:"name" validate:"omitempty,rualpha,max=255"`       // Alphabetic characters only (Russian and English).
	Patronymic string    `json:"patronymic" validate:"omitempty,rualpha,max=255"` // Alphabetic characters only (Russian and English).
	Sex        string    `json:"sex" validate:"omitempty,rualpha"`                // Alphabetic characters only (Russian and English).
	BirthDate  time.Time `json:"birthDate" format:"date"`                         // Date only
	Phone      string    `json:"phone" validate:"e164" example:"+79876543210"`    // E.164 phone number string.
	Email      string    `json:"email" validate:"email" example:"mail@mail.com"`  // E-mail string.
	Info       string    `json:"info" validate:"ruprintascii"`                    // Allowed characters: printable ASCII (Russian and English).
}

// ToModel creates models.Client from this DTO.
func (d *ClientDto) ToModel() *models.Client {
	return &models.Client{
		Surname:    d.Surname,
		Name:       d.Name,
		Patronymic: d.Patronymic,
		Sex:        d.Sex,
		BirthDate:  d.BirthDate,
		Phone:      d.Phone,
		Email:      d.Email,
		Info:       d.Info,
	}
}
