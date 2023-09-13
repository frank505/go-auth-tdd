package appconfig

import (
	"github.com/joho/godotenv"
	"log"
)

func GetEnvParam(envParam string) (envValue string) {
	config, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error reading .env file")
	}
	return config[envParam]
}
