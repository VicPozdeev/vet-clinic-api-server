package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"vet-clinic/container"
)

// APIError has a error code and a message.
type APIError struct {
	Code    int
	Message string
}

// ErrorController is a controller for handling errors.
type ErrorController interface {
	JSONError(err error, c echo.Context)
}

type DefaultErrorController struct {
	container container.Container
}

// NewErrorController is constructor.
func NewErrorController(container container.Container) *DefaultErrorController {
	return &DefaultErrorController{container: container}
}

// JSONError is cumstomize error handler
func (controller *DefaultErrorController) JSONError(err error, c echo.Context) {
	logger := controller.container.Logger()
	code := http.StatusInternalServerError
	msg := http.StatusText(code)

	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		msg = he.Message.(string)
	}

	var apiErr APIError
	apiErr.Code = code
	apiErr.Message = msg

	if !c.Response().Committed {
		if reserr := c.JSON(code, apiErr); reserr != nil {
			logger.Errorf(reserr.Error())
		}
	}
	logger.Debugf(err.Error())
}
