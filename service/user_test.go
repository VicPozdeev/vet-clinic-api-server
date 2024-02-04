package service

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"testing"
	"time"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/test"
)

func TestFindUserByID_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestFullData(cont, true)

	s := NewUserService(cont)
	result, err := s.Get("2")

	assert.Equal(t, uint(2), result.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestFindUserBySlug_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestFullData(cont, true)

	s := NewUserService(cont)
	resultByID, _ := s.Get("2")
	resultBySlug, err := s.Get(resultByID.Slug)

	assert.Equal(t, uint(2), resultBySlug.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resultBySlug)
}

func TestFindUserByID_IdNotNumeric(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	result, err := s.Get("ABCD")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindUserByID_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	result, err := s.Get("9999")

	assert.Nil(t, result)
	assert.Error(t, err, "failed to fetch data")
}

func TestFindAllUsers_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	result, err := s.GetAll()

	assert.Len(t, result, 2)
	assert.NoError(t, err)
}

func TestCreateUser_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	userDto := createUserForCreate()
	_, err := s.Create(userDto)

	result, _ := s.Get("2")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, userDto.Username, result.Username)
	assert.Equal(t, userDto.RoleID, result.RoleID)
	assert.Equal(t, true, result.Active)
	assert.Equal(t, strconv.Itoa(int(result.ID)), result.Slug)
}

func TestCreateUser_NotRole(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	result, err := s.Create(createUserForNotRole())

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestUpdateUser_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	userDto := createUserForUpdate()
	_, err := s.Update(userDto, "2", true)

	result, _ := s.Get("2")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, userDto.Email, result.Email)
	assert.Equal(t, userDto.Phone, result.Phone)
	assert.Equal(t, userDto.Active, result.Active)
	assert.Equal(t, userDto.Surname, result.Surname)
	assert.Equal(t, userDto.Name, result.Name)
	assert.Equal(t, userDto.Patronymic, result.Patronymic)
	assert.Equal(t, userDto.Sex, result.Sex)
	assert.Equal(t, userDto.Profession, result.Profession)
	assert.Equal(t, userDto.Info, result.Info)
	assert.Equal(t, makeSlugForUser(userDto.ToModel(true)), result.Slug)
	assert.Equal(t, userDto.RoleID, result.RoleID)
	assert.Equal(t, userDto.Departments[0], result.Departments[0].ID)
	assert.Equal(t, userDto.Services[0], result.Services[0].ID)
}

func TestUpdateUserNotOwner_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	userDto := createUserForUpdate()
	_, err := s.Update(userDto, "2", false)

	result, _ := s.Get("2")

	assert.NotEmpty(t, result)
	assert.Empty(t, err)
	assert.Equal(t, userDto.Email, result.Email)
	assert.Equal(t, userDto.Phone, result.Phone)
	assert.Equal(t, userDto.Active, result.Active)
	assert.Equal(t, userDto.Surname, result.Surname)
	assert.Equal(t, userDto.Name, result.Name)
	assert.Equal(t, userDto.Patronymic, result.Patronymic)
	assert.Equal(t, userDto.Sex, result.Sex)
	assert.NotEqual(t, userDto.Profession, result.Profession)
	assert.NotEqual(t, userDto.Info, result.Info)
	assert.Equal(t, makeSlugForUser(userDto.ToModel(true)), result.Slug)
	assert.NotEqual(t, userDto.RoleID, result.RoleID)
	assert.Empty(t, result.Departments)
	assert.Empty(t, result.Services)
}

func TestUpdateUser_NotEntity(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	result, err := s.Update(createUserForUpdate(), "99", true)

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestDeleteUser_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	data, _ := s.Get("2")

	result, err := s.Delete("2")

	assert.Equal(t, data, result)
	assert.Empty(t, err)
}

func TestDeleteUser_Error(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	result, err := s.Delete("99")

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestCompareHashAndPassword_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	s := NewUserService(cont)
	userDto := createUserForCreate()

	_, _ = s.Create(userDto)
	data, _ := s.Get("2")

	err := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(userDto.Password))

	assert.Empty(t, err)
}

func TestUpdatePassword_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	userDto := createUpdatePasswordDto()
	updateErr := s.UpdatePassword(userDto, "2")

	data, _ := s.Get("2")
	compareErr := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(userDto.NewPassword))

	assert.Empty(t, updateErr)
	assert.Empty(t, compareErr)
}

func TestUpdatePassword_WrongPassword(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	userDto := createUpdatePasswordDto()
	userDto.OldPassword = "wrong_password"
	err := s.UpdatePassword(userDto, "2")

	assert.Equal(t, "crypto/bcrypt: hashedPassword is not the hash of the given password", err.Error())
}

func TestLoginUserByUsername_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestFullData(cont, true)

	s := NewUserService(cont)
	result, err := s.Login(createLoginDto("Test2"))

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestLoginUserByEmail_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestFullData(cont, true)

	s := NewUserService(cont)
	result, err := s.Login(createLoginDto("test@test.com"))

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestLoginUserByPhone_Success(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestFullData(cont, true)

	s := NewUserService(cont)
	result, err := s.Login(createLoginDto("+79999999999"))

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestLoginUser_EntityNotFound(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	result, err := s.Login(createLoginDto("ABCD"))

	assert.Nil(t, result)
	assert.Equal(t, "record not found", err.Error())
}

func TestLoginUser_WrongPassword(t *testing.T) {
	cont := test.PrepareForServiceTest()

	setUpUserTestData(cont)

	s := NewUserService(cont)
	userDto := createLoginDto("Test1")
	userDto.Password = "pass1234"
	result, err := s.Login(userDto)

	assert.Nil(t, result)
	assert.Equal(t, "crypto/bcrypt: hashedPassword is not the hash of the given password", err.Error())
}

func setUpUserTestData(container container.Container) {
	rep := container.Repository()
	user := models.NewUser("Test2", "password", 1)
	_, _ = user.Create(rep)
}

func setUpUserTestFullData(container container.Container, owner bool) {
	setUpUserTestData(container)
	rep := container.Repository()

	user := createUserForUpdate().ToModel(owner)
	_, _ = user.Update(rep, 2, owner)
}

func createUserForCreate() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "Test2",
		Password: "password",
		RoleID:   1,
	}
}

func createUserForNotRole() *dto.UserCreateDto {
	return &dto.UserCreateDto{
		Username: "Test2",
		Password: "password",
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

func createLoginDto(username string) *dto.LoginDto {
	return &dto.LoginDto{
		Login:    username,
		Password: "password",
	}
}

func createUpdatePasswordDto() *dto.UpdatePasswordDto {
	return &dto.UpdatePasswordDto{
		OldPassword: "password",
		NewPassword: "new_password",
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
