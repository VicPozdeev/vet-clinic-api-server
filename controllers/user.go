package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"vet-clinic/container"
	"vet-clinic/models"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type UserController struct {
	container container.Container
	service   *service.UserService
}

// NewUserController is constructor.
func NewUserController(container container.Container) *UserController {
	return &UserController{container: container, service: service.NewUserService(container)}
}

// Get returns one record matched user's id or slug.
//
// @Summary Get a user.
// @Description Returns one record matched user's id or slug.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id_or_slug path string true "User ID"
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /users/{id_or_slug} [get]
func (u *UserController) Get(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	user, err := u.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, user)
}

// GetSelf retrieves and returns the user's own profile based on the session.
//
// @Summary Get user's profile.
// @Description Retrieves and returns the user's own profile based on the session.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /profile [get]
func (u *UserController) GetSelf(c echo.Context) error {
	user := getUser(c, u.container)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, user)
}

// GetAll returns the list of users.
//
// @Summary Get a user list.
// @Description Returns the list of users.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.User "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /users [get]
func (u *UserController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	users, err := u.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, users)
}

// Create creates a new user.
//
// @Summary Create a new user. Required user's role: Owner
// @Description Create a new user.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.UserCreateDto true "A new user data for creating."
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 400 {object} dto.UserCreateDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /users [post]
func (u *UserController) Create(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.UserCreateDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	user, err := u.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, user)
}

// Update updates the existing user.
//
// @Summary Update the existing user. Required user's role: Owner
// @Description Update the existing user.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Param data body dto.UserUpdateDto true "User data for update."
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 400 {object} dto.UserUpdateDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /users/{id} [put]
func (u *UserController) Update(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.UserUpdateDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	result, err := u.service.Update(data, c.Param("id"), true)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, result)
}

// UpdateSelf updates and returns the user's own profile based on the session.
//
// @Summary Update user's profile.
// @Description Updates and returns the user's own profile based on the session.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.UserUpdateDto true "User data for update."
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 400 {object} dto.UserUpdateDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Router /profile [put]
func (u *UserController) UpdateSelf(c echo.Context) error {
	user := getUser(c, u.container)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	data := &dto.UserUpdateDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	result, err := u.service.Update(data, strconv.Itoa(int(user.ID)), false)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, result)
}

// UpdatePassword updates the user's password based on the provided data.
//
// @Summary Update user's password.
// @Description Updates the user's password based on the provided data.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.UpdatePasswordDto true "Old password and new password for update."
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 400 {object} dto.UpdatePasswordDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Router /profile/password [put]
func (u *UserController) UpdatePassword(c echo.Context) error {
	user := getUser(c, u.container)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	data := &dto.UpdatePasswordDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	err := u.service.UpdatePassword(data, strconv.Itoa(int(user.ID)))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, user)
}

// Delete deletes the existing user.
//
// @Summary Delete the existing user. Required user's role: Superuser
// @Description Delete the existing user.
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.User "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /users/{id} [delete]
func (u *UserController) Delete(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	user, err := u.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, user)
}

// Login is the method to login using username, e-mail, or phone along with the password.
//
// @Summary Login with credentials.
// @Description Login using username, e-mail, or phone along with the password.
// @Tags Users
// @Accept json
// @Produce json
// @Param data body dto.LoginDto true "Login and Password for logged-in."
// @Success 200 {object} models.User "Success to the authentication."
// @Header 200 {string} Cookie "Authorization"
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /login [post]
func (u *UserController) Login(c echo.Context) error {
	data := &dto.LoginDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	sess := u.container.Session()
	if user := sess.GetUser(c); user != nil {
		return c.JSON(http.StatusOK, user)
	}

	if user, err := u.service.Login(data); err == nil {
		_ = sess.SetUser(c, user)
		_ = sess.Save(c)
		return c.JSON(http.StatusOK, user)
	} else {
		return c.JSON(http.StatusUnauthorized, Error(err))
	}
}

// Logout logs the user out by invalidating the current session.
//
// @Summary Logout user.
// @Description Logout the user by invalidating the current session.
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 "Successfully logged out."
// @Failure 401 "Failed to the authentication."
// @Router /logout [post]
func (u *UserController) Logout(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	sess := u.container.Session()
	_ = sess.SetUser(c, nil)
	_ = sess.Delete(c)
	return c.NoContent(http.StatusOK)
}

func getAccessLevel(c echo.Context, container container.Container) util.AccessLevel {
	user := getUser(c, container)
	if user == nil {
		user = &models.User{Role: models.NewRole(util.Unauthorized.ToString())}
	}
	return util.ToAccessLevel(user.Role.Name)
}

func getUser(c echo.Context, container container.Container) *models.User {
	var user *models.User
	if user = container.Session().GetUser(c); user != nil {
		_ = container.Session().Save(c)
		return user
	}
	return nil
}
