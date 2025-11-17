package routes

import (
	"gin-app/controllers"
	"gin-app/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	user.POST("/create", controllers.CreateUser)
	user.POST("/login", controllers.LoginUser)
	protectedUser := r.Group("/api")
	protectedUser.Use(middleware.AuthMiddleware())
	protectedUser.GET("/me", controllers.UserProfile)
}
