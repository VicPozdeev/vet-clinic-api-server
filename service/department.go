package service

import (
	"errors"
	"github.com/gosimple/slug"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type DepartmentService struct {
	container container.Container
}

// NewDepartmentService is constructor.
func NewDepartmentService(container container.Container) *DepartmentService {
	return &DepartmentService{container: container}
}

// Get returns department full matched given department ID or department slug.
func (s *DepartmentService) Get(param string) (*models.Department, error) {
	rep := s.container.Repository()
	department := &models.Department{}
	var err error

	if util.IsNumeric(param) {
		if department, err = department.Get(rep, util.ConvertToUint(param)); err != nil {
			s.container.Logger().Errorf("Failed to fetch department with ID %s: %v", param, err)
			return nil, err
		}
		return department, nil
	}
	if slug.IsSlug(param) {
		if department, err = department.GetBySlug(rep, param); err != nil {
			s.container.Logger().Errorf("Failed to fetch department with slug %s: %v", param, err)
			return nil, err
		}
		return department, nil
	}
	s.container.Logger().Errorf("Failed to fetch department with: %s", param)
	return nil, errors.New("failed to fetch data")
}

// GetAll returns a slice of all departments.
func (s *DepartmentService) GetAll() ([]*models.Department, error) {
	rep := s.container.Repository()
	model := &models.Department{}
	var departments []*models.Department
	var err error

	if departments, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch departments: %v", err)
		return nil, err
	}
	return departments, nil
}

// Create persists this department data.
func (s *DepartmentService) Create(dto *dto.DepartmentDto) (*models.Department, error) {
	rep := s.container.Repository()
	department := dto.ToModel()
	var err error

	if department, err = department.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create department: %v", err)
		return nil, err
	}

	return department, nil
}

// Update updates this department data.
func (s *DepartmentService) Update(dto *dto.DepartmentDto, id string) (*models.Department, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch department ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	department := dto.ToModel()
	var err error

	if department, err = department.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update department with ID %s: %v", id, err)
		return nil, err
	}
	return department, nil
}

// Delete deletes this department data.
func (s *DepartmentService) Delete(id string) (*models.Department, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch department ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	department := &models.Department{}
	var err error

	if department, err = department.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete department: %v", err)
		return nil, err
	}
	return department, nil
}
