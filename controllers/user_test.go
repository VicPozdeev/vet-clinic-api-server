package controllers

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
	"vet-clinic/util"
)

type UserDtoForCreateBindError struct {
	Username string
	Password string
	RoleID   string
}

type UserDtoForUpdateBindError struct {
	Email      string
	Phone      string
	Active     string
	Surname    string
	Name       string
	Patronymic string
	Sex        string
	BirthDate  time.Time
	Profession string
	Info       string
	RoleID     string
	Department string
	Services   string
}

func TestLoginByUsername_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })

	testUser := setAndGetFulfilledUser(cont, util.Staff)

	param := createLoginDtoForLogin(testUser.Username)
	req := test.NewJSONRequest("POST", config.APIv1Login, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestLoginByEmail_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })

	testUser := setAndGetFulfilledUser(cont, util.Staff)

	param := createLoginDtoForLogin(testUser.Email)
	req := test.NewJSONRequest("POST", config.APIv1Login, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestLoginByPhone_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })

	testUser := setAndGetFulfilledUser(cont, util.Staff)

	param := createLoginDtoForLogin(testUser.Phone)
	req := test.NewJSONRequest("POST", config.APIv1Login, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestLogin_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })

	setUpUserTestData(cont, util.Staff)

	param := createLoginDtoForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Login, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "LoginDto.Login")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'username|email|e164'")
}

func TestLogin_WrongPassword(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })

	setUpUserTestData(cont, util.Staff)

	param := createLoginDtoForWrongPassword()
	req := test.NewJSONRequest("POST", config.APIv1Login, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "crypto/bcrypt: hashedPassword is not the hash of the given password"}
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestLogin_EntityNotFound(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })

	setUpUserTestData(cont, util.Staff)

	param := createLoginDtoForLogin("NotExist")
	req := test.NewJSONRequest("POST", config.APIv1Login, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetUserByID_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1UsersID, func(c echo.Context) error { return user.Get(c) })

	setUpUserTestData(cont, util.Staff)

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1UsersID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetUserBySlug_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1UsersID, func(c echo.Context) error { return user.Get(c) })

	testUser := setAndGetFulfilledUser(cont, util.Staff)

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1UsersID, testUser.Slug), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetUser_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1UsersID, func(c echo.Context) error { return user.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1UsersID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetUser_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1UsersID, func(c echo.Context) error { return user.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1UsersID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetSelf_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1Profile, func(c echo.Context) error { return user.GetSelf(c) })

	setUpUserTestData(cont, util.Staff)

	req := httptest.NewRequest("GET", config.APIv1Profile, nil)
	rec := httptest.NewRecorder()

	m := &models.User{}
	userForLogin, _ := m.Get(cont.Repository(), 2)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetSelf_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1Profile, func(c echo.Context) error { return user.GetSelf(c) })

	req := httptest.NewRequest("GET", config.APIv1Profile, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetUserList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1Users, func(c echo.Context) error { return user.GetAll(c) })

	setUpUserTestData(cont, util.Staff)

	req := httptest.NewRequest("GET", config.APIv1Users, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetUserList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.GET(config.APIv1Users, func(c echo.Context) error { return user.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Users, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateUser_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
	assert.True(t, data.Active)
	assert.Equal(t, strconv.Itoa(int(data.ID)), data.Slug)
}

func TestCreateUser_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreateBindError()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultUserForCreateBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestCreateUser_ValidationNameError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreateValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "UserCreateDto.Username")
	assert.Contains(t, rec.Body.String(), "'username'")
	assert.Contains(t, rec.Body.String(), "UserCreateDto.Password")
	assert.Contains(t, rec.Body.String(), "'password'")
}

func TestCreateUser_ValidationMinError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreateValidationMinError()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "UserCreateDto.Password")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'min'")
}

func TestCreateUser_ValidationMaxError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreateValidationMaxError()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "UserCreateDto.Password")
	assert.Contains(t, rec.Body.String(), "Error:Field validation")
	assert.Contains(t, rec.Body.String(), "'max'")
}

func TestCreateUser_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreateUser_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })

	param := createUserForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Users, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateUser_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1UsersID, func(c echo.Context) error { return user.Update(c) })

	setUpUserTestData(cont, util.Staff)

	param := createUserForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1UsersID, "2"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
	assert.Equal(t, makeSlugForUser(param.ToModel(true)), data.Slug)
}

