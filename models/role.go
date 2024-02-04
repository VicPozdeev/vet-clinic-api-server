package models

import (
	"vet-clinic/repository"
)

// Role defines struct of role data.
type Role struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique;not null;size:255"`
}

// TableName returns the table name of role struct, and it is used by gorm.
func (*Role) TableName() string {
	return "role_master"
}

// NewRole is constructor.
func NewRole(name string) *Role {
	return &Role{Name: name}
}

// Exist returns true if a given role exits.
func (m *Role) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Role{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns role full matched given role ID.
func (m *Role) Get(rep repository.Repository, id uint) (*Role, error) {
	role := &Role{}
	if err := rep.First(role, id).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// GetAll returns a slice of all roles.
func (m *Role) GetAll(rep repository.Repository) ([]*Role, error) {
	var roles []*Role
	if err := rep.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Create persists this role data.
func (m *Role) Create(rep repository.Repository) (*Role, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return tx.Select("name").Create(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

// Update updates this role data.
func (m *Role) Update(rep repository.Repository, id uint) (*Role, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		if err := tx.First(&Role{}, "id = ?", id).Error; err != nil {
			return err
		}

		return tx.Model(&Role{}).Where("id = ?", id).
			Select("name").Updates(m).Error
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

// Delete deletes this role data.
func (m *Role) Delete(rep repository.Repository, id uint) (*Role, error) {
	role := &Role{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		if err := tx.First(role, "id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&Role{}, id).Error
	}); err != nil {
		return nil, err
	}
	return role, nil
}
