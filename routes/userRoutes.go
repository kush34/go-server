package routes

import (
	"github.com/gin-gonic/gin"
	"gin-app/controllers"
)

func UserRoutes(r *gin.Engine) {
	r.GET("/users", controllers.GetUsers)
	r.POST("/users", controllers.CreateUser)
}
