package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
	"vet-clinic/util"
)

func TestGetClient_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpClientTestData(cont)

	client := NewClientController(cont)
	e.GET(config.APIv1ClientsID, func(c echo.Context) error { return client.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1ClientsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Client{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetClient_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.GET(config.APIv1ClientsID, func(c echo.Context) error { return client.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1ClientsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetClient_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.GET(config.APIv1ClientsID, func(c echo.Context) error { return client.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1ClientsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetClientList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpClientTestData(cont)

	client := NewClientController(cont)
	e.GET(config.APIv1Clients, func(c echo.Context) error { return client.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Clients, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Client{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetClientList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.GET(config.APIv1Clients, func(c echo.Context) error { return client.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Clients, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateClient_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.POST(config.APIv1Clients, func(c echo.Context) error { return client.Create(c) })

	param := createClientForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Clients, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Client{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateClient_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.POST(config.APIv1Clients, func(c echo.Context) error { return client.Create(c) })

	param := createClientForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Clients, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ClientDto.Surname")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Patronymic")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Sex")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Phone")
	assert.Contains(t, rec.Body.String(), "'e164'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Email")
	assert.Contains(t, rec.Body.String(), "'email'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Info")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestCreateClient_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.POST(config.APIv1Clients, func(c echo.Context) error { return client.Create(c) })

	param := createClientForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Clients, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateClient_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.POST(config.APIv1Clients, func(c echo.Context) error { return client.Create(c) })

	param := createClientForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Clients, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateClient_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpClientTestData(cont)

	client := NewClientController(cont)
	e.PUT(config.APIv1ClientsID, func(c echo.Context) error { return client.Update(c) })

	param := createClientForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ClientsID, "2"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Client{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateClient_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.PUT(config.APIv1ClientsID, func(c echo.Context) error { return client.Update(c) })

	param := createClientForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ClientsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "ClientDto.Surname")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Patronymic")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Sex")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Phone")
	assert.Contains(t, rec.Body.String(), "'e164'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Email")
	assert.Contains(t, rec.Body.String(), "'email'")
	assert.Contains(t, rec.Body.String(), "ClientDto.Info")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateClient_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.PUT(config.APIv1ClientsID, func(c echo.Context) error { return client.Update(c) })

	param := createClientForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ClientsID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateClient_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.PUT(config.APIv1ClientsID, func(c echo.Context) error { return client.Update(c) })

	param := createClientForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1ClientsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteClient_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpClientTestData(cont)

	client := NewClientController(cont)
	e.DELETE(config.APIv1ClientsID, func(c echo.Context) error { return client.Delete(c) })

	m := &models.Client{}
	data, _ := m.Get(cont.Repository(), 2)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ClientsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteClient_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.DELETE(config.APIv1ClientsID, func(c echo.Context) error { return client.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ClientsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteClient_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.DELETE(config.APIv1ClientsID, func(c echo.Context) error { return client.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ClientsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteClient_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	client := NewClientController(cont)
	e.DELETE(config.APIv1ClientsID, func(c echo.Context) error { return client.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1ClientsID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpClientTestData(container container.Container) {
	rep := container.Repository()
	client := createClientForCreate().ToModel()
	_, _ = client.Create(rep)
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

func createClientForValidationError() *dto.ClientDto {
	return &dto.ClientDto{
		Surname:    "Фамилия2",
		Name:       "Имя2",
		Patronymic: "Отчество2",
		Sex:        "Мужской2",
		BirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Phone:      "2",
		Email:      "client2test.com",
		Info:       "Информация\n",
	}
}

func createClientForUpdate() *dto.ClientDto {
	return &dto.ClientDto{
		Surname:    "ФамилияUPD",
		Name:       "ИмяUPD",
		Patronymic: "ОтчествоUPD",
		Sex:        "МужскойUPD",
		BirthDate:  time.Date(1990, time.January, 4, 0, 0, 0, 0, time.Local),
		Phone:      "+71234560000",
		Email:      "client@test.upd",
		Info:       "ИнформацияUPD",
	}
}
