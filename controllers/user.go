package controllers

import (
	"jwt/initializers"
	"jwt/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context){
	var body struct{
		Email string
		Password string
	}
	if c.Bind(&body) !=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"参数错误！",
		})
		return 
	}

	hash,err:= bcrypt.GenerateFromPassword([]byte(body.Password),10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"参数错误！",
		})
		return
	}

	user:= models.User{Email:body.Email, Wxid:"filehelper", Password:string(hash)}
	result:=initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"创建用户失败！",
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{})
}


func Login(c *gin.Context){
	var body struct{
		Email string
		Password string
	}
	if c.Bind(&body) !=nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"参数错误！",
		})
		return 
	}
	var user models.User
	initializers.DB.First(&user,"email = ?", body.Email)
	if user.ID == 0{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"账号或密码无效！",
		})
		return
	}
	err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(body.Password))
	if err !=nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"账号或密码无效！",
		})
		return
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":user.ID,
		"exp":time.Now().Add(time.Hour * 24 *24*30).Unix(),
	})

	tokenString,err:=token.SignedString([]byte(os.Getenv("SECRET")))
	if err !=nil {
		c.JSON(http.StatusBadRequest,gin.H{
			"error":"生成 token 错误！",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"token":tokenString})
}

func Validate(c *gin.Context){
	user,_:=c.Get("user")
	c.JSON(http.StatusOK, gin.H{"message":user})
}