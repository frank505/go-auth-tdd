package integrated_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gomysqlapp/database"
	"gomysqlapp/router"
	"gomysqlapp/tests"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := "password"
	bytesPassHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println("failed to hash password")
	}

	newUser := map[string]interface{}{
		"Name":       name,
		"Password":   string(bytesPassHash),
		"Email":      email,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	database.GlobalDB.Table("users").Create(newUser)
	bodyData := map[string]string{
		"Email":    email,
		"Password": "password",
	}
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req, err := http.NewRequest("POST", "/api/public/login", bytes.NewBuffer(jsonData))
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
	fmt.Println("Err", rr.Body)
	assert.Equal(t, 201, rr.Code)
}

func TestSignup(t *testing.T) {
	errorSetup := tests.Setup()
	if errorSetup != nil {
		fmt.Println("Error with setup", errorSetup)
	}

	bodyData := map[string]string{
		"Name":     gofakeit.Name(),
		"Password": "password",
		"Email":    gofakeit.Email(),
	}
	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	req, err := http.NewRequest("POST", "/api/public/signup", bytes.NewBuffer(jsonData))
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
	fmt.Println("Err", rr.Body)
	assert.Equal(t, 200, rr.Code)
}
