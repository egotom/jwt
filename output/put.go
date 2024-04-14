package output

import (
	"io"
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"errors"
	"net/http"
	"encoding/json"
	"jwt/models"
	"jwt/initializers"
	
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


