package controllers

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gomysqlapp/appconfig"
	"gomysqlapp/auth"
	"gomysqlapp/database"
	"gomysqlapp/models"
	"gorm.io/gorm"
	"log"
)

// LoginPayload login body
// LoginPayload is a struct that contains the fields for a user's login credentials
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lp LoginPayload) LoginValidate() error {
	return validation.ValidateStruct(&lp,
		validation.Field(&lp.Email, validation.Required, is.Email),
		validation.Field(&lp.Password, validation.Required, validation.Length(6, 20)),
	)
}

type SignUpPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (sp SignUpPayload) RegisterValidate() error {

	return validation.ValidateStruct(&sp,
		validation.Field(&sp.Name, validation.Required),
		validation.Field(&sp.Email, validation.Required, is.Email),
		validation.Field(&sp.Password, validation.Required, validation.Length(6, 20)),
	)
}

// LoginResponse token response
// LoginResponse is a struct that contains the fields for a user's login response
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

// Signup is a function that handles user signup
// It takes in a gin context as an argument and binds the user data from the request body to a user struct
// It then hashes the user's password and creates a user record in the database
// If successful, it returns a 200 status code with a success message
// If unsuccessful, it returns a 400 or 500 status code with an error message

func Signup(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:    400,
			Context: c,
			Err:     err,
		})
		return
	}
	payload := SignUpPayload{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	if err = payload.RegisterValidate(); err != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:    400,
			Context: c,
			Err:     err,
		})
		return
	}

	err = user.HashPassword(user.Password)
	if err != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:    500,
			Context: c,
			Err:     err,
		})

		return
	}

	result := database.GlobalDB.Where("email = ?", payload.Email).First(&user)
	if result != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:     401,
			Context:  c,
			Messages: []string{"User with email address " + payload.Email + " already exist"},
		})
		return
	}

	createUser := user.CreateUserRecord()
	if createUser.Error != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:    401,
			Context: c,
			Err:     createUser.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"Message": "User Created Successfully",
	})
}

// Login is a function that handles user login
// It takes in a gin context as an argument and binds the user data from the request body to a LoginPayload struct
// It then checks if the user exists in the database and if the password is correct
// If successful, it generates a token and a refresh token and returns a 200 status code with the token and refresh token
// If unsuccessful, it returns a 401 or 500 status code with an error message

func Login(c *gin.Context) {
	var payload LoginPayload
	var user models.User
	errBind := c.ShouldBindJSON(&payload)
	if errBind != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:    400,
			Context: c,
			Err:     errBind,
		})
		return
	}
	if err := payload.LoginValidate(); err != nil {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:    400,
			Context: c,
			Err:     err,
		})
		return
	}
	result := database.GlobalDB.Where("email = ?", payload.Email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:     401,
			Context:  c,
			Messages: []string{"Invalid Username and Password"},
		})
		return
	}
	err := user.CheckPassword(payload.Password)
	if err != nil {
		log.Println(err)
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:     401,
			Context:  c,
			Messages: []string{"Invalid Username and Password"},
		})
		return
	}
	jwtWrapper := auth.JwtWrapper{
		SecretKey:       appconfig.GetEnvParam("JWT_SECRET"),
		Issuer:          appconfig.GetEnvParam("JWT_ISSUER"),
		ExpirationHours: 720,
	}
	signedToken, err := jwtWrapper.GenerateToken(user.Email)
	if err != nil {
		log.Println(err)
		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:     400,
			Context:  c,
			Messages: []string{"Error Signing Token"},
		})
		return
	}
	signedtoken, err := jwtWrapper.RefreshToken(user.Email)
	if err != nil {
		log.Println(err)

		appconfig.CustomErrResponse(appconfig.CustomErrorParams{
			Code:     400,
			Context:  c,
			Messages: []string{"Error Signing Token"},
		})
		return
	}
	tokenResponse := LoginResponse{
		Token:        signedToken,
		RefreshToken: signedtoken,
	}
	c.JSON(200, tokenResponse)
}