func TestUpdateUser_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1UsersID, func(c echo.Context) error { return user.Update(c) })

	param := createUserForUpdateBindError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1UsersID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultUserForUpdateBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestUpdateUser_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1UsersID, func(c echo.Context) error { return user.Update(c) })

	param := createUserForUpdateValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1UsersID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Email")
	assert.Contains(t, rec.Body.String(), "'email'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Phone")
	assert.Contains(t, rec.Body.String(), "'e164'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Surname")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Patronymic")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Sex")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Profession")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "UserUpdateDto.Info")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
}

func TestUpdateUser_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1UsersID, func(c echo.Context) error { return user.Update(c) })

	param := createUserForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1UsersID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdateUser_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1UsersID, func(c echo.Context) error { return user.Update(c) })

	param := createUserForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1UsersID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Administrator)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func TestUpdateSelf_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1Profile, func(c echo.Context) error { return user.UpdateSelf(c) })

	setUpUserTestData(cont, util.Staff)

	param := createUserForUpdate()
	req := test.NewJSONRequest("PUT", config.APIv1Profile, param)
	rec := httptest.NewRecorder()

	m := &models.User{}
	userForLogin, _ := m.Get(cont.Repository(), 2)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
	assert.Equal(t, makeSlugForUser(param.ToModel(true)), data.Slug)
}

func TestUpdateSelf_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1Profile, func(c echo.Context) error { return user.UpdateSelf(c) })

	param := createUserForUpdate()
	req := test.NewJSONRequest("PUT", config.APIv1Profile, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdatePassword_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1Password, func(c echo.Context) error { return user.UpdatePassword(c) })

	setUpUserTestData(cont, util.Staff)

	param := createPasswordForUpdate()
	req := test.NewJSONRequest("PUT", config.APIv1Password, param)
	rec := httptest.NewRecorder()

	m := &models.User{}
	userForLogin, _ := m.Get(cont.Repository(), 2)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdatePassword_NotEqualConfirmPassword(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1Password, func(c echo.Context) error { return user.UpdatePassword(c) })

	setUpUserTestData(cont, util.Staff)

	param := createPasswordForNotEqualConfirmPassword()
	req := test.NewJSONRequest("PUT", config.APIv1Password, param)
	rec := httptest.NewRecorder()

	m := &models.User{}
	userForLogin, _ := m.Get(cont.Repository(), 2)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "Key: 'UpdatePasswordDto.NewPassword' Error:Field validation for 'NewPassword' failed on the 'eqfield' tag"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestUpdatePassword_EqualOldPassword(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.PUT(config.APIv1Password, func(c echo.Context) error { return user.UpdatePassword(c) })

	setUpUserTestData(cont, util.Staff)

	param := createPasswordForEqualOldPassword()
	req := test.NewJSONRequest("PUT", config.APIv1Password, param)
	rec := httptest.NewRecorder()

	m := &models.User{}
	userForLogin, _ := m.Get(cont.Repository(), 2)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "Key: 'UpdatePasswordDto.NewPassword' Error:Field validation for 'NewPassword' failed on the 'nefield' tag"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteUser_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.DELETE(config.APIv1UsersID, func(c echo.Context) error { return user.Delete(c) })

	setUpUserTestData(cont, util.Staff)

	m := &models.User{}
	data, _ := m.Get(cont.Repository(), 1)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1UsersID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeleteUser_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.DELETE(config.APIv1UsersID, func(c echo.Context) error { return user.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1UsersID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeleteUser_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.DELETE(config.APIv1UsersID, func(c echo.Context) error { return user.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1UsersID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeleteUser_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	user := NewUserController(cont)
	e.DELETE(config.APIv1UsersID, func(c echo.Context) error { return user.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1UsersID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpUserTestData(container container.Container, level util.AccessLevel) {
	rep := container.Repository()
	role := &models.Role{}
	rep.First(role, "name = ?", level.ToString())
	user := models.NewUser("Test2", "Password2!", role.ID)
	_, _ = user.Create(rep)
}

func setAndGetFulfilledUser(container container.Container, level util.AccessLevel) *models.User {
	setUpUserTestData(container, level)
	rep := container.Repository()

	user := createUserForUpdate().ToModel(util.Owner.AccessAllowed(level))
	user, _ = user.Update(rep, 2, util.Owner.AccessAllowed(level))
	return user
}

func userWithAccessLevel(level util.AccessLevel) *models.User {
	return &models.User{Role: &models.Role{Name: level.ToString()}}
}

func createLoginDtoForLogin(username string) *dto.LoginDto {
	return &dto.LoginDto{
		Login:    username,
		Password: "Password2!",
	}
}

func createLoginDtoForValidationError() *dto.LoginDto {
	return &dto.LoginDto{
		Login:    "Test test",
		Password: "weak_password",
	}
}

func createLoginDtoForWrongPassword() *dto.LoginDto {
	return &dto.LoginDto{
		Login:    "Test2",
		Password: "Password2222!",
	}
}

func createUserForCreate() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "Test2",
		Password: "Password2!",
		RoleID:   1,
	}
}

func createUserForCreateBindError() *UserDtoForCreateBindError {
	return &UserDtoForCreateBindError{
		Username: "Test2",
		Password: "Password2!",
		RoleID:   "role",
	}
}

func createResultUserForCreateBindError() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "Test2",
		Password: "Password2!",
		RoleID:   0,
	}
}

func createUserForCreateValidationError() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "1startDigit",
		Password: "weak_password",
		RoleID:   1,
	}
}

func createUserForCreateValidationMinError() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "Test2",
		Password: "pass",
		RoleID:   1,
	}
}

