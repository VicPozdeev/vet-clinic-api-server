package container

import (
	"vet-clinic/config"
	"vet-clinic/logging"
	"vet-clinic/repository"
	"vet-clinic/session"
)

// Container represents a interface for accessing the data which sharing in overall application.
type Container interface {
	Repository() repository.Repository
	Session() session.Session
	Config() *config.Config
	Logger() logging.Logger
}

// DefaultContainer struct is for sharing data which such as database setting, the setting of application and logger in overall this application.
type DefaultContainer struct {
	rep     repository.Repository
	session session.Session
	config  *config.Config
	logger  logging.Logger
}

// NewContainer is constructor.
func NewContainer(rep repository.Repository, session session.Session, config *config.Config, logger logging.Logger) *DefaultContainer {
	return &DefaultContainer{rep: rep, session: session, config: config, logger: logger}
}

// Repository returns the object of repository.
func (c *DefaultContainer) Repository() repository.Repository {
	return c.rep
}

// Session returns the object of session.
func (c *DefaultContainer) Session() session.Session {
	return c.session
}

// Config returns the object of configuration.
func (c *DefaultContainer) Config() *config.Config {
	return c.config
}

// Logger returns the object of logger.
func (c *DefaultContainer) Logger() logging.Logger {
	return c.logger
}
