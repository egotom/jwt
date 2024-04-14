package scheduler

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"math/rand"
	"jwt/output"
	"jwt/models"
	"encoding/json"	
	"jwt/initializers"
)

func SmallTalk(style string){
	prompts:=initializers.Config.Prompts
	client := &http.Client{Timeout: 5*time.Second}	
	rs,_ := FetchWxUser()
	if len(rs)<1{
		fmt.Println("SmallTalk, 加载用户失败！")
		return
	}
	honey := []string{}
	if style=="honey"{
		initializers.DB.Table("honeys").Pluck("commit", &honey)
		if len(honey)<1{
			fmt.Println("SmallTalk , honeys 加载失败！")
			return
		}
	}
	for _,r := range(rs){
		if !strings.Contains(r.Wxid,"@chatroom") {
			continue
		}
		prompt := prompts["@chatroom"]
		if p, ok := prompts[r.Wxid]; ok{
			prompt=p
		}
		if len(prompt)<1 {
			continue
		}
		if style=="honey" {
			output.Reply(r.Wxid, RandChoice(honey))
		}else{
			payload := models.WxMsg{
				Type:1,
				Content:`st: @`+WxMe.Name+"-"+RandChoice(prompt),
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
			fmt.Printf("\nSamllTalk : %v------------%v\n", string(rb), time.Now())
			time.Sleep(3*time.Minute)
		}
	}
}

func RandChoice[T any](s []T) T{
	rand.Seed(time.Now().UnixNano())
	id:=rand.Intn(len(s))
	return s[id]
}
