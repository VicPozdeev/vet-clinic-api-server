package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type VisitService struct {
	container container.Container
}

// NewVisitService is constructor.
func NewVisitService(container container.Container) *VisitService {
	return &VisitService{container: container}
}

// Get returns visit full matched given visit ID or visit slug.
func (s *VisitService) Get(id string) (*models.Visit, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	visit := &models.Visit{}
	var err error

	if visit, err = visit.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch visit with ID %s: %v", id, err)
		return nil, err
	}
	return visit, nil
}

// GetAll returns a slice of all visits.
func (s *VisitService) GetAll() ([]*models.Visit, error) {
	rep := s.container.Repository()
	model := &models.Visit{}
	var visits []*models.Visit
	var err error

	if visits, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch visits: %v", err)
		return nil, err
	}
	return visits, nil
}

// Create persists this visit data.
func (s *VisitService) Create(dto *dto.VisitDto) (*models.Visit, error) {
	rep := s.container.Repository()
	visit := dto.ToModel()
	var err error

	if visit, err = visit.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create visit: %v", err)
		return nil, err
	}

	return visit, nil
}

// Update updates this visit data.
func (s *VisitService) Update(dto *dto.VisitDto, id string) (*models.Visit, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch visit ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	visit := dto.ToModel()
	var err error

	if visit, err = visit.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update visit with ID %s: %v", id, err)
		return nil, err
	}
	return visit, nil
}

// Delete deletes this visit data.
func (s *VisitService) Delete(id string) (*models.Visit, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch visit ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	visit := &models.Visit{}
	var err error

	if visit, err = visit.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete visit: %v", err)
		return nil, err
	}
	return visit, nil
}