func createUserForCreateValidationMaxError() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "Test2",
		Password: "Password2!Password2!Password2!Password2!Password2!Password2!Password2!Pas",
		RoleID:   1,
	}
}

func createUserForUpdate() *dto.UserUpdateDto {
	return &dto.UserUpdateDto{
		Email:       "test@test.com",
		Phone:       "+79999999999",
		Active:      true,
		Surname:     "Тестов",
		Name:        "Тест",
		Patronymic:  "Тестович",
		Sex:         "Man",
		BirthDate:   time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Profession:  "Кардиолог",
		Info:        "Информация ABCD0",
		RoleID:      2,
		Departments: []uint{1},
		Services:    []uint{1},
	}
}

func createUserForUpdateBindError() *UserDtoForUpdateBindError {
	return &UserDtoForUpdateBindError{
		Email:      "test@test.com",
		Phone:      "+79999999999",
		Active:     "true",
		Surname:    "Тестов",
		Name:       "Тест",
		Patronymic: "Тестович",
		Sex:        "Man",
		BirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Profession: "Кардиолог",
		Info:       "Информация ABCD0",
		RoleID:     "role",
		Department: "department",
		Services:   "services",
	}
}

func createResultUserForUpdateBindError() *dto.UserUpdateDto {
	return &dto.UserUpdateDto{
		Email:      "test@test.com",
		Phone:      "+79999999999",
		Active:     false,
		Surname:    "Тестов",
		Name:       "Тест",
		Patronymic: "Тестович",
		Sex:        "Man",
		BirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Profession: "Кардиолог",
		Info:       "Информация ABCD0",
		RoleID:     0,
	}
}

func createUserForUpdateValidationError() *dto.UserUpdateDto {
	return &dto.UserUpdateDto{
		Email:      "test1com",
		Phone:      "8",
		Active:     true,
		Surname:    "Тестов1",
		Name:       "Тест*",
		Patronymic: "Тесто вич",
		Sex:        "Im a man",
		BirthDate:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.Local),
		Profession: "\tКардиолог",
		Info:       "Информация \n ABCD0",
		RoleID:     1,
	}
}

func createPasswordForUpdate() *dto.UpdatePasswordDto {
	return &dto.UpdatePasswordDto{
		OldPassword:     "Password2!",
		NewPassword:     "Password3!",
		ConfirmPassword: "Password3!",
	}
}

func createPasswordForNotEqualConfirmPassword() *dto.UpdatePasswordDto {
	return &dto.UpdatePasswordDto{
		OldPassword:     "Password2!",
		NewPassword:     "Password3!",
		ConfirmPassword: "Password4!",
	}
}

func createPasswordForEqualOldPassword() *dto.UpdatePasswordDto {
	return &dto.UpdatePasswordDto{
		OldPassword:     "Password2!",
		NewPassword:     "Password2!",
		ConfirmPassword: "Password2!",
	}
}

func makeSlugForUser(user *models.User) string {
	slug.MaxLength = 40
	slug.EnableSmartTruncate = false
	slug.CustomSub = map[string]string{
		"ь": "",
		"Ь": "",
		"ъ": "",
		"Ъ": "",
	}
	if user.Surname != "" || user.Name != "" || user.Patronymic != "" {
		return slug.Make(fmt.Sprintf("%s %s %s", user.Surname, user.Name, user.Patronymic))
	} else {
		return strconv.Itoa(int(user.ID))
	}
}
