package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
	"vet-clinic/util"
)

type DepartmentDtoForBindError struct {
	Name     string
	Services string
}

func TestGetDepartmentByID_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpDepartmentTestData(cont)

	department := NewDepartmentController(cont)
	e.GET(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1DepartmentsID, "3"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Department{}
	data, _ := m.Get(cont.Repository(), 3)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetDepartmentBySlug_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.GET(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Get(c) })

	rep := cont.Repository()
	testDepartment, _ := createDepartmentForCreate().ToModel().Create(rep)

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1DepartmentsID, testDepartment.Slug), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Department{}
	data, _ := m.Get(cont.Repository(), 3)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetDepartment_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.GET(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1DepartmentsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetDepartment_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.GET(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1DepartmentsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetDepartmentList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpDepartmentTestData(cont)

	department := NewDepartmentController(cont)
	e.GET(config.APIv1Departments, func(c echo.Context) error { return department.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Departments, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Department{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetDepartmentList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.GET(config.APIv1Departments, func(c echo.Context) error { return department.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Departments, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateDepartment_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.POST(config.APIv1Departments, func(c echo.Context) error { return department.Create(c) })

	param := createDepartmentForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Departments, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Department{}
	data, _ := m.Get(cont.Repository(), 3)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateDepartment_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.POST(config.APIv1Departments, func(c echo.Context) error { return department.Create(c) })

	param := createDepartmentForBindError()
	req := test.NewJSONRequest("POST", config.APIv1Departments, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultDepartmentForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestCreateDepartment_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.POST(config.APIv1Departments, func(c echo.Context) error { return department.Create(c) })

	param := createDepartmentForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Departments, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "DepartmentDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
}

func TestCreateDepartment_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.POST(config.APIv1Departments, func(c echo.Context) error { return department.Create(c) })

	param := createDepartmentForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Departments, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateDepartment_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.POST(config.APIv1Departments, func(c echo.Context) error { return department.Create(c) })

	param := createDepartmentForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Departments, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateDepartment_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpDepartmentTestData(cont)

	department := NewDepartmentController(cont)
	e.PUT(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Update(c) })

	param := createDepartmentForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1DepartmentsID, "3"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Department{}
	data, _ := m.Get(cont.Repository(), 3)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateDepartment_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.PUT(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Update(c) })

	param := createDepartmentForBindError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1DepartmentsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultDepartmentForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestUpdateDepartment_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.PUT(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Update(c) })

	param := createDepartmentForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1DepartmentsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "DepartmentDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
}

func TestUpdateDepartment_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.PUT(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Update(c) })

	param := createDepartmentForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1DepartmentsID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateDepartment_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.PUT(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Update(c) })

	param := createDepartmentForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1DepartmentsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteDepartment_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpDepartmentTestData(cont)

	department := NewDepartmentController(cont)
	e.DELETE(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Delete(c) })

	m := &models.Department{}
	data, _ := m.Get(cont.Repository(), 3)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1DepartmentsID, "3"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteDepartment_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.DELETE(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1DepartmentsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteDepartment_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.DELETE(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1DepartmentsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteDepartment_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	department := NewDepartmentController(cont)
	e.DELETE(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1DepartmentsID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpDepartmentTestData(container container.Container) {
	rep := container.Repository()
	department := createDepartmentForCreate().ToModel()
	_, _ = department.Create(rep)
}

func createDepartmentForCreate() *dto.DepartmentDto {
	return &dto.DepartmentDto{
		Name:     "Хирургия",
		Services: []uint{1},
	}
}

func createDepartmentForBindError() *DepartmentDtoForBindError {
	return &DepartmentDtoForBindError{
		Name:     "Хирургия",
		Services: "Услуга",
	}
}

func createResultDepartmentForBindError() *dto.DepartmentDto {
	return &dto.DepartmentDto{
		Name:     "Хирургия",
		Services: nil,
	}
}

func createDepartmentForValidationError() *dto.DepartmentDto {
	return &dto.DepartmentDto{
		Name:     "Хирургия2",
		Services: []uint{1},
	}
}

func createDepartmentForUpdate() *dto.DepartmentDto {
	return &dto.DepartmentDto{
		Name:     "Лаборатория",
		Services: []uint{2},
	}
}
