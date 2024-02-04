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

func TestGetRole_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1RolesID, func(c echo.Context) error { return role.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1RolesID, "4"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Role{}
	data, _ := m.Get(cont.Repository(), 4)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetRole_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1RolesID, func(c echo.Context) error { return role.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1RolesID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetRole_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1RolesID, func(c echo.Context) error { return role.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1RolesID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetRole_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1RolesID, func(c echo.Context) error { return role.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1RolesID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestGetRoleList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1Roles, func(c echo.Context) error { return role.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Roles, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Role{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetRoleList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1Roles, func(c echo.Context) error { return role.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Roles, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetRoleList_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.GET(config.APIv1Roles, func(c echo.Context) error { return role.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Roles, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestCreateRole_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.POST(config.APIv1Roles, func(c echo.Context) error { return role.Create(c) })

	param := createRoleForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Roles, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Role{}
	data, _ := m.Get(cont.Repository(), 5)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreateRole_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.POST(config.APIv1Roles, func(c echo.Context) error { return role.Create(c) })

	param := createRoleForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Roles, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "RoleDto.Name")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
}

func TestCreateRole_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.POST(config.APIv1Roles, func(c echo.Context) error { return role.Create(c) })

	param := createRoleForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Roles, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateRole_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.POST(config.APIv1Roles, func(c echo.Context) error { return role.Create(c) })

	param := createRoleForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Roles, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateRole_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.PUT(config.APIv1RolesID, func(c echo.Context) error { return role.Update(c) })

	param := createRoleForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1RolesID, "4"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Role{}
	data, _ := m.Get(cont.Repository(), 4)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdateRole_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.PUT(config.APIv1RolesID, func(c echo.Context) error { return role.Update(c) })

	param := createRoleForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1RolesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "RoleDto.Name")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
}

func TestUpdateRole_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.PUT(config.APIv1RolesID, func(c echo.Context) error { return role.Update(c) })

	param := createRoleForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1RolesID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateRole_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.PUT(config.APIv1RolesID, func(c echo.Context) error { return role.Update(c) })

	param := createRoleForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1RolesID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestDeleteRole_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.DELETE(config.APIv1RolesID, func(c echo.Context) error { return role.Delete(c) })

	m := &models.Role{}
	data, _ := m.Get(cont.Repository(), 4)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1RolesID, "4"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteRole_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.DELETE(config.APIv1RolesID, func(c echo.Context) error { return role.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1RolesID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteRole_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.DELETE(config.APIv1RolesID, func(c echo.Context) error { return role.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1RolesID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteRole_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	role := NewRoleController(cont)
	e.DELETE(config.APIv1RolesID, func(c echo.Context) error { return role.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1RolesID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func createRoleForCreate() *dto.RoleDto {
	return &dto.RoleDto{
		Name: "Test",
	}
}

func createRoleForValidationError() *dto.RoleDto {
	return &dto.RoleDto{
		Name: "Test3@",
	}
}

func createRoleForUpdate() *dto.RoleDto {
	return &dto.RoleDto{
		Name: "TestU",
	}
}
