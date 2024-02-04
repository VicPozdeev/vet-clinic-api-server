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

type PetDtoForBindError struct {
	Name     string
	Type     string
	Breed    string
	Colour   string
	Sex      string
	ClientID string
}

func TestGetPetByID_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpPetTestData(cont)

	pet := NewPetController(cont)
	e.GET(config.APIv1PetsID, func(c echo.Context) error { return pet.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1PetsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Pet{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetPet_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.GET(config.APIv1PetsID, func(c echo.Context) error { return pet.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1PetsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestGetPet_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.GET(config.APIv1PetsID, func(c echo.Context) error { return pet.Get(c) })

	req := httptest.NewRequest("GET", test.SetParam(config.APIv1PetsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestGetPetList_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpPetTestData(cont)

	pet := NewPetController(cont)
	e.GET(config.APIv1Pets, func(c echo.Context) error { return pet.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Pets, nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Pet{}
	data, _ := m.GetAll(cont.Repository())

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestGetPetList_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.GET(config.APIv1Pets, func(c echo.Context) error { return pet.GetAll(c) })

	req := httptest.NewRequest("GET", config.APIv1Pets, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestCreatePet_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.POST(config.APIv1Pets, func(c echo.Context) error { return pet.Create(c) })

	param := createPetForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Pets, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Pet{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestCreatePet_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.POST(config.APIv1Pets, func(c echo.Context) error { return pet.Create(c) })

	param := createPetForBindError()
	req := test.NewJSONRequest("POST", config.APIv1Pets, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultPetForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestCreatePet_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.POST(config.APIv1Pets, func(c echo.Context) error { return pet.Create(c) })

	param := createPetForValidationError()
	req := test.NewJSONRequest("POST", config.APIv1Pets, param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "PetDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "PetDto.Type")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "PetDto.Breed")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "PetDto.Colour")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "PetDto.Sex")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
}

func TestCreatePet_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.POST(config.APIv1Pets, func(c echo.Context) error { return pet.Create(c) })

	param := createPetForCreate()
	req := test.NewJSONRequest("POST", config.APIv1Pets, param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestUpdatePet_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpPetTestData(cont)

	pet := NewPetController(cont)
	e.PUT(config.APIv1PetsID, func(c echo.Context) error { return pet.Update(c) })

	param := createPetForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1PetsID, "2"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	m := &models.Pet{}
	data, _ := m.Get(cont.Repository(), 2)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestUpdatePet_BindError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.PUT(config.APIv1PetsID, func(c echo.Context) error { return pet.Update(c) })

	param := createPetForBindError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1PetsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	result := createResultPetForBindError()
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(result), rec.Body.String())
}

func TestUpdatePet_ValidationError(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.PUT(config.APIv1PetsID, func(c echo.Context) error { return pet.Update(c) })

	param := createPetForValidationError()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1PetsID, "1"), param)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Staff)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "PetDto.Name")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
	assert.Contains(t, rec.Body.String(), "PetDto.Type")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "PetDto.Breed")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "PetDto.Colour")
	assert.Contains(t, rec.Body.String(), "'ruprintascii'")
	assert.Contains(t, rec.Body.String(), "PetDto.Sex")
	assert.Contains(t, rec.Body.String(), "'rualpha'")
}

func TestUpdatePet_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.PUT(config.APIv1PetsID, func(c echo.Context) error { return pet.Update(c) })

	param := createPetForUpdate()
	req := test.NewJSONRequest("PUT", test.SetParam(config.APIv1PetsID, "1"), param)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeletePet_Success(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	setUpPetTestData(cont)

	pet := NewPetController(cont)
	e.DELETE(config.APIv1PetsID, func(c echo.Context) error { return pet.Delete(c) })

	m := &models.Pet{}
	data, _ := m.Get(cont.Repository(), 2)

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1PetsID, "2"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(data), rec.Body.String())
}

func TestDeletePet_Failure(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.DELETE(config.APIv1PetsID, func(c echo.Context) error { return pet.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1PetsID, "9999"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Superuser)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	expected := ErrorResponse{Message: "record not found"}
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}

func TestDeletePet_Unauthorized(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.DELETE(config.APIv1PetsID, func(c echo.Context) error { return pet.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1PetsID, "1"), nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestDeletePet_Forbidden(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	pet := NewPetController(cont)
	e.DELETE(config.APIv1PetsID, func(c echo.Context) error { return pet.Delete(c) })

	req := test.NewJSONRequest("DELETE", test.SetParam(config.APIv1PetsID, "1"), nil)
	rec := httptest.NewRecorder()

	userForLogin := userWithAccessLevel(util.Owner)
	test.LoginUser(e, cont, req, rec, userForLogin)

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
}

func setUpPetTestData(container container.Container) {
	rep := container.Repository()
	pet := createPetForCreate().ToModel()
	_, _ = pet.Create(rep)
}

func createPetForCreate() *dto.PetDto {
	return &dto.PetDto{
		Name:     "Шарик",
		Type:     "Пёс",
		Breed:    "Дворняга",
		Colour:   "Коричневый",
		Sex:      "Самец",
		ClientID: 1,
	}
}

func createPetForBindError() *PetDtoForBindError {
	return &PetDtoForBindError{
		Name:     "Шарик",
		Type:     "Пёс",
		Breed:    "Дворняга",
		Colour:   "Коричневый",
		Sex:      "Самец",
		ClientID: "Client",
	}
}

func createResultPetForBindError() *dto.PetDto {
	return &dto.PetDto{
		Name:     "Шарик",
		Type:     "Пёс",
		Breed:    "Дворняга",
		Colour:   "Коричневый",
		Sex:      "Самец",
		ClientID: 0,
	}
}

func createPetForValidationError() *dto.PetDto {
	return &dto.PetDto{
		Name:     "Шарик2",
		Type:     "Пёс\n",
		Breed:    "Дворняга\n",
		Colour:   "Коричневый\n",
		Sex:      "Самец2",
		ClientID: 1,
	}
}

func createPetForUpdate() *dto.PetDto {
	return &dto.PetDto{
		Name:     "Ракета",
		Type:     "Енот",
		Breed:    "Полоскун",
		Colour:   "Буровато-серый",
		Sex:      "Самец",
		ClientID: 1,
	}
}
