package main

import (
	"github.com/gin-gonic/gin"
	"gin-app/routes"
)

func main() {
	r := gin.Default()

	routes.UserRoutes(r)

	r.Run(":8080")
}
