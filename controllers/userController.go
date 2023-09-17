package controllers

import (
	"github.com/gin-gonic/gin"
	"gomysqlapp/database"
	"gomysqlapp/models"
	"gorm.io/gorm"
)

type extendEmailToContext struct {
	c     *gin.Context
	email string
}

func Profile(c *gin.Context) {
	var user models.User
	//// Get the email from the authorization middleware
	email, _ := c.Get("email")
	result := database.GlobalDB.Where("email = ?", email.(string)).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		c.JSON(404, gin.H{
			"Error": "User not found",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"user": gin.H{
			"Name":       user.Name,
			"Email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
	return
}
