package appconfig

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

type CustomErrorParams struct {
	Code     int
	Err      interface{}
	Messages []string
	Context  *gin.Context
	// Add more fields as needed
}

func CustomErrResponse(params CustomErrorParams) {
	//log here
	if params.Err != nil {
		// Handle the validation error
		if validationErrors, ok := params.Err.(validator.ValidationErrors); ok {
			// Create a map to hold validation error messages
			errorMessages := make(map[string]string)

			// Loop through validationErrors to collect error messages
			for _, err := range validationErrors {
				fieldName := err.Field()
				errorMessage := fmt.Sprintf(err.Error())
				errorMessages[fieldName] = errorMessage
			}
			// continue by handling mysql errors

			params.Context.JSON(params.Code, gin.H{
				"Errors": errorMessages,
			})

			params.Context.Abort()
			return
		} else if mysqlErr, ok := params.Err.(*mysql.MySQLError); ok {
			fmt.Println("mysqlError", mysqlErr)

			params.Context.JSON(params.Code, gin.H{
				"Errors": "Failed to create user",
			})

			params.Context.Abort()
			return
		} else if len(params.Messages) < 1 {
			params.Context.JSON(params.Code, gin.H{
				"Errors": params.Err,
			})

			params.Context.Abort()
			return
		}

	}

	if len(params.Messages) > 0 {
		params.Context.JSON(params.Code, gin.H{
			"Error": params.Messages[0],
		})
		params.Context.Abort()
		return
	}

}
