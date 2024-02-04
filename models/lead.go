package models

import (
	"vet-clinic/repository"
)

// Lead defines struct of lead data.
type Lead struct {
	*BaseModel
	Name            string `json:"name" gorm:"size:255"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Comment         string `json:"comment"`
	Type            string `json:"type"`   // in clinic, online, callback
	Status          string `json:"status"` // open, in_progress, closed, rejected
	DoctorID        uint   `json:"doctorId"`
	Doctor          *User  `json:"doctor" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LastUpdatedByID uint   `json:"lastUpdatedById"`
	LastUpdatedBy   *User  `json:"lastUpdatedBy" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName returns the table name of lead struct and it is used by gorm.
func (*Lead) TableName() string {
	return "lead_master"
}

// Exist returns true if a given lead exits.
func (m *Lead) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Lead{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns lead full matched given lead ID.
func (m *Lead) Get(rep repository.Repository, id uint) (*Lead, error) {
	lead := &Lead{}
	if err := rep.Preload("Doctor").Preload("LastUpdatedBy").First(lead, id).Error; err != nil {
		return nil, err
	}
	return lead, nil
}

// GetAll returns a slice of all leads.
func (m *Lead) GetAll(rep repository.Repository) ([]*Lead, error) {
	var leads []*Lead
	if err := rep.Preload("Doctor").Preload("LastUpdatedBy").Find(&leads).Error; err != nil {
		return nil, err
	}
	return leads, nil
}

// Create persists this lead data.
func (m *Lead) Create(rep repository.Repository) (*Lead, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txCreateLead(tx, m)
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

func txCreateLead(tx repository.Repository, m *Lead) error {
	user := &User{}
	if _, err := user.Exist(tx, m.DoctorID); err != nil {
		return err
	}

	m.Status = "open"
	return tx.Select("name", "phone", "email", "comment", "type", "status", "doctor_id").Create(m).Error
}

// Update updates this lead data.
func (m *Lead) Update(rep repository.Repository, id uint) (*Lead, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txUpdateLead(tx, m, id)
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

func txUpdateLead(tx repository.Repository, m *Lead, id uint) error {
	if _, err := m.Exist(tx, id); err != nil {
		return err
	}

	user := &User{}
	if _, err := user.Exist(tx, m.DoctorID); err != nil {
		return err
	}
	if _, err := user.Exist(tx, m.LastUpdatedByID); err != nil {
		return err
	}

	return tx.Model(&Lead{}).Where("id = ?", id).
		Select("name", "phone", "email", "comment", "type",
			"status", "doctor_id", "last_updated_by_id").Updates(m).Error
}

// Delete deletes this lead data.
func (m *Lead) Delete(rep repository.Repository, id uint) (*Lead, error) {
	lead := &Lead{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if lead, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Lead{}, id).Error
	}); err != nil {
		return nil, err
	}
	return lead, nil
}
