package appconfig

import (
	"fmt"
	"gomysqlapp/database"
)

func GetEnvParam(envParam string) (envValue string) {
	var envVars map[string]string
	envVarsTest, err := database.LoadEnvFile("../../.env")

	if err != nil {
		envVarsproduction, errSecond := database.LoadEnvFile("./.env")
		if errSecond != nil {
			fmt.Println("Error env failed to open first time:", errSecond)
		} else {
			envVars = envVarsproduction
		}
	} else {
		envVars = envVarsTest
	}

	// Check if the environment variable exists
	if value, ok := envVars[envParam]; ok {
		return value
	}

	// Handle the case where the environment variable does not exist
	fmt.Printf("Environment variable '%s' not found\n", envParam)
	return "" // This line is only reached to satisfy Go's return requirements
}
