package main

import (
	"github.com/gin-gonic/gin"
	"kettkal/controllers"
	"kettkal/inits"
	"kettkal/middleware"
	"os"
)

func init() {
	inits.LoadEnvVariables()
	inits.ConnectToDB()
	inits.SyncDB()
}

func main() {
	port := os.Getenv("PORT")
	r := gin.New()
	r.GET("/", controllers.Welcome)
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.POST("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/generatePass", middleware.RequireAuth, controllers.GeneratePass)
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080
}
