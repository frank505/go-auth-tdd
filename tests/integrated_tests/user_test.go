package integrated_tests

import (
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"gomysqlapp/router"
	"gomysqlapp/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserProfile struct {
	User tests.UserData
}

func TestViewProfile(t *testing.T) {
	user := tests.UserData{
		Name:     gofakeit.Name(),
		Email:    gofakeit.Email(),
		Password: "password",
	}
	userToken := tests.GetJWTToken(user)

	req, err := http.NewRequest("GET", "/api/protected/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+userToken)
	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()
	r := router.SetupRouter()

	// Serve the test request using the test router
	r.ServeHTTP(rr, req)
	// Define the expected struct with the known user values

	// Create a struct to decode the response body into
	var actualResponse UserProfile

	// Decode the response body into the actualResponse struct
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Fatal(err)
	}

	// Use an assertion library (e.g., testify) to compare the structs
	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, actualResponse.User.Email, user.Email)
	assert.Equal(t, actualResponse.User.Name, user.Name)
}
