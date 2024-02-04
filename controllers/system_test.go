package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"vet-clinic/config"
	"vet-clinic/test"
)

func TestGetHealthCheck(t *testing.T) {
	e, cont := test.PrepareForControllerTest()

	System := NewSystemController(cont)
	e.GET(config.APIv1Health, func(c echo.Context) error { return System.GetHealthCheck(c) })

	req := httptest.NewRequest("GET", config.APIv1Health, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	expected := map[string]interface{}{
		"status": "available",
	}
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, test.ConvertToJSON(expected), rec.Body.String())
}
