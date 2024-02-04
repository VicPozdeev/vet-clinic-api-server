package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindPetByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindPetByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindPetByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllPets_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 1)
	assert.NoError(t, err)
}

func TestCreatePet_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	petDto := createPetForCreate()
	_, err := s.Create(petDto)

	result, _ := s.Get("2")
	result.BaseModel = nil
	result.Client = nil

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, petDto.ToModel(), result)
}

func TestCreatePet_NotClient(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.Create(createPetForNotClient())

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestUpdatePet_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	petDto := createPetForCreate()
	_, err := s.Update(petDto, "1")

	result, _ := s.Get("1")
	result.BaseModel = nil
	result.Client = nil

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, petDto.ToModel(), result)
}

func TestUpdatePet_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.Update(createPetForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeletePet_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeletePet_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewPetService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createPetForCreate() *dto.PetDto {
	return &dto.PetDto{
		Name:     "Шарик",
		Type:     "Пёс",
		Breed:    "Дворняга",
		Colour:   "Коричневый",
		Sex:      "Самец",
		ClientID: 1,
	}
}

func createPetForNotClient() *dto.PetDto {
	return &dto.PetDto{
		Name:   "Шарик",
		Type:   "Пёс",
		Breed:  "Дворняга",
		Colour: "Коричневый",
		Sex:    "Самец",
	}
}
