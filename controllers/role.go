package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type RoleController struct {
	container container.Container
	service   *service.RoleService
}

// NewRoleController is constructor.
func NewRoleController(container container.Container) *RoleController {
	return &RoleController{container: container, service: service.NewRoleService(container)}
}

// Get returns one record matched role's id.
//
// @Summary Get a role. Required user's role: Owner
// @Description Returns one record matched role's id.
// @Tags Roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Success 200 {object} models.Role "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /roles/{id} [get]
func (r *RoleController) Get(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	role, err := r.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, role)
}

// GetAll returns the list of roles.
//
// @Summary Get a role list. Required user's role: Owner
// @Description Returns the list of roles.
// @Tags Roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Role "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /roles [get]
func (r *RoleController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	roles, err := r.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, roles)
}

// Create creates a new role.
//
// @Summary Create a new role. Required user's role: Superuser
// @Description Create a new role.
// @Tags Roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.RoleDto true "A new role data for creating."
// @Success 200 {object} models.Role "Success to fetch data."
// @Failure 400 {object} dto.RoleDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /roles [post]
func (r *RoleController) Create(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.RoleDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	role, err := r.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, role)
}

// Update updates the existing role.
//
// @Summary Update the existing role. Required user's role: Superuser
// @Description Update the existing role.
// @Tags Roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Param data body dto.RoleDto true "Role data for update."
// @Success 200 {object} models.Role "Success to fetch data."
// @Failure 400 {object} dto.RoleDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /roles/{id} [put]
func (r *RoleController) Update(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.RoleDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	role, err := r.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, role)
}

// Delete deletes the existing role.
//
// @Summary Delete the existing role. Required user's role: Superuser
// @Description Delete the existing role.
// @Tags Roles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Role ID"
// @Success 200 {object} models.Role "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /roles/{id} [delete]
func (r *RoleController) Delete(c echo.Context) error {
	level := getAccessLevel(c, r.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Superuser.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	role, err := r.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, role)
}
