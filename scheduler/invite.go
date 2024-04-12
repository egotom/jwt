package scheduler

import (
	"fmt"
	"time"
	"bytes"
	"net/http"
	"jwt/models"
	"jwt/initializers"
)

func Invite(n int){
	rids:=[]string{"35031914979@chatroom"}
	Users := []models.User{}
	initializers.DB.Where("is_active=?","false").Not("wxid like ?","%@chatroom").Order("updated_at desc").Find(&Users)
	if len(Users)<1{
		fmt.Println("Invite(), 加载用户失败！")
	}
	for i,u := range(Users){
		fmt.Println(i,"==============================")
		if i >= n{
			break
		}
		if err := CRUInvite(rids,u.Wxid); err==nil{
			initializers.DB.Model(models.User{}).Where("id=?",u.ID).Updates(models.User{IsActive:true, IsVerified:true})
		}
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