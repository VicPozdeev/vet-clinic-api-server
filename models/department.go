package models

import (
	"github.com/gosimple/slug"
	"vet-clinic/repository"
)

// Department defines struct of department data.
type Department struct {
	*BaseModel
	Name     string     `json:"name" gorm:"unique;not null;size:255"`
	Slug     string     `json:"slug"`
	Users    []*User    `json:"users" gorm:"many2many:user_departments;"`
	Services []*Service `json:"services" gorm:"many2many:departments_services;"`
}

// TableName returns the table name of department struct and it is used by gorm.
func (*Department) TableName() string {
	return "department_master"
}

// NewDepartment is constructor.
func NewDepartment(name string) *Department {
	return &Department{Name: name}
}

// Exist returns true if a given department exits.
func (m *Department) Exist(rep repository.Repository, id uint) (bool, error) {
	if err := rep.First(&Department{}, id).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Get returns department full matched given department ID.
func (m *Department) Get(rep repository.Repository, id uint) (*Department, error) {
	department := &Department{}
	if err := rep.Preload("Users").Preload("Services").
		First(department, id).Error; err != nil {
		return nil, err
	}
	return department, nil
}

// GetBySlug returns department full matched given department slug.
func (m *Department) GetBySlug(rep repository.Repository, slug string) (*Department, error) {
	department := &Department{}
	if err := rep.Preload("Users").Preload("Services").
		First(department, "slug = ?", slug).Error; err != nil {
		return nil, err
	}

	return department, nil
}

// GetAll returns a slice of all departments.
func (m *Department) GetAll(rep repository.Repository) ([]*Department, error) {
	var departments []*Department
	if err := rep.Preload("Users").Preload("Services").
		Find(&departments).Error; err != nil {
		return nil, err
	}
	return departments, nil
}

// Create persists this department data.
func (m *Department) Create(rep repository.Repository) (*Department, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txCreateDepartment(tx, m)
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, m.ID)
}

func txCreateDepartment(tx repository.Repository, m *Department) error {
	i := 0
	for _, service := range m.Services {
		if ok, _ := service.Exist(tx, service.ID); ok {
			m.Services[i] = service
			i++
		}
	}
	for j := i; j < len(m.Services); j++ {
		m.Services[j] = nil
	}
	m.Services = m.Services[:i]

	makeDepartmentSlug(m)

	return tx.Select("name", "slug", "Services").Create(m).Error
}

// Update updates this department data.
func (m *Department) Update(rep repository.Repository, id uint) (*Department, error) {
	if err := rep.Transaction(func(tx repository.Repository) error {
		return txUpdateDepartment(tx, m, id)
	}); err != nil {
		return nil, err
	}
	return m.Get(rep, id)
}

func txUpdateDepartment(tx repository.Repository, m *Department, id uint) error {
	if _, err := m.Exist(tx, id); err != nil {
		return err
	}

	i := 0
	for _, service := range m.Services {
		if ok, _ := service.Exist(tx, service.ID); ok {
			m.Services[i] = service
			i++
		}
	}
	for j := i; j < len(m.Services); j++ {
		m.Services[j] = nil
	}
	m.Services = m.Services[:i]

	if err := tx.Model(&Department{BaseModel: &BaseModel{ID: id}}).
		Association("Services").Replace(m.Services); err != nil {
		return err
	}

	makeDepartmentSlug(m)

	return tx.Model(&Department{}).Where("id = ?", id).
		Select("name", "slug").Updates(m).Error
}

// Delete deletes this department data.
func (m *Department) Delete(rep repository.Repository, id uint) (*Department, error) {
	department := &Department{}
	if err := rep.Transaction(func(tx repository.Repository) error {
		var err error
		if department, err = m.Get(tx, id); err != nil {
			return err
		}
		return tx.Delete(&Department{}, id).Error
	}); err != nil {
		return nil, err
	}
	return department, nil
}

func makeDepartmentSlug(department *Department) {
	slug.MaxLength = 40
	slug.EnableSmartTruncate = false
	slug.CustomSub = map[string]string{
		"ь": "",
		"Ь": "",
		"ъ": "",
		"Ъ": "",
	}
	department.Slug = slug.Make(department.Name)
}
