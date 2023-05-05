package test

import (
	"net/http"
	"net/http/httptest"
	"rabbit/routes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleUserNotification(t *testing.T) {
	router := routes.InitRoutes()

	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/notification/beyond", nil)
	if err != nil {
		t.Fatal("Error from test request setup", err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}
