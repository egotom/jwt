package controllers

import (
	"fmt"
	"regexp"
	// "time"
	"strings"
	"net/http"	
	// "jwt/openai"
	"jwt/google"
	"jwt/models"
	"jwt/output"
	"jwt/scheduler"	
	"github.com/gin-gonic/gin"
)

func WxMsg(c *gin.Context){
	var body models.WxMsg
	if err:=c.Bind(&body); err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"参数错误！"})
		fmt.Println(err)
		return 
	}
	if body.Type == 1{
		fmt.Println("content: ",body.Content)
		fmt.Println("\n\nfrom: ",body.FromUser)
		fmt.Println("to: ",body.ToUser)
		fmt.Println("type: ",body.Type)
		fmt.Println("\n\n\n")
	}
	to:="filehelper"
	if body.ToUser != "filehelper"{
		to=body.FromUser
	}
	if to=="weixin" ||  body.Type!=1 {
		return
	}
	ticket, exist := scheduler.Ticket[to] 
	if !exist {
		scheduler.WxRegister()
		if len(body.Content)>0 && !strings.Contains(body.FromUser,"@"){
			output.Reply(to, "您好，我可以回答问题、提供信息、与您聊天和提供娱乐。如果您需要帮助，请随时告诉我。")
		}
		// time.Sleep(5*time.Second)
		return
	}

	if (body.ToUser ==scheduler.WxMe.Wxid || body.ToUser == "filehelper") && !strings.Contains(body.FromUser,"@chatroom") {
		if output.CMDS(body.Content, to, ticket[0]){
			return
		}
		// go openai.ChatGPT(to, body.Content)
		go google.ChatGemini(to, body.Content)
		return
	}

	if strings.Contains(body.FromUser,"@chatroom") && strings.Contains(body.Content, "@"+scheduler.WxMe.Name){
		re := regexp.MustCompile(`^(\w+):(\s)@`+scheduler.WxMe.Name+".?")
		content:=re.ReplaceAllString(body.Content, "")
		if output.CMDS(content, to, ticket[0]){
			return
		}
		// go openai.ChatGPT(to, body.Content)
		go google.ChatGemini(to, content)
		return
	}
}
