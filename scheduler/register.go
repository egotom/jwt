package scheduler

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"jwt/models"
	"encoding/json"	
	"jwt/initializers"
)

func WxRegister(){	
	Users := []models.User{}
	initializers.DB.Find(&Users)
	if len(Users)<10{
		panic("加载用户失败！")
	}
	// fmt.Printf("\nUsers : %v--------%v\n", len(Users), time.Now())
	client := &http.Client{Timeout: 3*time.Second}	
	requestBody := new(bytes.Buffer)
	res, err := client.Post("http://localhost:1909/api/getContactList", "application/json", requestBody)
	if err!=nil {
		fmt.Println(err)
		return
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	type User struct{
		Name 			string	 `json:"nickname"`
		Wxid			string 	 `json:"wxid"`
		Account			string 	 `json:"customAccount"`
	}
	type Body struct{
		Code 	int
		Data 	[]User 
		Msg 	string
	}
	var result Body
	json.Unmarshal(body, &result)
	for _, v := range(result.Data){
		if strings.Contains(v.Wxid,"@openim") {
			continue
		}
		bFond:=false
		for _,u:=range(Users){
			if u.Wxid==v.Wxid{
				bFond=true
				break
			}
		}
		if !bFond {
			initializers.DB.Create(&v)
			Ticket[v.Wxid]=[]int64{5000,5000}
			fmt.Printf("\nUsers : %v--------%v\n", v, time.Now())
		}
	}
}
