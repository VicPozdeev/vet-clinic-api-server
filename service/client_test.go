package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindClientByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindClientByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindClientByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllClients_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 1)
	assert.NoError(t, err)
}

func TestCreateClient_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	_, err := s.Create(createClientForCreate())
	result := createResultClient()

	data, _ := s.Get("2")
	data.BaseModel = nil
	data.BirthDate = data.BirthDate.Local()

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, data, result)
}

func TestUpdateClient_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	_, err := s.Update(createClientForCreate(), "1")
	result := createResultClient()

	data, _ := s.Get("1")
	data.BaseModel = nil
	data.BirthDate = data.BirthDate.Local()

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, data, result)
}

func TestUpdateClient_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	result, err := s.Update(createClientForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteClient_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteClient_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewClientService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createClientForCreate() *dto.ClientDto {
	return &dto.ClientDto{
		Surname:    "Фамилия",
		Name:       "Имя",
		Patronymic: "Отчество",
		Sex:        "Мужской",
		BirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Phone:      "+71234567899",
		Email:      "client@test.com",
		Info:       "Информация",
	}
}

func createResultClient() *models.Client {
	return &models.Client{
		Surname:    "Фамилия",
		Name:       "Имя",
		Patronymic: "Отчество",
		Sex:        "Мужской",
		BirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Phone:      "+71234567899",
		Email:      "client@test.com",
		Info:       "Информация",
	}
}
