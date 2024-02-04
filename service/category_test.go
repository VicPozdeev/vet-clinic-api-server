package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindCategoryByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	result, err := s.Get("1")

	assert.Equal(t, uint(1), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindCategoryByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindCategoryByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllCategories_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 4)
	assert.NoError(t, err)
}

func TestCreateCategory_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	_, err := s.Create(createCategoryForCreate())
	result := createResultCategory()
	result.ID = 5

	data, _ := s.Get("5")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, data, result)
}

func TestUpdateCategory_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	_, err := s.Update(createCategoryForCreate(), "1")
	result := createResultCategory()
	result.ID = 1

	data, _ := s.Get("1")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, result, data)
}

func TestUpdateCategory_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	result, err := s.Update(createCategoryForCreate(), "99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteCategory_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	data, _ := s.Get("1")

	result, err := s.Delete("1")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteCategory_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewCategoryService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func createCategoryForCreate() *dto.CategoryDto {
	return &dto.CategoryDto{
		Name: "Category",
	}
}

func createResultCategory() *models.Category {
	return &models.Category{
		Name: "Category",
	}
}
