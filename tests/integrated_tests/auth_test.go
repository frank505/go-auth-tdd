package integrated_tests

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gomysqlapp/router"
	"gomysqlapp/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWelcomeRoute(t *testing.T) {
	// Create a test request to the welcome route
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a test router and serve the request
	r := router.SetupRouter()
	r.ServeHTTP(rr, req)

	// Check the response status code and body
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Welcome To This Website")
}

func TestLoginApi(t *testing.T) {
	errorSetup := tests.Setup()
	if errorSetup != nil {
		fmt.Println("Error with setup", errorSetup)
	}
	bodyData := `{"email": "example@email.com", "password": "password123"}`
	req, err := http.NewRequest("POST", "/api/public/login", bytes.NewBufferString(bodyData))
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	r := router.SetupRouter()

	// Serve the test request using the test router
	r.ServeHTTP(rr, req)
	// Create a test HTTP request for login with JSON payload
	// Check the response status code and body
	assert.Equal(t, http.StatusOK, rr.Code)
}
