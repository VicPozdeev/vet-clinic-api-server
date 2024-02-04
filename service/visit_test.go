package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindVisitByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindVisitByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindVisitByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllVisits_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 1)
	assert.NoError(t, err)
}

func TestCreateVisit_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	visitDto := createVisitForCreate()
	_, err := s.Create(visitDto)

	result, _ := s.Get("2")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, visitDto.DateTime, result.DateTime.Local())
	assert.Equal(t, visitDto.Info, result.Info)
	assert.Equal(t, visitDto.ClientID, result.ClientID)
	assert.Equal(t, visitDto.PetID, result.PetID)
	assert.Equal(t, visitDto.DoctorID, result.DoctorID)
	assert.Equal(t, visitDto.ServiceID, result.ServiceID)
	assert.Equal(t, visitDto.LastUpdatedByID, result.LastUpdatedByID)
}

func TestUpdateVisit_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	visitDto := createVisitForCreate()
	_, err := s.Update(visitDto, "1")

	result, _ := s.Get("1")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, visitDto.DateTime, result.DateTime.Local())
	assert.Equal(t, visitDto.Info, result.Info)
	assert.Equal(t, visitDto.ClientID, result.ClientID)
	assert.Equal(t, visitDto.PetID, result.PetID)
	assert.Equal(t, visitDto.DoctorID, result.DoctorID)
	assert.Equal(t, visitDto.LastUpdatedByID, result.LastUpdatedByID)
}

func TestUpdateVisit_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	result, err := s.Update(createVisitForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteVisit_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteVisit_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewVisitService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createVisitForCreate() *dto.VisitDto {
	return &dto.VisitDto{
		DateTime:        time.Date(2024, time.January, 2, 0, 0, 0, 0, time.Local),
		Info:            "Информация",
		ClientID:        1,
		PetID:           1,
		DoctorID:        1,
		ServiceID:       4,
		LastUpdatedByID: 1,
	}
}
