package routes

import (
	"go-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userController *controllers.UserController) {
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUser)
			users.POST("", userController.CreateUser)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
	}
}
