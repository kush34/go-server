package routes

import (
	"gin-app/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
    user := r.Group("/user")
    user.POST("/create", controllers.CreateUser)
    user.POST("/login", controllers.LoginUser)
}
