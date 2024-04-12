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
	if len(Users)<1{
		panic("加载用户失败！")
	}
	// fmt.Printf("\nUsers : %v--------%v\n", len(Users), time.Now())

	wxUsers,err:= FetchWxUser()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range(wxUsers){
		if strings.Contains(v.Wxid,"@openim") {
			continue
		}
		if !WSIn(Users,v.Wxid) {
			nu:=models.User{
				Wxid:v.Wxid,
				Account:v.Account,
				Name:v.Name,
			}
			initializers.DB.Create(&nu)
			Ticket[v.Wxid]=[]int64{5000,5000}
			fmt.Printf("\nUsers : %v--------%v\n", v, time.Now())
		}
	}
}

type Slice []models.ChatRoomUser
func (u Slice)In(id string) bool {
	for _,i:=range(u){
		if id==i.Wxid{
			return true
		}
	}
	return false
}


func WIn(u []WxUser, id string) bool {
	for _,i:=range(u){
		if id==i.Wxid{
			return true
		}
	}
	return false
}

func WSIn(u []models.User, id string) bool {
	for _,i:=range(u){
		if id==i.Wxid{
			return true
		}
	}
	return false
}


func CRURegister(roomId string){
	users:= Slice{}
	initializers.DB.Find(&users)
	wxUsers,err := FetchWxUser()
	if err != nil {
		panic(err)
	}
	if len(roomId) > 5 {
		mbs,err:=FetchRoomUser(roomId)
		if err!=nil{
			return
		}
		mba := strings.Split(mbs, "^G")
		for _,mb:=range(mba){
			if !users.In( mb) {
				cru:=models.ChatRoomUser{
					Wxid:mb,
					Room:roomId,
					IsFriend:WIn(wxUsers,mb),
				}
				initializers.DB.Create(&cru)
				fmt.Printf("\nCRURegister : %s--------%v\n", mb, time.Now())
			}
		} 
		
		return
	}
	for _, v := range(wxUsers){
		if !strings.Contains(v.Wxid,"@chatroom") {
			continue
		}
		mbs,err:=FetchRoomUser(v.Wxid)
		if err!=nil{
			continue
		}
		mba := strings.Split(mbs, "^G")
		for _,mb:=range(mba){
			if !users.In(mb) {
				cru:=models.ChatRoomUser{
					Wxid:mb,
					Room:roomId,
					IsFriend:WIn(wxUsers,mb),
				}
				initializers.DB.Create(&cru)
				fmt.Printf("\nCRURegister : %s--------%v\n", mb, time.Now())
			}
		} 
	}
}


type WxUser struct{
	Name 			string	 `json:"nickname"`
	Wxid			string 	 `json:"wxid"`
	Account			string 	 `json:"customAccount"`
}

func FetchWxUser() ([]WxUser, error) {
	client := &http.Client{Timeout: 5*time.Second}	
	requestBody := new(bytes.Buffer)
	res, err := client.Post("http://localhost:1909/api/getContactList", "application/json", requestBody)
	if err!=nil {
		fmt.Println(err)
		return nil, err
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	type Body struct{
		Code 	int			`json:"code"`
		Data 	[]WxUser 	`json:"data"`
		Msg 	string		`json:"msg"`
	}
	var result Body
	err=json.Unmarshal(body, &result)
	return result.Data, err
}

func FetchRoomUser(roomId string) (string, error) {
	client := &http.Client{Timeout: 5*time.Second}
	requestBody := bytes.NewBuffer([]byte("{\"chatRoomId\":\""+roomId+"\"}"))
	res, err := client.Post("http://localhost:1909/api/getMemberFromChatRoom", "application/json", requestBody)
	if err!=nil {
		fmt.Println(err)
		return "", err
	}
	body, _ := io.ReadAll(res.Body)
	res.Body.Close()
	
	type Menber struct{
		Admin			string	`json:"admin"`
		AdminNickname	string	`json:"adminNickname"`
		ChatRoomId		string	`json:"chatRoomId"`
		MemberNickname	string	`json:"memberNickname"`
		Members			string	`json:"members"`
	}
	type Body struct{
		Code 	int			`json:"code"`
		Data 	Menber	 	`json:"data"`
		Msg 	string		`json:"msg"`
	}
	var result Body
	err=json.Unmarshal(body, &result)
	return result.Data.Members, err
}