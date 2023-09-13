package router

import (
	"github.com/gin-gonic/gin"
	"gomysqlapp/controllers"
	"gomysqlapp/middlewares"
)

func SetupRouter() *gin.Engine {
	// Create a new router
	r := gin.Default()
	// Add a welcome route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome To This Website")
	})
	r.GET("/list-routes", func(c *gin.Context) {
		if gin.Mode() == gin.DebugMode {
			var routesList []string // Declare an empty slice
			for _, item := range r.Routes() {
				route := "method: " + item.Method + " path:" + item.Path
				routesList = append(routesList, route)
			}
			c.JSON(200, routesList)

		}

	})
	// Create a new group for the API
	api := r.Group("/api")
	{
		// Create a new group for the public routes
		public := api.Group("/public")
		{
			// Add the login route
			public.POST("/login", controllers.Login)
			// Add the signup route
			public.POST("/signup", controllers.Signup)
		}
		// Add the signup route
		protected := api.Group("/protected").Use(middlewares.Authz())
		{
			// Add the profile route
			protected.GET("/profile", controllers.Profile)
		}
	}
	// Return the router
	return r
}
