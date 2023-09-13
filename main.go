package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gomysqlapp/database"
	"gomysqlapp/models"
	"gomysqlapp/router"
)

func initializeDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}

	return db, nil
}

func main() {

	err := database.InitDatabase("production")
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	// Automigrate the User model
	// AutoMigrate() automatically migrates our schema, to keep our schema upto date.
	errAutoMigrate := database.GlobalDB.AutoMigrate(&models.User{})
	if errAutoMigrate != nil {
		fmt.Println("fails migration", errAutoMigrate)
		return
	}
	// Set up the router
	r := router.SetupRouter()
	// Start the server
	r.Run(":8080")

	// Use the db object for database operations
}
