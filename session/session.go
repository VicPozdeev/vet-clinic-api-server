package session

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"gopkg.in/boj/redistore.v1"
	"vet-clinic/config"
	"vet-clinic/logging"
	"vet-clinic/models"
)

const (
	// Auth represents a string of session key.
	Auth = "Authorization"
	// User is the key of account data in the session.
	User = "User"
)

// Session represents an interface for accessing the session within the application.
type Session interface {
	Store() sessions.Store

	Get(c echo.Context) *sessions.Session
	Save(c echo.Context) error
	Delete(c echo.Context) error
	SetValue(c echo.Context, key string, value interface{}) error
	GetValue(c echo.Context, key string) string
	SetUser(c echo.Context, user *models.User) error
	GetUser(c echo.Context) *models.User
}

type GorillaSession struct {
	store sessions.Store
}

// NewSession is constructor.
func NewSession(logger logging.Logger, conf *config.Config) *GorillaSession {
	if !conf.Redis.Enabled {
		logger.Infof("use CookieStore for session")
		return &GorillaSession{sessions.NewCookieStore([]byte("secret"))}
	}

	logger.Infof("use redis for session")
	logger.Infof("Try redis connection")
	address := fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port)
	store, err := redistore.NewRediStore(conf.Redis.ConnectionPoolSize, "tcp", address, "", []byte("secret"))
	if err != nil {
		logger.Panicf("Failure redis connection, %s", err.Error())
	}
	logger.Infof(fmt.Sprintf("Success redis connection, %s", address))
	return &GorillaSession{store: store}
}

func (s *GorillaSession) Store() sessions.Store {
	return s.store
}

// Get returns a session for the current request.
func (s *GorillaSession) Get(c echo.Context) *sessions.Session {
	sess, _ := s.store.Get(c.Request(), Auth)
	return sess
}

// Save saves the current session.
func (s *GorillaSession) Save(c echo.Context) error {
	sess := s.Get(c)
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	}
	return s.saveSession(c, sess)
}

// Delete the current session.
func (s *GorillaSession) Delete(c echo.Context) error {
	sess := s.Get(c)
	sess.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	return s.saveSession(c, sess)
}

func (s *GorillaSession) saveSession(c echo.Context, sess *sessions.Session) error {
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return fmt.Errorf("error occurred while save session")
	}
	return nil
}

// SetValue sets a key and a value.
func (s *GorillaSession) SetValue(c echo.Context, key string, value interface{}) error {
	sess := s.Get(c)
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json marshal error while set value in session")
	}
	sess.Values[key] = string(bytes)
	return nil
}

// GetValue returns value of session.
func (s *GorillaSession) GetValue(c echo.Context, key string) string {
	sess := s.Get(c)
	if sess != nil {
		if v, ok := sess.Values[key]; ok {
			data, result := v.(string)
			if result && data != "null" {
				return data
			}
		}
	}
	return ""
}

func (s *GorillaSession) SetUser(c echo.Context, user *models.User) error {
	return s.SetValue(c, User, user)
}

func (s *GorillaSession) GetUser(c echo.Context) *models.User {
	if v := s.GetValue(c, User); v != "" {
		a := &models.User{}
		_ = json.Unmarshal([]byte(v), a)
		return a
	}
	return nil
}
