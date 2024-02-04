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

type VisitDtoForBindError struct {
	DateTime        time.Time
	Info            string
	ClientID        string
	PetID           string
	DoctorID        string
	ServiceID       string
	LastUpdatedByID string
}

func TestGetVisitByID_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpVisitTestData(cont)

	visit := NewVisitController(cont)
	e.GET(config.APIv1VisitsID, func(c echo.Context) error { return visit.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1VisitsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Visit{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetVisit_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.GET(config.APIv1VisitsID, func(c echo.Context) error { return visit.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1VisitsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetVisit_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.GET(config.APIv1VisitsID, func(c echo.Context) error { return visit.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1VisitsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetVisitList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpVisitTestData(cont)

	visit := NewVisitController(cont)
	e.GET(config.APIv1Visits, func(c echo.Context) error { return visit.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Visits, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Visit{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetVisitList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.GET(config.APIv1Visits, func(c echo.Context) error { return visit.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Visits, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateVisit_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.POST(config.APIv1Visits, func(c echo.Context) error { return visit.Create(c) })

	param := createVisitForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Visits, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	userForLogin.BaseModel = &models.BaseModel{ID: 1}
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Visit{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateVisit_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.POST(config.APIv1Visits, func(c echo.Context) error { return visit.Create(c) })

	param := createVisitForBindError()
	req := test.NewJSONRequest("POST", config.APIv1Visits, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultVisitForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestCreateVisit_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.POST(config.APIv1Visits, func(c echo.Context) error { return visit.Create(c) })

	param := createVisitForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Visits, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "VisitDto.Info")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestCreateVisit_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.POST(config.APIv1Visits, func(c echo.Context) error { return visit.Create(c) })

	param := createVisitForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Visits, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateVisit_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.POST(config.APIv1Visits, func(c echo.Context) error { return visit.Create(c) })

	param := createVisitForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Visits, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateVisit_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpVisitTestData(cont)

	visit := NewVisitController(cont)
	e.PUT(config.APIv1VisitsID, func(c echo.Context) error { return visit.Update(c) })

	param := createVisitForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1VisitsID, "2"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	userForLogin.BaseModel = &models.BaseModel{ID: 1}
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Visit{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateVisit_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.PUT(config.APIv1VisitsID, func(c echo.Context) error { return visit.Update(c) })

	param := createVisitForBindError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1VisitsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultVisitForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestUpdateVisit_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.PUT(config.APIv1VisitsID, func(c echo.Context) error { return visit.Update(c) })

	param := createVisitForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1VisitsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "VisitDto.Info")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateVisit_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.PUT(config.APIv1VisitsID, func(c echo.Context) error { return visit.Update(c) })

	param := createVisitForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1VisitsID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateVisit_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.PUT(config.APIv1VisitsID, func(c echo.Context) error { return visit.Update(c) })

	param := createVisitForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1VisitsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteVisit_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpVisitTestData(cont)

	visit := NewVisitController(cont)
	e.DELETE(config.APIv1VisitsID, func(c echo.Context) error { return visit.Delete(c) })

	m := &models.Visit{}
	data, _ := m.Get(cont.Repository(), 2)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1VisitsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteVisit_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.DELETE(config.APIv1VisitsID, func(c echo.Context) error { return visit.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1VisitsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteVisit_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.DELETE(config.APIv1VisitsID, func(c echo.Context) error { return visit.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1VisitsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteVisit_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	visit := NewVisitController(cont)
	e.DELETE(config.APIv1VisitsID, func(c echo.Context) error { return visit.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1VisitsID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpVisitTestData(container container.Container) {
	rep := container.Repository()
	visit := createVisitForCreate().ToModel()
	visit.LastUpdatedByID = 1
	_, _ = visit.Create(rep)
}

func createVisitForCreate() *dto.VisitDto {
	return &dto.VisitDto{
		DateTime:  time.Date(2024, time.January, 2, 0, 0, 0, 0, time.Local),
		Info:      "Информация",
		ClientID:  1,
		PetID:     1,
		DoctorID:  1,
		ServiceID: 4,
	}
}

func createVisitForBindError() *VisitDtoForBindError {
	return &VisitDtoForBindError{
		DateTime:  time.Date(2024, time.January, 2, 0, 0, 0, 0, time.Local),
		Info:      "Информация",
		ClientID:  "Client",
		PetID:     "Pet",
		DoctorID:  "Doctor",
		ServiceID: "Service",
	}
}

func createResultVisitForBindError() *dto.VisitDto {
	return &dto.VisitDto{
		DateTime:  time.Date(2024, time.January, 2, 0, 0, 0, 0, time.Local),
		Info:      "Информация",
		ClientID:  0,
		PetID:     0,
		DoctorID:  0,
		ServiceID: 0,
	}
}

func createVisitForValidationError() *dto.VisitDto {
	return &dto.VisitDto{
		DateTime:  time.Date(2024, time.January, 2, 0, 0, 0, 0, time.Local),
		Info:      "Информация\n",
		ClientID:  1,
		PetID:     1,
		DoctorID:  1,
		ServiceID: 4,
	}
}

func createVisitForUpdate() *dto.VisitDto {
	return &dto.VisitDto{
		DateTime:  time.Date(2024, time.January, 2, 0, 0, 0, 0, time.Local),
		Info:      "ИнформацияUPD",
		ClientID:  1,
		PetID:     1,
		DoctorID:  1,
		ServiceID: 1,
	}
}
