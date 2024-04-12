package scheduler

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"jwt/models"
	"math/rand"
	"encoding/json"	
)

func SmallTalk(){
	str:=os.Getenv("PROMPTS")
	prompts:=strings.Split(str,"^")
	client := &http.Client{Timeout: 5*time.Second}	
	rs,_ := FetchWxUser()
	if len(rs)<1{
		fmt.Println("SmallTalk, 加载用户失败！")
	}
	for _,r := range(rs){
		if !strings.Contains(r.Wxid,"@chatroom") {
			continue
		}
		payload := models.WxMsg{
			Type:1,
			Content:`st @`+WxMe.Name+RandChoice(prompts),
			FromUser:r.Wxid,	//"35031914979@chatroom", //
			ToUser:WxMe.Wxid,
		}
		rb, err := json.Marshal(&payload)
		if err!=nil{
			fmt.Println("SmallTalk ", err)
		}
		_, err = client.Post("http://localhost:8080/wxmsg", "application/json", bytes.NewBuffer(rb))
		if err!=nil {
			fmt.Println("SmallTalk ", err)
		}
		fmt.Printf("\nSamllTalk : %v------------%v\n", payload, time.Now())
		time.Sleep(60*time.Second)
	}
}

func RandChoice[T any](s []T) T{
	rand.Seed(time.Now().UnixNano())
	id:=rand.Intn(len(s))
	return s[id]
}
