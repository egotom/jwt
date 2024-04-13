package middleware

import (
	"fmt"
	"jwt/initializers"
	"jwt/models"
	"log"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context){
	tokenString,err:=c.Cookie("Authorization") 
	// tokenString=c.Request.Header["Token"][0]
	if err !=nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return  []byte(initializers.Config.Secret), nil
	})
	if err != nil {
		log.Fatal(err,"-----------------",tokenString)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix())> claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var user models.User
		initializers.DB.First(&user,claims["sub"])
		if user.ID==0{
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user",user)
		c.Next()
	} else {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}	
}