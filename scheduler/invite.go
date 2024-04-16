package scheduler

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/http"
	"jwt/models"
	"jwt/initializers"
)

func Invite(){
	rids:=[]string{"35031914979@chatroom"}
	Users := []models.User{}
	initializers.DB.Where("is_active = ?","false").Not("wxid like ?","%@%").Limit(5).Order("updated_at desc").Find(&Users)
	if len(Users)<1{
		fmt.Println("Invite(), 加载用户失败！")
		return
	}
	var uids string
	for _,u := range(Users){
		uids = fmt.Sprintf("%s,%s", uids,u.Wxid)
	}
	client := &http.Client{Timeout: 5*time.Second}	

	for _, r := range(rids){
		time.Sleep(6*time.Second)
		menbers, err := FetchRoomUser(r)
		if err!=nil {
			fmt.Println(err)
			continue
		}
		rmc := len(strings.Split(menbers,"^G"))
		if rmc < 40 {
			rb := bytes.NewBuffer([]byte("{\"chatRoomId\":\""+r+"\",\"memberIds\":\""+uids+"\"}"))
			_, err := client.Post("http://localhost:1909/api/addMemberToChatRoom", "application/json", rb)
			if err!=nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("\n Invite : %s------------%v\n", uids, time.Now())
		}else{
			rb := bytes.NewBuffer([]byte("{\"chatRoomId\":\""+r+"\",\"memberIds\":\""+uids+"\"}"))
			_, err := client.Post("http://localhost:1909/api/InviteMemberToChatRoom", "application/json", rb)
			if err!=nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("\n Invite : %s------------%v\n", uids, time.Now())
		}
		// US := models.User{}
		// initializers.DB.Model(&US).Where("wxid in ?", strings.Split(uids,",")).Update(
		// 	models.User {IsActive: true,IsVerified: true}
		// )
	}
}

func CRUInvite(rids []string, uids string) error{
	client := &http.Client{Timeout: 5*time.Second}	
	for _, r := range(rids){
		rb := bytes.NewBuffer([]byte("{\"chatRoomId\":\""+r+"\",\"memberIds\":\""+uids+"\"}"))
		_, err := client.Post("http://localhost:1909/api/InviteMemberToChatRoom", "application/json", rb)
		if err!=nil {
			fmt.Println(err)
			return err
		}
		fmt.Printf("\nCRUInvite : %s------------%v\n", uids, time.Now())
		time.Sleep(6*time.Second)
	}
	return nil
}