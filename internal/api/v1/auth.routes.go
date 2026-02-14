package v1

import (
	"centralService/internal/controllers/authcontroller"
	"centralService/internal/controllers/usercontroller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, auth_controller *authcontroller.AuthController) {

	auth := router.Group("/auth")
	{
		auth.POST("/callback/google", auth_controller.Login)
		auth.GET("/me", auth_controller.Me)
		auth.GET("/generate-token", auth_controller.GenerateToken)
	}
}

func RegisterUserRoutes(router *gin.RouterGroup, usercontroller *usercontroller.UserController) {

	user := router.Group("/admin")
	{
		user.POST("/users", usercontroller.Create)
		user.GET("/users", usercontroller.List)
	}
}
