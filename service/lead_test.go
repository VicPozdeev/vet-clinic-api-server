package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindLeadByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	result, err := s.Get("1")

	assert.NotEmpty(t, result)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
}

func TestFindLeadByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindLeadByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllLeads_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 1)
	assert.NoError(t, err)
}

func TestCreateLead_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	leadDto := createLeadForCreate()
	_, err := s.Create(leadDto)

	result, _ := s.Get("2")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, leadDto.Name, result.Name)
	assert.Equal(t, leadDto.Phone, result.Phone)
	assert.Equal(t, leadDto.Email, result.Email)
	assert.Equal(t, leadDto.Comment, result.Comment)
	assert.Equal(t, leadDto.Type, result.Type)
	assert.Equal(t, "open", result.Status)
	assert.Equal(t, leadDto.DoctorID, result.DoctorID)
	assert.Equal(t, uint(0), result.LastUpdatedByID)
}

func TestUpdateLead_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	leadDto := createLeadForCreate()
	_, err := s.Update(leadDto, "1")

	result, _ := s.Get("1")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, leadDto.Name, result.Name)
	assert.Equal(t, leadDto.Phone, result.Phone)
	assert.Equal(t, leadDto.Email, result.Email)
	assert.Equal(t, leadDto.Comment, result.Comment)
	assert.Equal(t, leadDto.Type, result.Type)
	assert.Equal(t, leadDto.Status, result.Status)
	assert.Equal(t, leadDto.DoctorID, result.DoctorID)
	assert.Equal(t, leadDto.LastUpdatedByID, result.LastUpdatedByID)
}

func TestUpdateLead_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	result, err := s.Update(createLeadForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteLead_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteLead_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewLeadService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createLeadForCreate() *dto.LeadDto {
	return &dto.LeadDto{
		Name:            "Клиент",
		Phone:           "+79998887766",
		Email:           "client@test.com",
		Comment:         "Комментарий",
		Type:            "consult-online",
		Status:          "in-progress",
		DoctorID:        1,
		LastUpdatedByID: 1,
	}
}
