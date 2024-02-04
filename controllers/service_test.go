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

type ServiceDtoForBindError struct {
	Name       string
	Price      string
	CategoryID string
}

func TestGetServiceByID_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpServiceTestData(cont)

	service := NewServiceController(cont)
	e.GET(config.APIv1ServicesID, func(c echo.Context) error { return service.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1ServicesID, "9"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Service{}
	data, _ := m.Get(cont.Repository(), 9)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetService_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.GET(config.APIv1ServicesID, func(c echo.Context) error { return service.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1ServicesID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetService_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.GET(config.APIv1ServicesID, func(c echo.Context) error { return service.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1ServicesID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetServiceList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpServiceTestData(cont)

	service := NewServiceController(cont)
	e.GET(config.APIv1Services, func(c echo.Context) error { return service.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Services, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Service{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetServiceList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.GET(config.APIv1Services, func(c echo.Context) error { return service.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Services, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateService_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.POST(config.APIv1Services, func(c echo.Context) error { return service.Create(c) })

	param := createServiceForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Services, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Service{}
	data, _ := m.Get(cont.Repository(), 9)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateService_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.POST(config.APIv1Services, func(c echo.Context) error { return service.Create(c) })

	param := createServiceForBindError()
	req := test.NewJSONRequest("POST", config.APIv1Services, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultServiceForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestCreateService_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.POST(config.APIv1Services, func(c echo.Context) error { return service.Create(c) })

	param := createServiceForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Services, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ServiceDto.Name")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestCreateService_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.POST(config.APIv1Services, func(c echo.Context) error { return service.Create(c) })

	param := createServiceForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Services, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateService_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.POST(config.APIv1Services, func(c echo.Context) error { return service.Create(c) })

	param := createServiceForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Services, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateService_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpServiceTestData(cont)

	service := NewServiceController(cont)
	e.PUT(config.APIv1ServicesID, func(c echo.Context) error { return service.Update(c) })

	param := createServiceForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ServicesID, "9"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Service{}
	data, _ := m.Get(cont.Repository(), 9)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateService_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.PUT(config.APIv1ServicesID, func(c echo.Context) error { return service.Update(c) })

	param := createServiceForBindError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ServicesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultServiceForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestUpdateService_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.PUT(config.APIv1ServicesID, func(c echo.Context) error { return service.Update(c) })

	param := createServiceForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ServicesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ServiceDto.Name")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateService_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.PUT(config.APIv1ServicesID, func(c echo.Context) error { return service.Update(c) })

	param := createServiceForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ServicesID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateService_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.PUT(config.APIv1ServicesID, func(c echo.Context) error { return service.Update(c) })

	param := createServiceForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ServicesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteService_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpServiceTestData(cont)

	service := NewServiceController(cont)
	e.DELETE(config.APIv1ServicesID, func(c echo.Context) error { return service.Delete(c) })

	m := &models.Service{}
	data, _ := m.Get(cont.Repository(), 9)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ServicesID, "9"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteService_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.DELETE(config.APIv1ServicesID, func(c echo.Context) error { return service.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ServicesID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteService_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.DELETE(config.APIv1ServicesID, func(c echo.Context) error { return service.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ServicesID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteService_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	service := NewServiceController(cont)
	e.DELETE(config.APIv1ServicesID, func(c echo.Context) error { return service.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ServicesID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpServiceTestData(container container.Container) {
	rep := container.Repository()
	service := createServiceForCreate().ToModel()
	_, _ = service.Create(rep)
}

func createServiceForCreate() *dto.ServiceDto {
	return &dto.ServiceDto{
		Name:       "Общий анализ крови",
		Price:      2200,
		CategoryID: 2,
	}
}

func createServiceForBindError() *ServiceDtoForBindError {
	return &ServiceDtoForBindError{
		Name:       "Общий анализ крови",
		Price:      "Price",
		CategoryID: "Category",
	}
}

func createResultServiceForBindError() *dto.ServiceDto {
	return &dto.ServiceDto{
		Name:       "Общий анализ крови",
		Price:      0,
		CategoryID: 0,
	}
}

func createServiceForValidationError() *dto.ServiceDto {
	return &dto.ServiceDto{
		Name:       "Общий анализ крови\n",
		Price:      2200,
		CategoryID: 2,
	}
}

func createServiceForUpdate() *dto.ServiceDto {
	return &dto.ServiceDto{
		Name:       "Электрическая кардиоверсия",
		Price:      14200,
		CategoryID: 3,
	}
}
