package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindRoleByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindRoleByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindRoleByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllRoles_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 4)
	assert.NoError(t, err)
}

func TestCreateRole_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	_, err := s.Create(createRoleForCreate())
	result := createResultRole()
	result.ID = 5

	data, _ := s.Get("5")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, data, result)
}

func TestUpdateRole_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	_, err := s.Update(createRoleForCreate(), "1")
	result := createResultRole()
	result.ID = 1

	data, _ := s.Get("1")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, data, result)
}

func TestUpdateRole_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	result, err := s.Update(createRoleForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteRole_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteRole_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewRoleService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createRoleForCreate() *dto.RoleDto {
	return &dto.RoleDto{
		Name: "Role",
	}
}

func createResultRole() *models.Role {
	return &models.Role{
		Name: "Role",
	}
}
