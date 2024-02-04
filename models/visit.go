package models

import (
	"time"
	"vet-clinic/repository"
)

// Visit defines struct of visit data.
type Visit struct {
	*BaseModel
	DateTime        time.Time `json:"dateTime"`
	Info            string    `json:"info"`
	ClientID        uint      `json:"clientId"`
	Client          *Client   `json:"client" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PetID           uint      `json:"petId"`
	Pet             *Pet      `json:"pet" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DoctorID        uint      `json:"userId"`
	Doctor          *User     `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ServiceID       uint      `json:"serviceId"`
	Service         *Service  `json:"services" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LastUpdatedByID uint      `json:"lastUpdatedById"`
	LastUpdatedBy   *User     `json:"lastUpdatedBy" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName returns the table name of visit struct and it is used by gorm.
func (*Visit) TableName() string {
	return "visit_master"
}

// Exist returns true if a given pet exits.
func (m *Visit) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Visit{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns visit full matched given visit ID.
func (m *Visit) Get(rep repository.Repository, id uint) (*Visit, error) {
	visit := &Visit{}
	if err := rep.Preload("Client").Preload("Pet").Preload("Doctor").
		Preload("LastUpdatedBy").Preload("Service").First(visit, id).Error; err != nil {
		return nil, err
	}
	return visit, nil
}

// GetAll returns a slice of all visits.
func (m *Visit) GetAll(rep repository.Repository) ([]*Visit, error) {
	var visits []*Visit
	if err := rep.Preload("Client").Preload("Pet").Preload("Doctor").
		Preload("LastUpdatedBy").Preload("Service").Find(&visits).Error; err != nil {
		return nil, err
	}
	return visits, nil
}

// Create persists this visit data.
func (m *Visit) Create(rep repository.Repository) (*Visit, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		client := &Client{}
		if _, err := client.Exist(tx, m.ClientID); err != nil {
			return err
		}

		pet := &Pet{}
		if _, err := pet.Exist(tx, m.PetID); err != nil {
			return err
		}

		user := &User{}
		if _, err := user.Exist(tx, m.DoctorID); err != nil {
			return err
		}
		if _, err := user.Exist(tx, m.LastUpdatedByID); err != nil {
			return err
		}

		service := &Service{}
		if _, err := service.Exist(tx, m.ServiceID); err != nil {
			return err
		}

		return tx.Select("date_time", "info", "client_id", "pet_id",
			"doctor_id", "last_updated_by_id", "service_id").Create(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

// Update updates this visit data.
func (m *Visit) Update(rep repository.Repository, id uint) (*Visit, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		if _, err := m.Exist(tx, id); err != nil {
			return err
		}

		client := &Client{}
		if _, err := client.Exist(tx, m.ClientID); err != nil {
			return err
		}

		pet := &Pet{}
		if _, err := pet.Exist(tx, m.PetID); err != nil {
			return err
		}

		user := &User{}
		if _, err := user.Exist(tx, m.DoctorID); err != nil {
			return err
		}
		if _, err := user.Exist(tx, m.LastUpdatedByID); err != nil {
			return err
		}

		service := &Service{}
		if _, err := service.Exist(tx, m.ServiceID); err != nil {
			return err
		}

		return tx.Model(&Visit{}).Where("id = ?", id).
			Select("date_time", "info", "client_id", "pet_id",
				"doctor_id", "last_updated_by_id", "service_id").Updates(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

// Delete deletes this visit data.
func (m *Visit) Delete(rep repository.Repository, id uint) (*Visit, error) {
	visit := &Visit{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if visit, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Visit{}, id).Error
	}); err != nil {
		return nil, err
	}
	return visit, nil
}
