package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
	"vet-clinic/models/dto"
	"vet-clinic/service"
	"vet-clinic/util"
)

type DepartmentController struct {
	container container.Container
	service   *service.DepartmentService
}

// NewDepartmentController is constructor.
func NewDepartmentController(container container.Container) *DepartmentController {
	return &DepartmentController{container: container, service: service.NewDepartmentService(container)}
}

// Get returns one record matched department's id.
//
// @Summary Get a department.
// @Description Returns one record matched department's id.
// @Tags Departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id_or_slug path string true "Department ID"
// @Success 200 {object} models.Department "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /departments/{id_or_slug} [get]
func (u *DepartmentController) Get(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	department, err := u.service.Get(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, department)
}

// GetAll returns the list of departments.
//
// @Summary Get a department list.
// @Description Returns the list of departments.
// @Tags Departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} []models.Department "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to fetch data."
// @Failure 401 "Failed to the authentication."
// @Router /departments [get]
func (u *DepartmentController) GetAll(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}

	departments, err := u.service.GetAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, departments)
}

// Create creates a new department.
//
// @Summary Create a new department. Required user's role: Owner
// @Description Create a new department.
// @Tags Departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DepartmentDto true "A new department data for creating."
// @Success 200 {object} models.Department "Success to fetch data."
// @Failure 400 {object} dto.DepartmentDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /departments [post]
func (u *DepartmentController) Create(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.DepartmentDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	department, err := u.service.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, department)
}

// Update updates the existing department.
//
// @Summary Update the existing department. Required user's role: Owner
// @Description Update the existing department.
// @Tags Departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Department ID"
// @Param data body dto.DepartmentDto true "Department data for update."
// @Success 200 {object} models.Department "Success to fetch data."
// @Failure 400 {object} dto.DepartmentDto "Failed to the registration."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /departments/{id} [put]
func (u *DepartmentController) Update(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	data := &dto.DepartmentDto{}
	if err := c.Bind(data); err != nil {
		return c.JSON(http.StatusBadRequest, data)
	}
	if err := c.Validate(data); err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}

	result, err := u.service.Update(data, c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, result)
}

// Delete deletes the existing department.
//
// @Summary Delete the existing department. Required user's role: Owner
// @Description Delete the existing department.
// @Tags Departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Department ID"
// @Success 200 {object} models.Department "Success to fetch data."
// @Failure 400 {object} ErrorResponse "Failed to the registration."
// @Failure 401 "Failed to the authentication."
// @Failure 403 "Access denied."
// @Router /departments/{id} [delete]
func (u *DepartmentController) Delete(c echo.Context) error {
	level := getAccessLevel(c, u.container)
	if !util.Staff.AccessAllowed(level) {
		return c.NoContent(http.StatusUnauthorized)
	}
	if !util.Owner.AccessAllowed(level) {
		return c.NoContent(http.StatusForbidden)
	}

	department, err := u.service.Delete(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error(err))
	}
	return c.JSON(http.StatusOK, department)
}
