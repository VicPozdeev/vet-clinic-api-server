package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"vet-clinic/config"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
	"vet-clinic/util"
)

func TestGetCategory_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.GET(config.APIv1CategoriesID, func(c echo.Context) error { return category.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1CategoriesID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Category{}
	data, _ := m.Get(cont.Repository(), 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetCategory_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.GET(config.APIv1CategoriesID, func(c echo.Context) error { return category.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1CategoriesID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetCategory_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.GET(config.APIv1CategoriesID, func(c echo.Context) error { return category.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1CategoriesID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetCategoryList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.GET(config.APIv1Categories, func(c echo.Context) error { return category.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Categories, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Category{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetCategoryList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.GET(config.APIv1Categories, func(c echo.Context) error { return category.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Categories, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateCategory_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.POST(config.APIv1Categories, func(c echo.Context) error { return category.Create(c) })

	param := createCategoryForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Categories, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Category{}
	data, _ := m.Get(cont.Repository(), 5)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateCategory_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.POST(config.APIv1Categories, func(c echo.Context) error { return category.Create(c) })

	param := createCategoryForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Categories, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "CategoryDto.Name")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestCreateCategory_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.POST(config.APIv1Categories, func(c echo.Context) error { return category.Create(c) })

	param := createCategoryForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Categories, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateCategory_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.POST(config.APIv1Categories, func(c echo.Context) error { return category.Create(c) })

	param := createCategoryForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Categories, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateCategory_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.PUT(config.APIv1CategoriesID, func(c echo.Context) error { return category.Update(c) })

	param := createCategoryForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1CategoriesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Category{}
	data, _ := m.Get(cont.Repository(), 1)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateCategory_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.PUT(config.APIv1CategoriesID, func(c echo.Context) error { return category.Update(c) })

	param := createCategoryForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1CategoriesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "CategoryDto.Name")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateCategory_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.PUT(config.APIv1CategoriesID, func(c echo.Context) error { return category.Update(c) })

	param := createCategoryForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1CategoriesID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateCategory_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.PUT(config.APIv1CategoriesID, func(c echo.Context) error { return category.Update(c) })

	param := createCategoryForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1CategoriesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteCategory_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.DELETE(config.APIv1CategoriesID, func(c echo.Context) error { return category.Delete(c) })

	m := &models.Category{}
	data, _ := m.Get(cont.Repository(), 1)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1CategoriesID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteCategory_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.DELETE(config.APIv1CategoriesID, func(c echo.Context) error { return category.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1CategoriesID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteCategory_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.DELETE(config.APIv1CategoriesID, func(c echo.Context) error { return category.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1CategoriesID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteCategory_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	category := NewCategoryController(cont)
	e.DELETE(config.APIv1CategoriesID, func(c echo.Context) error { return category.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1CategoriesID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func createCategoryForCreate() *dto.CategoryDto {
	return &dto.CategoryDto{
		Name: "Test",
	}
}

func createCategoryForValidationError() *dto.CategoryDto {
	return &dto.CategoryDto{
		Name: "Test3@\n",
	}
}

func createCategoryForUpdate() *dto.CategoryDto {
	return &dto.CategoryDto{
		Name: "TestU",
	}
}
