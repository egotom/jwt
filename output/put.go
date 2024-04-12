package output

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"errors"
	"strings"
	"net/http"
	"encoding/json"
	"jwt/models"
	"jwt/scheduler"
	"jwt/initializers"
	"github.com/google/uuid"
)

func Log2DB(to, msg, role string, vendor uint8){
	// v := models.Content{
	// 	Vendor:vendor,
	// 	Wxid: to,
	// 	Msg: fmt.Sprintf("{\"role\":\"%s\",\"content\":\"%s\"}", role, msg),
	// }
	// initializers.DB.Create(&v)
	v := models.Context{
		Vendor:vendor,
		Wxid: to,
		Role: role,
		Msg: msg,
	}
	initializers.DB.Create(&v)
}


func Log2file(text string){
	f, err := os.OpenFile("msgLog.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	
	log.SetOutput(f)
	log.Println(text,"\n")
}

func Promote(user string, id string) error {
	var count int64
	initializers.DB.Model(&models.Promotion{}).Where("publisher = ?", user).Where("created_at > ?", "curdate()").Count(&count)
	if count>10{
		return errors.New("今天已经生成了10个推荐码，额度已用完。")
	}
	v:=models.Promotion{
		ID:id,
		Publisher: user,
	}
	result := initializers.DB.Create(&v)
	if result.Error == nil {
		return nil
	}
	return errors.New("生成推荐码失败！")
}

func Reply(to, msg string){
	type Body struct{
		Wxid 		string	`json:"wxid"`
		Msg 		string	`json:"msg"`
	}
	payload := Body{
		Wxid: to,
		Msg:msg,
	}
	
	requestBody,err := json.Marshal(&payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("requestBody:", string(requestBody))
	client := &http.Client{Timeout: 3*time.Second}
	req, err := http.NewRequest("POST", "http://localhost:1909/api/sendTextMsg", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func ReplyImg(to, imagePath string){
	type Body struct{
		Wxid 			string	`json:"wxid"`
		ImagePath 		string	`json:"imagePath"`
	}
	payload := Body{
		Wxid: to,
		ImagePath:imagePath,
	}
	
	requestBody,err := json.Marshal(&payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("requestBody:", string(requestBody))
	client := &http.Client{Timeout: 3*time.Second}
	req, err := http.NewRequest("POST", "http://localhost:1909/api/sendImagesMsg", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func Help(to string){
	msg:=fmt.Sprintf("发送\"?\"或\"帮助\":  显示本说明。\n\n发送\"推荐码\":  生成推荐码。被使用后, 你每天的使用额度可获得300 token 奖励。被推荐人充值可获得33%%的提成。\n\n使用推荐码: 请直接把别人分享给你的推荐码转发给 @%s，充值获得9折优惠。\n\n发送\"额度\"或\"余额\": 查询每天使用额度和当天剩余额度。\n\n充值或任何问题: 请联系微信: mtldswz_03", scheduler.WxMe.Name)
	Reply(to, msg)
	ReplyImg(to, "E:/USB/mtl.jpg")
}

func CMDS(Content,to string, ticket int64) bool{
	if Content== "推荐码"{
		id:= uuid.Must(uuid.NewRandom()).String()		//uuid.New()
		err:= Promote(to, id)
		if err == nil {
			Reply(to, fmt.Sprintf("%s推荐码:%s",scheduler.WxMe.Name, id))
			return true
		}
		Reply(to, fmt.Sprintf("%s\n联系微信: mtldswz_03 ", err))
		ReplyImg(to, "E:/USB/mtl.jpg")
		return true
	}

	if strings.HasPrefix(Content, scheduler.WxMe.Name+"推荐码") {
		id:= strings.Replace(Content,scheduler.WxMe.Name+"推荐码:", "", 1)
		fmt.Println(id)
		var pm = models.Promotion{ID:id}
		result := initializers.DB.Where("publisher != ?",to).First(&pm)
		if result.RowsAffected == 0{
			Reply(to, "该推荐码无效，请联系推荐者重新生成。")
			Help(to)
			return true
		}
		if strings.Contains(pm.Consumer, to) {
			Reply(to, "你已经使用过推荐码。")
			Help(to)
			return true
		}
		pm.Consumer=pm.Consumer+fmt.Sprintf("%s,",to)
		result = initializers.DB.Save(&pm)
		if result.RowsAffected>0{
			var pub = models.User{}
			result = initializers.DB.Where("wxid = ?",  pm.Publisher).First(&pub)
			if result.RowsAffected>0{
				pub.Ticket = pub.Ticket+300
				result = initializers.DB.Save(&pub)
				scheduler.Ticket[pm.Publisher]=[]int64{scheduler.Ticket[pm.Publisher][0]+300,scheduler.Ticket[pm.Publisher][1]+300}
				// congrats:=fmt.Sprintf("你的推荐码:%s，已经被使用。恭喜获得每天300 token 使用额度奖励。", id)
				// Reply(pm.Publisher, congrats)
				Reply(to, "推荐码使用成功！充值可获9折优惠。请联系微信:  mtldswz_03")
			}
		}
		return true
	}

	if Content== "帮助" || Content== "？"|| Content== "?"{
		Help(to)
		return true
	}

	if Content== "额度" || Content== "余额"{
		quota:=fmt.Sprintf("每天使用额度: %d Token\n今天剩余额度: %d Token", scheduler.Ticket[to][1], scheduler.Ticket[to][0])
		Reply(to,quota)
		return true
	} 

	if ticket < 0 {
		Reply(to, "今日的免费额度已用完。推荐新用户获得每天300 token免费额度奖励。\n充值30元，每天限额加倍。联系微信: mtldswz_03。")
		Help(to)
		if to == "wxid_ievhmifil5vf22"{
			scheduler.Ticket[to] = []int64{5000, 5000}
		}
		return true
	}
	return false
}