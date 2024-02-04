package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type ClientService struct {
	container container.Container
}

// NewClientService is constructor.
func NewClientService(container container.Container) *ClientService {
	return &ClientService{container: container}
}

// Get returns client full matched given client ID.
func (s *ClientService) Get(id string) (*models.Client, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	client := &models.Client{}
	var err error

	if client, err = client.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch client with ID %s: %v", id, err)
		return nil, err
	}
	return client, nil
}

// GetAll returns a slice of all clients.
func (s *ClientService) GetAll() ([]*models.Client, error) {
	rep := s.container.Repository()
	model := &models.Client{}
	var clients []*models.Client
	var err error

	if clients, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch clients: %v", err)
		return nil, err
	}
	return clients, nil
}

// Create persists this client data.
func (s *ClientService) Create(dto *dto.ClientDto) (*models.Client, error) {
	rep := s.container.Repository()
	client := dto.ToModel()
	var err error

	if client, err = client.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create client: %v", err)
		return nil, err
	}
	return client, nil
}

// Update updates this client data.
func (s *ClientService) Update(dto *dto.ClientDto, id string) (*models.Client, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	client := dto.ToModel()
	var err error

	if client, err = client.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update client with ID %s: %v", id, err)
		return nil, err
	}
	return client, nil
}

// Delete deletes this client data.
func (s *ClientService) Delete(id string) (*models.Client, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	client := &models.Client{}
	var err error

	if client, err = client.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete client: %v", err)
		return nil, err
	}
	return client, nil
}
