package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindServiceByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindServiceByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindServiceByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllServices_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 8)
	assert.NoError(t, err)
}

func TestCreateService_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	serviceDto := createServiceForCreate()
	_, err := s.Create(serviceDto)

	result, _ := s.Get("9")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, serviceDto.Name, result.Name)
	assert.Equal(t, serviceDto.Price, result.Price)
	assert.Equal(t, serviceDto.CategoryID, result.CategoryID)
}

func TestUpdateService_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	serviceDto := createServiceForCreate()
	_, err := s.Update(serviceDto, "1")

	result, _ := s.Get("1")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, serviceDto.Name, result.Name)
	assert.Equal(t, serviceDto.Price, result.Price)
	assert.Equal(t, serviceDto.CategoryID, result.CategoryID)
}

func TestUpdateService_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	result, err := s.Update(createServiceForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteService_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteService_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewServiceService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createServiceForCreate() *dto.ServiceDto {
	return &dto.ServiceDto{
		Name:       "Общий анализ крови",
		Price:      2200,
		CategoryID: 2,
	}
}
