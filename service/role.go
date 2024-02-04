package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type RoleService struct {
	container container.Container
}

// NewRoleService is constructor.
func NewRoleService(container container.Container) *RoleService {
	return &RoleService{container: container}
}

// Get returns role full matched given role ID.
func (s *RoleService) Get(id string) (*models.Role, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch role ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	role := &models.Role{}
	var err error

	if role, err = role.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch role with ID %s: %v", id, err)
		return nil, err
	}
	return role, nil
}

// GetAll returns a slice of all roles.
func (s *RoleService) GetAll() ([]*models.Role, error) {
	rep := s.container.Repository()
	model := &models.Role{}
	var roles []*models.Role
	var err error

	if roles, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch roles: %v", err)
		return nil, err
	}
	return roles, nil
}

// Create persists this role data.
func (s *RoleService) Create(dto *dto.RoleDto) (*models.Role, error) {
	rep := s.container.Repository()
	role := dto.ToModel()
	var err error

	if role, err = role.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create role: %v", err)
		return nil, err
	}
	return role, nil
}

// Update updates this role data.
func (s *RoleService) Update(dto *dto.RoleDto, id string) (*models.Role, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch role ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	role := dto.ToModel()
	var err error

	if role, err = role.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update role with ID %s: %v", id, err)
		return nil, err
	}
	return role, nil
}

// Delete deletes this role data.
func (s *RoleService) Delete(id string) (*models.Role, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch role ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	role := &models.Role{}
	var err error

	if role, err = role.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete role: %v", err)
		return nil, err
	}
	return role, nil
}
