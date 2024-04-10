package scheduler

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"net/http"
	"jwt/models"
	"encoding/json"	
	"jwt/initializers"
	"github.com/go-co-op/gocron/v2"
)

func Scheduler() {
	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Println("error", err)
		return
	}
	// j, err := s.NewJob(gocron.DurationJob(55*time.Second), gocron.NewTask(WxRegister))
	j, err := s.NewJob(
		gocron.DailyJob(1,
			gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0)),
		),
		gocron.NewTask(
			func(){
				fmt.Println("gocron.DailyJob: ",time.Now())
				var users []models.User
				initializers.DB.Select("wxid","ticket").Find(&users)
				for _,u := range(users){
					Ticket[u.Wxid]=[]int64{u.Ticket,u.Ticket}
				}
			},
		),
	)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println("job has a unique id: ",j.ID())
	s.Start()
	// when you're done, shut it down
	// err = s.Shutdown()
	// if err != nil {
	// 	fmt.Println("error", err)
	// }
	HookSyncMsg()
	UserInfo()
	WxRegister()
	LoadCache()
}

func HookSyncMsg(){
	type Body struct{
		Port 		string	`json:"port"`
		Ip 			string	`json:"ip"`
		Url 		string	`json:"url"`
		Timeout 	string	`json:"timeout"`
		EnableHttp  string	`json:"enableHttp"`
	}
	payload := Body{
		Port: "8080",
		Ip:"127.0.0.1",
		Url:"http://localhost:8080/wxmsg",
		Timeout:"50000",
		EnableHttp:"1",
	}
	requestBody, err := json.Marshal(&payload)
	if err!=nil{
		panic(err)
	}
	// fmt.Println("HookSyncMsg requestBody:", string(requestBody))
	client := &http.Client{Timeout: 3*time.Second}	
	res, err := client.Post("http://localhost:1909/api/hookSyncMsg", "application/json", bytes.NewBuffer(requestBody))
	if err!=nil{
		panic(err)
	}
	body, _ := io.ReadAll(res.Body)
	fmt.Println("HookSyncMsg response Body:", string(body))
}

var(
	WxMe models.WxUser
	UserCtx = map[string]([]models.OpenaiMsg){}
	GUserCtx = map[string]([]models.Contents){}
	Ticket 	= map[string][]int64{}
	Ticker 	= map[string]time.Time{}
)

func UserInfo(){
	client := &http.Client{Timeout: 3*time.Second}	
	res, err := client.Post("http://localhost:1909/api/userInfo", "application/json", bytes.NewBuffer([]byte("{}")))
	if err!=nil{
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	if err!=nil{
		panic(err)
	}
	type Body struct{
		Code	uint 			`json:"code"`
		Data 	models.WxUser 	`json:"data"`
	}
	var result Body
	json.Unmarshal(body, &result)
	WxMe=result.Data
	fmt.Printf("%v\tAccount: %v\tWxID: %v\n", WxMe.Name, WxMe.Account, WxMe.Wxid)
}

func LoadCache(){
	var users []models.User
	initializers.DB.Select("wxid","ticket").Find(&users)
	for _,u := range(users){
		Ticket[u.Wxid]=[]int64{u.Ticket, u.Ticket}
		Ticker[u.Wxid]=time.Now()
	}

	rows,_ := initializers.DB.Raw("select t.wxid, t.msg from (SELECT *, row_number() over (partition by wxid order by created_at desc) as seqnum FROM contexts where created_at >= DATE_SUB(NOW(), INTERVAL 30 MINUTE) AND role='user' ) t where seqnum < 7 order by created_at; ").Rows()
	defer rows.Close()
	var Wxid	string
	var Msg 	string
	for rows.Next() {
		rows.Scan(&Wxid, &Msg)
		m :=models.OpenaiMsg{Role:"user", Content:Msg}
		ctx, wx := UserCtx[Wxid]
		if wx {
			ctx = append(ctx, m)
			UserCtx[Wxid] = ctx
		}else{
			UserCtx[Wxid] = []models.OpenaiMsg{m}
		}

		u :=[]models.Contents{
			{Role: "user", Parts:[]models.Parts{{Text:Msg}}},
			{Role: "model", Parts:[]models.Parts{{Text:""}}},
		}
		gctx, wx := GUserCtx[Wxid]
		if wx {
			gctx = append(gctx, u...)
			GUserCtx[Wxid] = gctx
		}else{
			GUserCtx[Wxid] = u
		}
	}

	// fmt.Printf("\nTicket[\"wxid_ievhmifil5vf22\"]: %v\n", Ticket["wxid_ievhmifil5vf22"])
	// fmt.Printf("\nUserCtx: %v\n", UserCtx)
	// fmt.Printf("\nGUserCtx: %v\n", GUserCtx)
}