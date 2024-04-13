package openai

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"net/url"
	"net/http"
	"encoding/json"
	"jwt/models"
	"jwt/output"
	"jwt/scheduler"
	"jwt/initializers"
)

func ChatGPT(to, msg string){
	// go output.Log2DB(to, msg,"user",1)
	msgs := scheduler.UserCtx[to]
	// if len(msgs) > 6{
	// 	msgs=msgs[1:]
	// }
	// msgs=append(msgs, models.OpenaiMsg{ Role: "user",Content: msg})
	// scheduler.UserCtx[to]=msgs
	rb:=models.OAChatGPTRequestBody{
		Model:"gpt-3.5-turbo-0125",
		Messages:msgs,
	}
	m, err := json.Marshal(&rb)
	if err!=nil{
		fmt.Println(err)
	}
	// fmt.Println(string(m))
	r, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(m))
	// OPuri := fmt.Sprintf("https://main--zingy-heliotrope-e35efb.netlify.app?key=%s", initializers.Config.APIKey["OpenAI"])
	// r, err := http.NewRequest("POST",OPuri, bytes.NewBuffer(m))
	// r, err := http.NewRequest("POST", "https://sibylla.org", bytes.NewBuffer(m))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer " + initializers.Config.APIKey["OpenAI"])
	client:= &http.Client{
		Timeout: 15*time.Second,
		Transport: &http.Transport {
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host: initializers.Config.HTTPProxy,
			}),
		},
	}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		output.Reply(to, scheduler.WxMe.Name+"正在忙, 请稍后再试!")
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err!=nil{
		// panic(err)
		output.Reply(to, "服务器出错了，请稍后再试！")
		return
	}
	// fmt.Println(string(body))
	gmsg := scheduler.GUserCtx[to]
	if len(gmsg) > 12{
		gmsg=gmsg[2:]
	}
	gmsg = append(gmsg, models.Contents{Role: "user", Parts:[]models.Parts{{Text:msg}}})

	answer := models.Answer{}
	err =json.Unmarshal(body, &answer)
	if err==nil && len(answer.Choices)>0 {
		text := answer.Choices[0].Message.Content
		// fmt.Println("Total_tokens: ",answer.Usage.Total_tokens)
		// fmt.Println(answer.Choices[0].Message.Content)
		output.Reply(to, text)
		output.Log2DB(to, text, "system",1)
		scheduler.Ticket[to][0]=scheduler.Ticket[to][0]-answer.Usage.Total_tokens
		gmsg = append(gmsg, models.Contents{Role: "model", Parts:[]models.Parts{{Text:text}}})
	}else{
		output.Reply(to, scheduler.WxMe.Name+"没有回应。发一个句号提醒他一下。")
		gmsg = append(gmsg, models.Contents{Role: "model", Parts:[]models.Parts{{Text:""}}})
	}
	scheduler.GUserCtx[to]=gmsg
}
