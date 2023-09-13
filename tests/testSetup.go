package tests

import (
	"fmt"
	"gomysqlapp/database"
	"gomysqlapp/models"
)

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
