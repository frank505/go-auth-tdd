package tests

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gomysqlapp/appconfig"
	"gomysqlapp/auth"
	"gomysqlapp/database"
	"gomysqlapp/models"
)

type UserData struct {
	Name     string
	Email    string
	Password string
}

func Setup() error {
	// Load environment variables from .env
	err := database.InitDatabase("testing")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return err
	}
	// Automigrate the User model
	// AutoMigrate() automatically migrates our schema, to keep our schema up to date.
	errAutoMigrate := database.GlobalDB.AutoMigrate(&models.User{})
	if errAutoMigrate != nil {
		fmt.Println("Failed migration", errAutoMigrate)
		return errAutoMigrate
	}

	return nil
}

func CreateNewUser(user UserData) UserData {
	errorSetup := Setup()
	if errorSetup != nil {
		fmt.Println("Error with setup", errorSetup)
	}

	database.GlobalDB.Table("users").Create(&user)
	return user
}

func GetJWTToken(user UserData) string {
	password := user.Password
	bytes, errPassword := bcrypt.GenerateFromPassword([]byte(password), 14)
	if errPassword != nil {
		fmt.Println("Error:", errPassword)
	}
	user.Password = string(bytes)
	CreateNewUser(user)
	jwtWrapperParams := auth.JwtWrapper{
		SecretKey:         appconfig.GetEnvParam("JWT_SECRET"),
		Issuer:            appconfig.GetEnvParam("JWT_ISSUER"),
		ExpirationHours:   720,
		ExpirationMinutes: 166640,
	}

	signedToken, errorTokenSigning := jwtWrapperParams.GenerateToken(user.Email)
	if errorTokenSigning != nil {
		fmt.Println("Error:", errorTokenSigning)
	}
	return signedToken
}
