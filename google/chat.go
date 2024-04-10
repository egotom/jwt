package google

import (
	"io"
	"os"
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/url"
	"net/http"
	"jwt/models"
	"jwt/output"
	"jwt/openai"
	"jwt/scheduler"
	"encoding/json"
)

func ChatGemini(to, msg string){
	go output.Log2DB(to, msg, "user", 2)
	Now:=time.Now()
	msgs :=[]models.Contents{}
	if Now.Sub(scheduler.Ticker[to])<30*time.Minute {
		msgs = scheduler.GUserCtx[to]
		if len(msgs) > 12{
			msgs=msgs[2:]
		}
	}
	msgs = append(msgs, models.Contents{Role: "user", Parts:[]models.Parts{{Text:msg}}})

	omsg :=[]models.OpenaiMsg{}
	if Now.Sub(scheduler.Ticker[to])<30*time.Minute {
		omsg := scheduler.UserCtx[to]
		if len(omsg) > 6{
			omsg=omsg[1:]
		}
	}
	scheduler.UserCtx[to]= append(omsg, models.OpenaiMsg{ Role: "user", Content: msg})
	scheduler.Ticker[to]=Now 

	rb:=models.GUserCtx{Contents:msgs}
	m, err := json.Marshal(&rb)
	if err!=nil{
		fmt.Println(err)
	}
	uriGemini:=fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", os.Getenv("GEMINI_API_KEY"))
	// uriGemini:=fmt.Sprintf("https://deluxe-trifle-53ce77.netlify.app/v1/models/gemini-pro:generateContent?key=%s",os.Getenv("GEMINI_API_KEY"))
	r, err := http.NewRequest("POST",uriGemini, bytes.NewBuffer(m))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	// client:= &http.Client{Timeout: 15*time.Second}
	client:= &http.Client{
		Timeout: 15*time.Second,
		Transport: &http.Transport {
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host: os.Getenv("HTTP_PROXY"),
			}),
		},
	}
	res, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		// output.Reply(to, scheduler.WxMe.Name+"正在忙, 请稍后再试!")
		openai.ChatGPT(to, msg)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err!=nil{
		// panic(err)
		// output.Reply(to, "服务器出错了，请稍后再试！")
		openai.ChatGPT(to, msg)
		return
	}
	
	answer := models.AnswerG{}
	err =json.Unmarshal(body, &answer)
	if err==nil && len(answer.Candidates)>0 && len(answer.Candidates[0].Content.Parts) >0{
		res:=strings.Replace(answer.Candidates[0].Content.Parts[0].Text,"*","",-1)
		// res:=answer.Candidates[0].Content.Parts[0].Text
		// fmt.Println(len(res))
		// fmt.Println(res)
		output.Reply(to, res)
		output.Log2DB(to, res, "model", 2)
		scheduler.Ticket[to][0]=scheduler.Ticket[to][0]-int64(len(res)/5)
		msgs = append(msgs, models.Contents{Role: "model", Parts:[]models.Parts{{Text:res}}})	
		scheduler.GUserCtx[to]=msgs	
	}else{
		msgs = append(msgs, models.Contents{Role: "model", Parts:[]models.Parts{{Text:""}}})
		fmt.Println(string(body))
		// output.Reply(to, scheduler.WxMe.Name+"没有回应。发一个句号提醒他一下。")
		openai.ChatGPT(to, msg)
	}
	// scheduler.GUserCtx[to]=msgs
	// fmt.Println(string(body))
	// fmt.Println(msgs)
	// fmt.Println(answer)
}