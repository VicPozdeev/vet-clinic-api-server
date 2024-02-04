package models

import (
	"vet-clinic/repository"
)

// Service defines struct of service data.
type Service struct {
	*BaseModel
	Name        string        `json:"name" gorm:"unique;not null;size:255"`
	Price       float64       `json:"price" gorm:"not null"`
	CategoryID  uint          `json:"categoryId"`
	Category    *Category     `json:"category"`
	Users       []*User       `json:"users" gorm:"many2many:users_services;"`
	Departments []*Department `json:"departments" gorm:"many2many:departments_services;"`
}

// TableName returns the table name of service struct and it is used by gorm.
func (*Service) TableName() string {
	return "service_master"
}

// NewService is constructor.
func NewService() *Service {
	return &Service{}
}

// Exist returns true if a given service exits.
func (m *Service) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Service{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns service full matched given service ID.
func (m *Service) Get(rep repository.Repository, id uint) (*Service, error) {
	service := &Service{}
	if err := rep.Preload("Category").Preload("Users").Preload("Departments").
		First(service, id).Error; err != nil {
		return nil, err
	}
	return service, nil
}

// GetAll returns a slice of all services.
func (m *Service) GetAll(rep repository.Repository) ([]*Service, error) {
	var services []*Service

	if err := rep.Preload("Category").Preload("Users").
		Preload("Departments").Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

// Create persists this service data.
func (m *Service) Create(rep repository.Repository) (*Service, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txCreateService(tx, m)
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

func txCreateService(tx repository.Repository, m *Service) error {
	category := &Category{}
	if _, err := category.Exist(tx, m.CategoryID); err != nil {
		return err
	}

	return tx.Select("name", "price", "category_id").Create(m).Error
}

// Update updates this service data.
func (m *Service) Update(rep repository.Repository, id uint) (*Service, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txUpdateService(tx, m, id)
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

func txUpdateService(tx repository.Repository, m *Service, id uint) error {
	if _, err := m.Exist(tx, id); err != nil {
		return err
	}

	category := &Category{}
	if _, err := category.Exist(tx, m.CategoryID); err != nil {
		return err
	}

	return tx.Model(&Service{}).Where("id = ?", id).
		Select("name", "price", "category_id").Updates(m).Error
}

// Delete deletes this service data.
func (m *Service) Delete(rep repository.Repository, id uint) (*Service, error) {
	service := &Service{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if service, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Service{}, id).Error
	}); err != nil {
		return nil, err
	}
	return service, nil
}
