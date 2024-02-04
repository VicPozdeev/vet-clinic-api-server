package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type CategoryService struct {
	container container.Container
}

// NewCategoryService is constructor.
func NewCategoryService(container container.Container) *CategoryService {
	return &CategoryService{container: container}
}

// Get returns category full matched given category ID.
func (s *CategoryService) Get(id string) (*models.Category, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch category ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	category := &models.Category{}
	var err error

	if category, err = category.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch category with ID %s: %v", id, err)
		return nil, err
	}
	return category, nil
}

// GetAll returns a slice of all categories.
func (s *CategoryService) GetAll() ([]*models.Category, error) {
	rep := s.container.Repository()
	model := &models.Category{}
	var categories []*models.Category
	var err error

	if categories, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch categories: %v", err)
		return nil, err
	}
	return categories, nil
}

// Create persists this category data.
func (s *CategoryService) Create(dto *dto.CategoryDto) (*models.Category, error) {
	rep := s.container.Repository()
	category := &models.Category{}
	var err error

	if category, err = dto.ToModel().Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create category: %v", err)
		return nil, err
	}
	return category, nil
}

// Update updates this category data.
func (s *CategoryService) Update(dto *dto.CategoryDto, id string) (*models.Category, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch category ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	category := &models.Category{}
	var err error

	if category, err = dto.ToModel().Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update category with ID %s: %v", id, err)
		return nil, err
	}
	return category, nil
}

// Delete deletes this category data.
func (s *CategoryService) Delete(id string) (*models.Category, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch category ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	category := &models.Category{}
	var err error

	if category, err = category.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete category: %v", err)
		return nil, err
	}
	return category, nil
}
