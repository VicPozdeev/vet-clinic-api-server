package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type ServiceService struct {
	container container.Container
}

// NewServiceService is constructor.
func NewServiceService(container container.Container) *ServiceService {
	return &ServiceService{container: container}
}

// Get returns service full matched given service ID or service slug.
func (s *ServiceService) Get(id string) (*models.Service, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	service := &models.Service{}
	var err error

	if service, err = service.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch service with ID %s: %v", id, err)
		return nil, err
	}
	return service, nil
}

// GetAll returns a slice of all services.
func (s *ServiceService) GetAll() ([]*models.Service, error) {
	rep := s.container.Repository()
	model := &models.Service{}
	var services []*models.Service
	var err error

	if services, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch services: %v", err)
		return nil, err
	}
	return services, nil
}

// Create persists this service data.
func (s *ServiceService) Create(dto *dto.ServiceDto) (*models.Service, error) {
	rep := s.container.Repository()
	service := dto.ToModel()
	var err error

	if service, err = service.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create service: %v", err)
		return nil, err
	}

	return service, nil
}

// Update updates this service data.
func (s *ServiceService) Update(dto *dto.ServiceDto, id string) (*models.Service, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch service ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	service := dto.ToModel()
	var err error

	if service, err = service.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update service with ID %s: %v", id, err)
		return nil, err
	}
	return service, nil
}

// Delete deletes this service data.
func (s *ServiceService) Delete(id string) (*models.Service, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch service ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	service := &models.Service{}
	var err error

	if service, err = service.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete service: %v", err)
		return nil, err
	}
	return service, nil
}
