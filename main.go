package main

import (
	// "flag"
	// "fmt"
	"jwt/controllers"
	"jwt/initializers"
	"jwt/middleware"
	"jwt/scheduler"
	"github.com/gin-gonic/gin"
)

func init(){
	initializers.LoadEnv()
	initializers.ConnDb()
	// initializers.SyncDb()		//merge table 
}

func main() {
	// ph := flag.String("ph", "127.0.0.1", "代理ip")
	// pp := flag.Int("pp", 1081, "代理端口")
	// p :=flag.String("p", "127.0.0.1:1081", "http代理")
	// flag.Parse()
	// fmt.Println(*ph, *pp, *p)
	
	go scheduler.Scheduler()
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)
	r.POST("/wxmsg", controllers.WxMsg)
	r.Run() // listen and serve on 0.0.0.0:8080
}