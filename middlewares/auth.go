package middlewares

import (
	"github.com/gin-gonic/gin"
	"gomysqlapp/appconfig"
	"gomysqlapp/auth"
	"strings"
)

func Authz() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header from the request
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
			// If the Authorization header is not present, return a 403 status code
			appconfig.CustomErrResponse(appconfig.CustomErrorParams{
				Context:  c,
				Code:     403,
				Messages: []string{"No Authorization header provided"},
			})

			return
		}
		//// Split the Authorization header to get the token
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) == 2 {
			// Trim the token
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			appconfig.CustomErrResponse(appconfig.CustomErrorParams{
				Context:  c,
				Code:     400,
				Messages: []string{"No Authorization header provided"},
			})
			return
		}
		//// Create a JwtWrapper with the secret key and issuer
		jwtWrapper := auth.JwtWrapper{
			SecretKey:         appconfig.GetEnvParam("JWT_SECRET"),
			Issuer:            appconfig.GetEnvParam("JWT_ISSUER"),
			ExpirationHours:   720,
			ExpirationMinutes: 166640,
		}

		err, claims := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			appconfig.CustomErrResponse(appconfig.CustomErrorParams{
				Context:  c,
				Code:     403,
				Messages: []string{"error: " + err.Error()},
			})
			c.Abort()
			return
		}

		// Set the claims in the context
		c.Set("email", claims.Email)
		// Continue to the next handler
		c.Next()
	}
}
