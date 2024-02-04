package models

import (
	"vet-clinic/repository"
)

// Pet defines struct of pet data.
type Pet struct {
	*BaseModel
	Name     string  `json:"name" gorm:"size:255"`
	Type     string  `json:"type" gorm:"size:255"`
	Breed    string  `json:"breed" gorm:"size:255"`
	Colour   string  `json:"colour" gorm:"size:255"`
	Sex      string  `json:"sex" gorm:"size:255"`
	ClientID uint    `json:"clientId"`
	Client   *Client `json:"client" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName returns the table name of pet struct and it is used by gorm.
func (*Pet) TableName() string {
	return "pet_master"
}

// Exist returns true if a given pet exits.
func (m *Pet) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Pet{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns pet full matched given pet ID.
func (m *Pet) Get(rep repository.Repository, id uint) (*Pet, error) {
	pet := &Pet{}
	if err := rep.Preload("Client").First(pet, id).Error; err != nil {
		return nil, err
	}
	return pet, nil
}

// GetAll returns a slice of all pets.
func (m *Pet) GetAll(rep repository.Repository) ([]*Pet, error) {
	var pets []*Pet

	if err := rep.Preload("Client").Find(&pets).Error; err != nil {
		return nil, err
	}
	return pets, nil
}

// Create persists this pet data.
func (m *Pet) Create(rep repository.Repository) (*Pet, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		client := &Client{}
		if _, err := client.Exist(tx, m.ClientID); err != nil {
			return err
		}

		return tx.Select("name", "type", "breed", "colour", "sex", "client_id").Create(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

// Update updates this pet data.
func (m *Pet) Update(rep repository.Repository, id uint) (*Pet, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		if _, err := m.Exist(tx, id); err != nil {
			return err
		}

		client := &Client{}
		if _, err := client.Exist(tx, m.ClientID); err != nil {
			return err
		}

		return tx.Model(&Pet{}).Where("id = ?", id).
			Select("name", "type", "breed", "colour", "sex", "client_id").Updates(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

// Delete deletes this pet data.
func (m *Pet) Delete(rep repository.Repository, id uint) (*Pet, error) {
	pet := &Pet{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if pet, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Pet{}, id).Error
	}); err != nil {
		return nil, err
	}
	return pet, nil
}
