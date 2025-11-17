	package main

	import (
		"gin-app/config"
		"gin-app/controllers"
		"gin-app/routes"
		"log"
		"os"

		"github.com/gin-gonic/gin"
		"github.com/joho/godotenv"
	)

	func main() {
		loadEnvErr := godotenv.Load()
		if loadEnvErr != nil {
			log.Fatal("Error loading .env file")
		}
		PORT := os.Getenv("PORT")

		r := gin.Default()
		db := config.ConnectDB()
		controllers.InitUserController(db)
		routes.UserRoutes(r)

		r.Run(PORT)
	}
