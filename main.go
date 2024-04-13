package main

import (
	"github.com/gin-gonic/gin"
	"jwt/initializers"
	"jwt/controllers"
	"jwt/middleware"
	"jwt/scheduler"
	// "jwt/output"
	"fmt"
)

func init(){
	initializers.LoadConfig("config.yaml")
	initializers.ConnDb()
	// initializers.SyncDb()		//merge table 
}

func main() {
	go scheduler.Scheduler()
	// go output.Honey()
	gin.SetMode(initializers.Config.GinMode)
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/wxmsg", controllers.WxMsg)
	s:=fmt.Sprint("0.0.0.0:",initializers.Config.Port)
	r.Run(s) // listen and serve on 0.0.0.0:8080
}