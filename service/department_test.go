package service

import (
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"testing"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindDepartmentByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindDepartmentBySlug_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	resultByID, _ := s.Get("2")
	resultBySlug, err := s.Get(resultByID.Slug)

	assert.Equal(t, uint(2), resultBySlug.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resultBySlug)
}

func TestFindDepartmentByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindDepartmentByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllDepartments_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 2)
	assert.NoError(t, err)
}

func TestCreateDepartment_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	departmentDto := createDepartmentForCreate()
	_, err := s.Create(departmentDto)

	result, _ := s.Get("3")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, departmentDto.Name, result.Name)
	assert.Equal(t, makeSlugForDepartment(departmentDto.ToModel()), result.Slug)
	assert.Equal(t, departmentDto.Services[0], result.Services[0].ID)
}

func TestUpdateDepartment_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	departmentDto := createDepartmentForCreate()
	_, err := s.Update(departmentDto, "1")

	result, _ := s.Get("1")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, departmentDto.Name, result.Name)
	assert.Equal(t, makeSlugForDepartment(departmentDto.ToModel()), result.Slug)
	assert.Equal(t, departmentDto.Services[0], result.Services[0].ID)
}

func TestUpdateDepartment_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	result, err := s.Update(createDepartmentForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteDepartment_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteDepartment_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewDepartmentService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createDepartmentForCreate() *dto.DepartmentDto {
	return &dto.DepartmentDto{
		Name:     "Хирургия",
		Services: []uint{1},
	}
}

func makeSlugForDepartment(department *models.Department) string {
	slug.MaxLength = 40
	slug.EnableSmartTruncate = false
	slug.CustomSub = map[string]string{
		"ь": "",
		"Ь": "",
		"ъ": "",
		"Ъ": "",
	}
	return slug.Make(department.Name)
}
