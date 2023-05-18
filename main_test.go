package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"otus-sonet/controllers"
	"testing"
)

func TestGetUser(t *testing.T) {
	var url = "/health"
	r := gin.Default()
	r.GET(url, controllers.Health)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "{\"status\":\"OK\"}")
}
