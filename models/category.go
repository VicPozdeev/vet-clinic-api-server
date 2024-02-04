package models

import (
	"vet-clinic/repository"
)

// Category defines struct of category data.
type Category struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique;not null;size:255"`
}

// TableName returns the table name of category struct and it is used by gorm.
func (*Category) TableName() string {
	return "category_master"
}

// NewCategory is constructor.
func NewCategory(name string) *Category {
	return &Category{Name: name}
}

// Exist returns true if a given category exits.
func (m *Category) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Category{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns category full matched given category ID.
func (m *Category) Get(rep repository.Repository, id uint) (*Category, error) {
	category := &Category{}
	if err := rep.First(category, id).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// GetAll returns a slice of all categories.
func (m *Category) GetAll(rep repository.Repository) ([]*Category, error) {
	var categories []*Category
	if err := rep.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Create persists this category data.
func (m *Category) Create(rep repository.Repository) (*Category, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return tx.Select("name").Create(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

// Update updates this category data.
func (m *Category) Update(rep repository.Repository, id uint) (*Category, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		if _, err := m.Exist(tx, id); err != nil {
			return err
		}

		return tx.Model(&Category{}).Where("id = ?", id).
			Select("name").Updates(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

// Delete deletes this category data.
func (m *Category) Delete(rep repository.Repository, id uint) (*Category, error) {
	category := &Category{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if category, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Category{}, id).Error
	}); err != nil {
		return nil, err
	}
	return category, nil
}
