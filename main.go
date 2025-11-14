package main

import (
	"gin-app/config"
	"gin-app/controllers"
	"gin-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	db := config.ConnectDB()
	controllers.InitUserController(db)
	routes.UserRoutes(r)

	r.Run(":8080")
}
