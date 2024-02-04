package test

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"strings"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/logging"
	"vet-clinic/middleware"
	"vet-clinic/migration"
	"vet-clinic/models"
	"vet-clinic/repository"
	"vet-clinic/session"
	"vet-clinic/validate"

	"github.com/labstack/echo/v4"
)

// PrepareForControllerTest func prepares the controllers for testing.
func PrepareForControllerTest() (*echo.Echo, container.Container) {
	e := echo.New()
	e.Validator = validate.NewValidator(validator.New())

	conf := createConfig()
	logger := initTestLogger(false)
	cont := initContainer(conf, logger)

	middleware.Init(e, cont)

	migration.CreateTables(cont)
	migration.InitMasterData(cont)
	return e, cont
}

// PrepareForServiceTest func prepares the services for testing.
func PrepareForServiceTest() container.Container {
	conf := createConfig()
	logger := initTestLogger(false)
	cont := initContainer(conf, logger)

	migration.CreateTables(cont)
	migration.InitMasterData(cont)

	return cont
}

func createConfig() *config.Config {
	conf := &config.Config{}
	conf.Database.Dialect = "sqlite3"
	conf.Database.Host = "file::memory:?cache=shared"
	conf.Database.Migration = true
	conf.Extension.MasterGenerator = true

	return conf
}

func initTestLogger(debug bool) logging.Logger {
	zapConf := zap.NewDevelopmentConfig()
	if debug {
		zapConf.Level.SetLevel(zapcore.DebugLevel)
	} else {
		zapConf.Level.SetLevel(zapcore.PanicLevel)
	}
	zapLogger, err := zapConf.Build()
	if err != nil {
		fmt.Println("Failed to create zap new development logger")
	}
	logger := logging.NewLogger(zapLogger.Sugar(), &gormLogger.Config{})
	return logger
}

func initContainer(conf *config.Config, logger logging.Logger) container.Container {
	rep := repository.NewRepository(logger, conf)
	sess := session.NewSession(logger, conf)

	cont := container.NewContainer(rep, sess, conf, logger)
	return cont
}

// ConvertToJSON func converts model to string.
func ConvertToJSON(model interface{}) string {
	bytes, _ := json.Marshal(model)
	return string(bytes)
}

// NewJSONRequest func creates a new request using JSON format.
func NewJSONRequest(method string, target string, param interface{}) *http.Request {
	req := httptest.NewRequest(method, target, strings.NewReader(ConvertToJSON(param)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	return req
}

// LoginUser saves the user in the session.
func LoginUser(e *echo.Echo, container container.Container,
	r *http.Request, w http.ResponseWriter, user *models.User) {
	ctx := e.NewContext(r, w)
	_ = container.Session().SetUser(ctx, user)
}

// SetParam sets the parameter in the URI.
func SetParam(s string, param string) string {
	return strings.Replace(s, ":id", param, 1)
}
