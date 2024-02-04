package middleware

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
	"github.com/valyala/fasttemplate"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"vet-clinic/config"
	"vet-clinic/container"
	"vet-clinic/logging"
)

// Init initializes the middleware for the application.
func Init(e *echo.Echo, container container.Container) {
	conf := container.Config()
	logger := container.Logger()

	// Recovery middleware
	e.Use(setErrorHandler(e, container))

	// Logger middleware
	e.Use(actionLoggerMiddleware(logger))
	e.Use(requestLoggerMiddleware(logger))

	// SecurityMiddleware
	if conf.Extension.SecurityEnabled {
		e.Use(setSecureMiddleware())
	}

	// CORS middleware
	if conf.Extension.CorsEnabled {
		e.Use(setCORSMiddleware())
	}

	// CSRF middleware
	if conf.Extension.CsrfEnabled {
		e.Use(setCSRFMiddleware())
	}

	// Session middleware
	e.Use(session.Middleware(container.Session().Store()))

	// Gzip middleware
	e.Use(echomw.Gzip())

	// Static middleware
	if conf.StaticContents.Enabled {
		e.Use(staticContentsMiddleware(conf))
		logger.Infof("Served the static contents.")
	}
}

func setErrorHandler(e *echo.Echo, container container.Container) echo.MiddlewareFunc {
	errorHandler := NewErrorController(container)
	e.HTTPErrorHandler = errorHandler.JSONError
	return echomw.Recover()
}

// requestLoggerMiddleware is middleware for logging the contents of requests.
func requestLoggerMiddleware(logger logging.Logger) echo.MiddlewareFunc {
	return echomw.RequestLoggerWithConfig(echomw.RequestLoggerConfig{
		LogRemoteIP: true,
		LogURI:      true,
		LogMethod:   true,
		LogStatus:   true,
		LogValuesFunc: func(c echo.Context, v echomw.RequestLoggerValues) error {
			template := fasttemplate.New("${remote_ip} ${uri} ${method} ${status}",
				"${", "}")
			result := template.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
				switch tag {
				case "remote_ip":
					return w.Write([]byte(v.RemoteIP))
				case "uri":
					return w.Write([]byte(v.URI))
				case "method":
					return w.Write([]byte(v.Method))
				case "status":
					return w.Write([]byte(strconv.Itoa(v.Status)))
				default:
					return w.Write([]byte(""))
				}
			})
			logger.Infof(result)
			return nil
		},
	})
}

// actionLoggerMiddleware is middleware for logging the start and end of controller processes.
func actionLoggerMiddleware(logger logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Debugf(c.Path() + " Action Start")
			if err := next(c); err != nil {
				c.Error(err)
			}
			logger.Debugf(c.Path() + " Action End")
			return nil
		}
	}
}

func setSecureMiddleware() echo.MiddlewareFunc {
	return echomw.SecureWithConfig(echomw.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	})
}

func setCORSMiddleware() echo.MiddlewareFunc {
	return echomw.CORSWithConfig(echomw.CORSConfig{
		AllowCredentials:                         true,
		UnsafeWildcardOriginWithAllowCredentials: true,
		AllowOrigins:                             []string{"*"},
		AllowHeaders: []string{
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAcceptEncoding,
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		MaxAge: 86400,
	})

}

func setCSRFMiddleware() echo.MiddlewareFunc {
	return echomw.CSRFWithConfig(echomw.CSRFConfig{
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})
}

// StaticContentsMiddleware is the middleware for loading the static files.
func staticContentsMiddleware(conf *config.Config) echo.MiddlewareFunc {
	staticConfig := echomw.StaticConfig{
		Root:   "assets",
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
		//Filesystem: http.FS(staticFile),
	}
	if conf.Swagger.Enabled {
		staticConfig.Skipper = func(c echo.Context) bool {
			return equalPath(c.Path(), []string{conf.Swagger.Path})
		}
	}
	return echomw.StaticWithConfig(staticConfig)
}

// equalPath judges whether a given path contains in the path list.
func equalPath(cpath string, paths []string) bool {
	for i := range paths {
		if regexp.MustCompile(paths[i]).Match([]byte(cpath)) {
			return true
		}
	}
	return false
}
