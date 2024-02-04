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

type LeadDtoForBindError struct {
	Name     string
	Phone    string
	Email    string
	Comment  string
	Type     string
	Status   string
	DoctorID string
}

func TestGetLeadByID_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpLeadTestData(cont)

	lead := NewLeadController(cont)
	e.GET(config.APIv1LeadsID, func(c echo.Context) error { return lead.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1LeadsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Lead{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetLead_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.GET(config.APIv1LeadsID, func(c echo.Context) error { return lead.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1LeadsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetLead_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.GET(config.APIv1LeadsID, func(c echo.Context) error { return lead.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1LeadsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetLeadList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpLeadTestData(cont)

	lead := NewLeadController(cont)
	e.GET(config.APIv1Leads, func(c echo.Context) error { return lead.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Leads, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Lead{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetLeadList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.GET(config.APIv1Leads, func(c echo.Context) error { return lead.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Leads, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateLead_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.POST(config.APIv1Leads, func(c echo.Context) error { return lead.Create(c) })

	param := createLeadForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Leads, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	m := &models.Lead{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateLead_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.POST(config.APIv1Leads, func(c echo.Context) error { return lead.Create(c) })

	param := createLeadForBindError()
	req := test.NewJSONRequest("POST", config.APIv1Leads, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	result := createResultLeadForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestCreateLead_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.POST(config.APIv1Leads, func(c echo.Context) error { return lead.Create(c) })

	param := createLeadForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Leads, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "LeadDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Phone")
	assert.Contains(t, rec.Body.String(), "'e164'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Email")
	assert.Contains(t, rec.Body.String(), "'email'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Comment")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Type")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Status")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateLead_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpLeadTestData(cont)

	lead := NewLeadController(cont)
	e.PUT(config.APIv1LeadsID, func(c echo.Context) error { return lead.Update(c) })

	param := createLeadForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1LeadsID, "2"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	userForLogin.BaseModel = &models.BaseModel{ID: 1}
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Lead{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateLead_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.PUT(config.APIv1LeadsID, func(c echo.Context) error { return lead.Update(c) })

	param := createLeadForBindError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1LeadsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultLeadForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestUpdateLead_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.PUT(config.APIv1LeadsID, func(c echo.Context) error { return lead.Update(c) })

	param := createLeadForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1LeadsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "LeadDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Phone")
	assert.Contains(t, rec.Body.String(), "'e164'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Email")
	assert.Contains(t, rec.Body.String(), "'email'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Comment")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Type")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "LeadDto.Status")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateLead_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.PUT(config.APIv1LeadsID, func(c echo.Context) error { return lead.Update(c) })

	param := createLeadForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1LeadsID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteLead_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpLeadTestData(cont)

	lead := NewLeadController(cont)
	e.DELETE(config.APIv1LeadsID, func(c echo.Context) error { return lead.Delete(c) })

	m := &models.Lead{}
	data, _ := m.Get(cont.Repository(), 2)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1LeadsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteLead_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.DELETE(config.APIv1LeadsID, func(c echo.Context) error { return lead.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1LeadsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteLead_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.DELETE(config.APIv1LeadsID, func(c echo.Context) error { return lead.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1LeadsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteLead_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	lead := NewLeadController(cont)
	e.DELETE(config.APIv1LeadsID, func(c echo.Context) error { return lead.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1LeadsID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpLeadTestData(container container.Container) {
	rep := container.Repository()
	lead := createLeadForCreate().ToModel()
	_, _ = lead.Create(rep)
}

func createLeadForCreate() *dto.LeadDto {
	return &dto.LeadDto{
		Name:     "Клиент",
		Phone:    "+79998887766",
		Email:    "client@test.com",
		Comment:  "Комментарий",
		Type:     "consult-online",
		Status:   "in-progress",
		DoctorID: 1,
	}
}

func createLeadForBindError() *LeadDtoForBindError {
	return &LeadDtoForBindError{
		Name:     "Клиент",
		Phone:    "+79998887766",
		Email:    "client@test.com",
		Comment:  "Комментарий",
		Type:     "consult-online",
		Status:   "in-progress",
		DoctorID: "Doctor",
	}
}

func createResultLeadForBindError() *dto.LeadDto {
	return &dto.LeadDto{
		Name:     "Клиент",
		Phone:    "+79998887766",
		Email:    "client@test.com",
		Comment:  "Комментарий",
		Type:     "consult-online",
		Status:   "in-progress",
		DoctorID: 0,
	}
}

func createLeadForValidationError() *dto.LeadDto {
	return &dto.LeadDto{
		Name:            "Клиент2",
		Phone:           "2",
		Email:           "client2test.com",
		Comment:         "Комментарий\n",
		Type:            "consult-online\n",
		Status:          "in-progress\n",
		DoctorID:        1,
		LastUpdatedByID: 1,
	}
}

func createLeadForUpdate() *dto.LeadDto {
	return &dto.LeadDto{
		Name:            "КлиентUPD",
		Phone:           "+79998880000",
		Email:           "client@test.UPD",
		Comment:         "КомментарийUPD",
		Type:            "consult-in-clinic",
		Status:          "rejected",
		DoctorID:        1,
		LastUpdatedByID: 1,
	}
}
