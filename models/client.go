package models

import (
	"time"
	"vet-clinic/repository"
)

// Client defines struct of client data.
type Client struct {
	*BaseModel
	Surname    string    `json:"surname" gorm:"size:255"`
	Name       string    `json:"name" gorm:"size:255"`
	Patronymic string    `json:"patronymic" gorm:"size:255"`
	Sex        string    `json:"sex"`
	BirthDate  time.Time `json:"birthDate"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Info       string    `json:"info"`
}

// TableName returns the table name of client struct and it is used by gorm.
func (*Client) TableName() string {
	return "client_master"
}

// Exist returns true if a given client exits.
func (m *Client) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Client{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns client full matched given client ID.
func (m *Client) Get(rep repository.Repository, id uint) (*Client, error) {
	client := &Client{}
	if err := rep.First(client, id).Error; err != nil {
		return nil, err
	}
	return client, nil
}

// GetAll returns a slice of all clients.
func (m *Client) GetAll(rep repository.Repository) ([]*Client, error) {
	var clients []*Client
	if err := rep.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

// Create persists this client data.
func (m *Client) Create(rep repository.Repository) (*Client, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return tx.Select("surname", "name", "patronymic", "sex",
			"birth_date", "phone", "email", "info").Create(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

// Update updates this client data.
func (m *Client) Update(rep repository.Repository, id uint) (*Client, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		if _, err := m.Exist(tx, id); err != nil {
			return err
		}

		return tx.Model(&Client{}).Where("id = ?", id).
			Select("surname", "name", "patronymic", "sex",
				"birth_date", "phone", "email", "info").Updates(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

// Delete deletes this client data.
func (m *Client) Delete(rep repository.Repository, id uint) (*Client, error) {
	client := &Client{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if client, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Client{}, id).Error
	}); err != nil {
		return nil, err
	}
	return client, nil
}
