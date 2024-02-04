package router

import (
	"github.com/labstack/echo/v4"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/controllers"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "vet-clinic/docs" // for using echo-swagger
)

// Init initialize the routing of this application.
func Init(e *echo.Echo, container container.Container) {
	setSystemRoutes(e, container)
	setRoleRoutes(e, container)
	setUserRoutes(e, container)
	setDepartmentRoutes(e, container)
	setCategoryRoutes(e, container)
	setServiceRoutes(e, container)
	setClientRoutes(e, container)
	setPetRoutes(e, container)
	setVisitRoutes(e, container)
	setLeadRoutes(e, container)
}

func setSystemRoutes(e *echo.Echo, container container.Container) {
	system := controllers.NewSystemController(container)
	e.GET(config.APIv1Health, func(c echo.Context) error { return system.GetHealthCheck(c) })
	if container.Config().Swagger.Enabled {
		e.GET(container.Config().Swagger.Path, echoSwagger.WrapHandler)
	}
}

func setRoleRoutes(e *echo.Echo, container container.Container) {
	role := controllers.NewRoleController(container)
	e.GET(config.APIv1RolesID, func(c echo.Context) error { return role.Get(c) })
	e.GET(config.APIv1Roles, func(c echo.Context) error { return role.GetAll(c) })
	e.POST(config.APIv1Roles, func(c echo.Context) error { return role.Create(c) })
	e.PUT(config.APIv1RolesID, func(c echo.Context) error { return role.Update(c) })
	e.DELETE(config.APIv1RolesID, func(c echo.Context) error { return role.Delete(c) })
}

func setUserRoutes(e *echo.Echo, container container.Container) {
	user := controllers.NewUserController(container)
	e.GET(config.APIv1UsersID, func(c echo.Context) error { return user.Get(c) })
	e.GET(config.APIv1Users, func(c echo.Context) error { return user.GetAll(c) })
	e.POST(config.APIv1Users, func(c echo.Context) error { return user.Create(c) })
	e.PUT(config.APIv1UsersID, func(c echo.Context) error { return user.Update(c) })
	e.DELETE(config.APIv1UsersID, func(c echo.Context) error { return user.Delete(c) })
	e.GET(config.APIv1Profile, func(c echo.Context) error { return user.GetSelf(c) })
	e.PUT(config.APIv1Profile, func(c echo.Context) error { return user.UpdateSelf(c) })
	e.PUT(config.APIv1Password, func(c echo.Context) error { return user.UpdatePassword(c) })
	e.POST(config.APIv1Login, func(c echo.Context) error { return user.Login(c) })
	e.POST(config.APIv1Logout, func(c echo.Context) error { return user.Logout(c) })
}

func setDepartmentRoutes(e *echo.Echo, container container.Container) {
	department := controllers.NewDepartmentController(container)
	e.GET(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Get(c) })
	e.GET(config.APIv1Departments, func(c echo.Context) error { return department.GetAll(c) })
	e.POST(config.APIv1Departments, func(c echo.Context) error { return department.Create(c) })
	e.PUT(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Update(c) })
	e.DELETE(config.APIv1DepartmentsID, func(c echo.Context) error { return department.Delete(c) })
}

func setCategoryRoutes(e *echo.Echo, container container.Container) {
	category := controllers.NewCategoryController(container)
	e.GET(config.APIv1CategoriesID, func(c echo.Context) error { return category.Get(c) })
	e.GET(config.APIv1Categories, func(c echo.Context) error { return category.GetAll(c) })
	e.POST(config.APIv1Categories, func(c echo.Context) error { return category.Create(c) })
	e.PUT(config.APIv1CategoriesID, func(c echo.Context) error { return category.Update(c) })
	e.DELETE(config.APIv1CategoriesID, func(c echo.Context) error { return category.Delete(c) })
}

func setServiceRoutes(e *echo.Echo, container container.Container) {
	service := controllers.NewServiceController(container)
	e.GET(config.APIv1ServicesID, func(c echo.Context) error { return service.Get(c) })
	e.GET(config.APIv1Services, func(c echo.Context) error { return service.GetAll(c) })
	e.POST(config.APIv1Services, func(c echo.Context) error { return service.Create(c) })
	e.PUT(config.APIv1ServicesID, func(c echo.Context) error { return service.Update(c) })
	e.DELETE(config.APIv1ServicesID, func(c echo.Context) error { return service.Delete(c) })
}

func setClientRoutes(e *echo.Echo, container container.Container) {
	client := controllers.NewClientController(container)
	e.GET(config.APIv1ClientsID, func(c echo.Context) error { return client.Get(c) })
	e.GET(config.APIv1Clients, func(c echo.Context) error { return client.GetAll(c) })
	e.POST(config.APIv1Clients, func(c echo.Context) error { return client.Create(c) })
	e.PUT(config.APIv1ClientsID, func(c echo.Context) error { return client.Update(c) })
	e.DELETE(config.APIv1ClientsID, func(c echo.Context) error { return client.Delete(c) })
}

func setPetRoutes(e *echo.Echo, container container.Container) {
	pet := controllers.NewPetController(container)
	e.GET(config.APIv1PetsID, func(c echo.Context) error { return pet.Get(c) })
	e.GET(config.APIv1Pets, func(c echo.Context) error { return pet.GetAll(c) })
	e.POST(config.APIv1Pets, func(c echo.Context) error { return pet.Create(c) })
	e.PUT(config.APIv1PetsID, func(c echo.Context) error { return pet.Update(c) })
	e.DELETE(config.APIv1PetsID, func(c echo.Context) error { return pet.Delete(c) })
}

func setVisitRoutes(e *echo.Echo, container container.Container) {
	visit := controllers.NewVisitController(container)
	e.GET(config.APIv1VisitsID, func(c echo.Context) error { return visit.Get(c) })
	e.GET(config.APIv1Visits, func(c echo.Context) error { return visit.GetAll(c) })
	e.POST(config.APIv1Visits, func(c echo.Context) error { return visit.Create(c) })
	e.PUT(config.APIv1VisitsID, func(c echo.Context) error { return visit.Update(c) })
	e.DELETE(config.APIv1VisitsID, func(c echo.Context) error { return visit.Delete(c) })
}

func setLeadRoutes(e *echo.Echo, container container.Container) {
	lead := controllers.NewLeadController(container)
	e.GET(config.APIv1LeadsID, func(c echo.Context) error { return lead.Get(c) })
	e.GET(config.APIv1Leads, func(c echo.Context) error { return lead.GetAll(c) })
	e.POST(config.APIv1Leads, func(c echo.Context) error { return lead.Create(c) })
	e.PUT(config.APIv1LeadsID, func(c echo.Context) error { return lead.Update(c) })
	e.DELETE(config.APIv1LeadsID, func(c echo.Context) error { return lead.Delete(c) })
}
