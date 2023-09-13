package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

// GlobalDB is a global db object that will be used across different packages
var GlobalDB *gorm.DB

func LoadEnvFile(filePath string) (map[string]string, error) {
	// Load the .env file
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if err := godotenv.Load(absPath); err != nil {
		return nil, err
	}

	// Read the environment variables into a map
	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}

	return envVars, nil
}

// InitDatabase creates a mysql db connection and stores it in the GlobalDB variable
func InitDatabase(state string) (err error) {
	// Read the environment variables from the .env file
	var dsn string

	// Create the data source name (DSN) using the environment variables
	if state == "testing" {

		// Load and read the .env file
		envVars, err := LoadEnvFile("../../.env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		dsn = fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			envVars["DB_USERNAME"],
			envVars["DB_PASSWORD"],
			envVars["DATABASE_HOST"],
			envVars["DB_DATABASE"],
		)
	} else {

		config, err := LoadEnvFile("./.env")
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

		dsn = fmt.Sprintf(
			"%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config["DB_USERNAME"],
			config["DB_PASSWORD"],
			config["DATABASE_HOST"],
			config["TEST_DB"],
		)
	}

	// Create the connection and store it in the GlobalDB variable
	GlobalDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	return
}
