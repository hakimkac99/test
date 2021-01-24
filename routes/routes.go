package routes

import (
	"net/http"
	user "test-CRUD/controllers/users"

	"github.com/gin-gonic/gin"
)

func StartGin() {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/add/users", user.InsertUsers)
		api.POST("/login", user.Login)
		api.DELETE("/delete/user/:id", user.DeleteUser)
		api.GET("/users/list", user.GetAllUser)
		api.GET("/user/:id", user.GetUser)
		api.PUT("/users/:id", user.UpdateUser)

	}
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	router.Run(":8000")
}
