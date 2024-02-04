package service

import (
	"errors"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/util"
)

type PetService struct {
	container container.Container
}

// NewPetService is constructor.
func NewPetService(container container.Container) *PetService {
	return &PetService{container: container}
}

// Get returns pet full matched given pet ID or pet slug.
func (s *PetService) Get(id string) (*models.Pet, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch client ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	pet := &models.Pet{}
	var err error

	if pet, err = pet.Get(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to fetch pet with ID %s: %v", id, err)
		return nil, err
	}
	return pet, nil
}

// GetAll returns a slice of all pets.
func (s *PetService) GetAll() ([]*models.Pet, error) {
	rep := s.container.Repository()
	model := &models.Pet{}
	var pets []*models.Pet
	var err error

	if pets, err = model.GetAll(rep); err != nil {
		s.container.Logger().Errorf("Failed to fetch pets: %v", err)
		return nil, err
	}
	return pets, nil
}

// Create persists this pet data.
func (s *PetService) Create(dto *dto.PetDto) (*models.Pet, error) {
	rep := s.container.Repository()
	pet := dto.ToModel()
	var err error

	if pet, err = pet.Create(rep); err != nil {
		s.container.Logger().Errorf("Failed to create pet: %v", err)
		return nil, err
	}

	return pet, nil
}

// Update updates this pet data.
func (s *PetService) Update(dto *dto.PetDto, id string) (*models.Pet, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch pet ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	pet := dto.ToModel()
	var err error

	if pet, err = pet.Update(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to update pet with ID %s: %v", id, err)
		return nil, err
	}
	return pet, nil
}

// Delete deletes this pet data.
func (s *PetService) Delete(id string) (*models.Pet, error) {
	if !util.IsNumeric(id) {
		s.container.Logger().Errorf("Failed to fetch pet ID: %s", id)
		return nil, errors.New("failed to fetch data")
	}

	rep := s.container.Repository()
	pet := &models.Pet{}
	var err error

	if pet, err = pet.Delete(rep, util.ConvertToUint(id)); err != nil {
		s.container.Logger().Errorf("Failed to delete pet: %v", err)
		return nil, err
	}
	return pet, nil
}
