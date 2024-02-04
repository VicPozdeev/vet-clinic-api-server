package service

import (
	"errors"
	"github.com/gosimple/slug"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type UserService struct {
	container container.Container
}

// NewUserService is constructor.
func NewUserService(container container.Container) *UserService {
	return &UserService{container: container}
}

// Get returns user full matched given user ID or user slug.
func (s *UserService) Get(param string) (*models.User, error) {
	rep := s.container.Repository()
	user := &models.User{}
	var err error

	if util.IsNumeric(param) {
		if user, err = user.Get(rep, util.ConvertToUint(param)); err != nil {
			s.container.Logger().Debugf("Failed to fetch user with ID %s: %v", param, err)
			return nil, err
		}
		return user, nil
	}
	if slug.IsSlug(param) {
		if user, err = user.GetBySlug(rep, param); err != nil {
			s.container.Logger().Debugf("Failed to fetch user with slug %s: %v", param, err)
			return nil, err
		}
		return user, nil
	}
	s.container.Logger().Debugf("Failed to fetch user with: %s", param)
	return nil, errors.New("failed to fetch data")
}

// GetAll returns a slice of all users.
func (s *UserService) GetAll() ([]*models.User, error) {
	rep := s.container.Repository()
	model := &models.User{}
	var users []*models.User
	var err error

	if users, err = model.GetAll(rep); err != nil {
		s.container.Logger().Debugf("Failed to fetch users: %v", err)
		return nil, err
	}
	return users, nil
}

// Create persists this user data.
func (s *UserService) Create(dto *dto.UserCreateDto) (*models.User, error) {
	rep := s.container.Repository()
	user := dto.ToModel()
	var err error

	if user, err = user.Create(rep); err != nil {
		s.container.Logger().Debugf("Failed to create user: %v", err)
		return nil, err
	}

	return user, nil
}

// Update updates this user data.
func (s *UserService) Update(dto *dto.UserUpdateDto, id string, owner bool) (*models.User, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Debugf("Failed to fetch user ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	user := dto.ToModel(owner)
	var err error

	if user, err = user.Update(rep, util.ConvertToUint(id), owner); err != nil {
		s.container.Logger().Debugf("Failed to update user with ID %s: %v", id, err)
		return nil, err
	}
	return user, nil
}

// UpdatePassword updates this user data.
func (s *UserService) UpdatePassword(dto *dto.UpdatePasswordDto, id string) error {
	if !util.IsNumeric(id) {
		s.container.Logger().Debugf("Failed to fetch user ID: %s", id)
		return errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	user := &models.User{}

	if err := user.UpdatePassword(rep, util.ConvertToUint(id), dto.OldPassword, dto.NewPassword); err != nil {
		s.container.Logger().Debugf("Failed to update user with ID %s: %v", id, err)
		return err
	}
	return nil
}

// Delete deletes this user data.
func (s *UserService) Delete(id string) (*models.User, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Debugf("Failed to fetch user ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	user := &models.User{}
	var err error

	if user, err = user.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Debugf("Failed to delete user: %v", err)
		return nil, err
	}
	return user, nil
}

// Login authenticates by using login DTO.
func (s *UserService) Login(dto *dto.LoginDto) (*models.User, error) {
	rep := s.container.Repository()
	user := &models.User{}
	var err error

	if user, err = user.Login(rep, dto.Login, dto.Password); err != nil {
		s.container.Logger().Debugf("Failed to login: %v", err)
		return nil, err
	}

	return user, nil
}
