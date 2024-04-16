package controllers

import (
	"fmt"
	"regexp"
	// "time"
	"strings"
	"net/http"	
	// "jwt/openai"
	"jwt/google"
	"jwt/models"
	"jwt/output"
	"encoding/xml"
	"jwt/scheduler"	
	"jwt/initializers"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)

func WxMsg(c *gin.Context){
	// str, _ := c.GetRawData()
	// go output.Log2file(string(str))
	var body models.WxMsg
	if err:=c.Bind(&body); err !=nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"参数错误！"})
		fmt.Println(err)
		return 
	}
	
	// if (body.Type == 10000){
	//	群增加新成员,红包、系统消息。风险，会触发微信服务器调用 
	// 	go scheduler.CRURegister(body.FromUser)
	// 	return
	// }
	// 好友确认
	if (body.Type == 37 ){
		fmt.Println("\n\ncontent: ",body.Content)
		fmt.Println("from: ",body.FromUser)
		fmt.Println("to: ",body.ToUser)
		fmt.Println("type: ",body.Type)
		fmt.Println("\n\n")
	}
	// 共享实时位置、文件、转账、链接
	// if (body.Type == 49 ){
	// 	fmt.Println("\n\ncontent: ",body.Content)
	// 	fmt.Println("from: ",body.FromUser)
	// 	fmt.Println("to: ",body.ToUser)
	// 	fmt.Println("type: ",body.Type)
	// 	fmt.Println("\n\n")
	// }
	// if (body.Type != 49 && body.Type != 47&& body.Type != 43&& body.Type != 10002){
	if (body.Type == 49 ){
		// fmt.Println("\n\ncontent: ",body.Content)
		// fmt.Println("from: ",body.FromUser)
		// fmt.Println("to: ",body.ToUser)
		// fmt.Println("type: ",body.Type)
		// fmt.Println("\n\n")

		re := regexp.MustCompile(`^\w+:`)
		content:=re.ReplaceAllString(body.Content, "")
		var msg49 models.QuoteMsg 
		if err := xml.Unmarshal([]byte(content), &msg49); err != nil {
			fmt.Println(err)
			return
		}
		if msg49.Appmsg.Refermsg.Type == 1 {
			body.Content=fmt.Sprintf("\"%s\", %s", msg49.Appmsg.Refermsg.Content, msg49.Appmsg.Title)
			body.Type=1
		}else{
			go output.Log2file(content)
		}
	}
	
	to:="filehelper"
	if body.ToUser != "filehelper"{
		to=body.FromUser
	}

	if to=="weixin" ||  body.Type!=1 {
		return
	}
	ticket, exist := scheduler.Ticket[to] 
	if !exist {
		scheduler.WxRegister()
		if len(body.Content)>0 && !strings.Contains(body.FromUser,"@"){
			output.Reply(to, scheduler.RandChoice(initializers.Config.Greets))
		}
		// time.Sleep(5*time.Second)
		return
	}

	if (body.ToUser ==scheduler.WxMe.Wxid || body.ToUser == "filehelper") && !strings.Contains(body.FromUser,"@chatroom") {
		if CMDS(body.Content, to, ticket[0]){
			return
		}
		// go openai.ChatGPT(to, body.Content)
		go google.ChatGemini(to, body.Content)
		return
	}

	if strings.Contains(body.FromUser,"@chatroom") && strings.Contains(body.Content, "@"+scheduler.WxMe.Name){
		re := regexp.MustCompile(`^\w+:\s@`+scheduler.WxMe.Name+".?")
		content := re.ReplaceAllString(body.Content, "")
		content = strings.Replace(content, "@"+scheduler.WxMe.Name, "",-1)
		if CMDS(content, to, ticket[0]){
			return
		}
		// go openai.ChatGPT(to, body.Content)
		go google.ChatGemini(to, content)
		return
	}
}

func CMDS(Content,to string, ticket int64) bool{
	if Content== "推荐码"{
		id:= uuid.Must(uuid.NewRandom()).String()		//uuid.New()
		err:= output.Promote(to, id)
		if err == nil {
			output.Reply(to, fmt.Sprintf("%s推荐码:%s",scheduler.WxMe.Name, id))
			return true
		}
		output.Reply(to, fmt.Sprintf("%s\n联系微信: mtldswz_03 ", err))
		output.ReplyImg(to, "E:/USB/mtl.jpg")
		return true
	}

	if strings.HasPrefix(Content, scheduler.WxMe.Name+"推荐码") {
		id:= strings.Replace(Content,scheduler.WxMe.Name+"推荐码:", "", 1)
		fmt.Println(Content,"--------------------------------------------")
		var pm = models.Promotion{ID:id}
		result := initializers.DB.Where("publisher != ?",to).First(&pm, id)
		if result.RowsAffected == 0{
			output.Reply(to, "该推荐码无效，或自己生成？请联系推荐者重新生成。")
			Help(to)
			return true
		}
		pmu := models.Promotion{}
		result = initializers.DB.Where("consumer like ?","%"+to+"%").First(&pmu)
		if result.RowsAffected > 0 {
			output.Reply(to, "你已经使用过推荐码。")
			Help(to)
			return true
		}
		pm.Consumer=fmt.Sprintf("%s,%s", pm.Consumer, to)
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
				output.Reply(to, "推荐码使用成功！充值可获9折优惠。请联系微信:  mtldswz_03")
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
		output.Reply(to, quota)
		return true
	} 

	if ticket < 0 {
		output.Reply(to, "今日的免费额度已用完。推荐新用户获得每天300 token免费额度奖励。\n充值30元，每天使用额度增加5000token。联系微信: mtldswz_03。")
		Help(to)
		if to == "wxid_ievhmifil5vf22"{
			scheduler.Ticket[to] = []int64{5000, 5000}
		}
		return true
	}
	return false
}

func Help(to string){
	msg:=fmt.Sprintf("发送\"?\"或\"帮助\":  显示本说明。\n\n发送\"推荐码\":  生成推荐码。被使用后, 你每天的使用额度可获得300 token 奖励。被推荐人充值可获得33%%的提成。\n\n使用推荐码: 请直接把别人分享给你的推荐码转发给 @%s，充值获得9折优惠。\n\n充值：充值30元，每天使用额度增加5000token。\n\n发送\"额度\"或\"余额\": 查询每天使用额度和当天剩余额度。\n\n充值或任何问题: 请联系微信: mtldswz_03", scheduler.WxMe.Name)
	output.Reply(to, msg)
	output.ReplyImg(to, "E:/USB/mtl.jpg")
}