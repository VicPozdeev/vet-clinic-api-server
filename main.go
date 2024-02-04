package main

import (
	"embed"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/logging"
	"vet-clinic/middleware"
	"vet-clinic/migration"
	"vet-clinic/repository"
	"vet-clinic/router"
	"vet-clinic/session"
	"vet-clinic/validate"
)

//go:embed config/*.yml
var configFile embed.FS

// @title Vet clinic API
// @version v0.1.0
// @description This is API specification for vet-clinic API server.

// @contact.name Victor Pozdeev
// @contact.email vic.pozdeev@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/mit-license.php

// @host localhost:8080
// @BasePath /v1

// @schemes http
func main() {
	e := echo.New()
	e.HideBanner = true
	e.Validator = validate.NewValidator(validator.New())

	conf := config.LoadConfig(configFile)
	logger := logging.Init(conf)
	rep := repository.NewRepository(logger, conf)
	sess := session.NewSession(logger, conf)
	cont := container.NewContainer(rep, sess, conf, logger)

	middleware.Init(e, cont)
	router.Init(e, cont)

	migration.CreateTables(cont)
	migration.InitMasterData(cont)

	logger.Fatalf(e.Start(":8080").Error())
}
