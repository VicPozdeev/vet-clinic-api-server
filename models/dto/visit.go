package dto

import (
	"time"
	"vet-clinic/models"
)

// VisitDto defines a data transfer object for visit.
type VisitDto struct {
	DateTime        time.Time `json:"dateTime" format:"date-time"`  // Date and time.
	Info            string    `json:"info" validate:"ruprintascii"` // Allowed characters: printable ASCII (Russian and English).
	ClientID        uint      `json:"clientId"`
	PetID           uint      `json:"petId"`
	DoctorID        uint      `json:"doctorId"`
	ServiceID       uint      `json:"serviceId"`
	LastUpdatedByID uint      `json:"-"`
}

// ToModel creates models.Visit from this DTO.
func (d *VisitDto) ToModel() *models.Visit {
	return &models.Visit{
		DateTime:        d.DateTime,
		Info:            d.Info,
		ClientID:        d.ClientID,
		PetID:           d.PetID,
		DoctorID:        d.DoctorID,
		ServiceID:       d.ServiceID,
		LastUpdatedByID: d.LastUpdatedByID,
	}
}
