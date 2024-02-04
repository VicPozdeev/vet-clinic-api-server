package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type LeadService struct {
	container container.Container
}

// NewLeadService is constructor.
func NewLeadService(container container.Container) *LeadService {
	return &LeadService{container: container}
}

// Get returns lead full matched given lead ID or lead slug.
func (s *LeadService) Get(id string) (*models.Lead, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	lead := &models.Lead{}
	var err error

	if lead, err = lead.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch lead with ID %s: %v", id, err)
		return nil, err
	}
	return lead, nil
}

// GetAll returns a slice of all leads.
func (s *LeadService) GetAll() ([]*models.Lead, error) {
	rep := s.container.Repository()
	model := &models.Lead{}
	var leads []*models.Lead
	var err error

	if leads, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch leads: %v", err)
		return nil, err
	}
	return leads, nil
}

// Create persists this lead data.
func (s *LeadService) Create(dto *dto.LeadDto) (*models.Lead, error) {
	rep := s.container.Repository()
	lead := dto.ToModel()
	var err error

	if lead, err = lead.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create lead: %v", err)
		return nil, err
	}

	return lead, nil
}

// Update updates this lead data.
func (s *LeadService) Update(dto *dto.LeadDto, id string) (*models.Lead, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch lead ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	lead := dto.ToModel()
	var err error

	if lead, err = lead.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update lead with ID %s: %v", id, err)
		return nil, err
	}
	return lead, nil
}

// Delete deletes this lead data.
func (s *LeadService) Delete(id string) (*models.Lead, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch lead ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	lead := &models.Lead{}
	var err error

	if lead, err = lead.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete lead: %v", err)
		return nil, err
	}
	return lead, nil
}
