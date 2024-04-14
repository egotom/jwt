package output

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"strings"
	"net/url"
	"net/http"
	"jwt/models"
	"encoding/json"	
	"jwt/initializers"
)

type Task struct{
	Style 	string 
	Formart	string
	Words	[]string
}

func Honey(){
	tasks := []Task{
		// Task{
		// 	Style:"书摘-SAT官方推荐高中必读",
		// 	Formart:"写一则美国SAT官方推荐高中必读书《%s》的原文书摘，用以下格式回答：{书摘原文} -- 摘自：{作者}，{书名}。不需要评论".
		// 	Words:[]string{"玻璃动物园","道林格雷的画像","草叶集","尤多拉韦尔蒂全集","欢乐之家","紫色","第五号屠宰场","老实人","哈克贝利费恩历险记","父与子","战争与和平","瓦尔登湖","名利场","格列佛游记","金银岛","愤怒的葡萄","俄狄浦斯王","安提戈涅","伊凡杰尼索维奇的一天","典仪","皮格马利翁","罗密欧与朱丽叶","仲夏夜之梦","麦克白","哈姆雷特","麦田里的守望者","就说是睡着了","大鼻子情圣","西线无战事","拍卖第49批","追忆似水年华","爱伦坡短片小说选","钟形罩","日瓦戈医生","动物庄园","长夜漫漫路迢迢","好人难寻","宠儿","萨勒姆的女巫","白鲸记","文书巴特尔比","百年孤独","魔山","野性的呼唤","巴比特","杀死一只知更鸟","女战士","变形记","一个青年艺术家的画像","螺丝在拧紧","一个贵妇人的画像","玩偶之家","美丽新世界","凝望上帝","巴黎圣母院","奥德赛","伊利亚特","永别了，武器","第22条军规","红字","苔丝","蝇王","浮士德","好兵","包法利夫人","了不起的盖茨比","汤姆琼斯","喧哗与骚动","我弥留之际","爱默生散文精选","隐形人","佛洛斯河上的磨坊","三个火枪手","美国的悲剧","一个美国黑奴的自传","罪与罚","双城记","鲁宾逊漂流记","唐吉可德","神曲地狱篇","红色英勇勋章","最后的莫希干人","黑暗之心","觉醒","樱桃园","坎特伯雷故事集","大主教之死","局外人","呼啸山庄","简爱","奥吉玛琪历险记","等待戈多","高山上的呼喊","傲慢与偏见","家中丧事","瓦解","贝奥武夫"},
		// },
		// Task{
		// 	Style:"BBC推荐，影响世界的77部文学名著",
		// 	Formart:"写一则文学名著《%s》的原文书摘，用以下格式回答：{书摘原文} -- 摘自：{作者}，{书名}。不需要评论",
		// 	Words: []string {"阿伽门农","法律门前","五卷书","杀死一只知更鸟","少年维特之烦恼","狼厅","哈利波特系列","唐吉诃德","炼狱","圣经","源氏物语","金色笔记","哈姆雷特","凝望上帝","奥梅罗斯","水浒传","嚎叫","伊凡杰尼索维奇的一天","乌鸦","黑暗之心","伊索寓言","白鲸","审判","包法利夫人","佩德罗巴拉莫","玛特辽娜的家","欧墨尼得","罪与罚","抒情诗集","使女的故事","奥德赛","穿破裤子的慈善家","革命之路","傲慢与偏见","拉德茨基进行曲","俄瑞斯忒亚","米德尔马契","罗密欧与朱丽叶","德拉库拉","神曲","渴望真实的狮子","看不见的人","宠儿","大师和玛格丽特","鲁宾逊漂流记","向日葵箴言","一只右手","纽约时报","弗兰肯斯坦","追忆似水年华","李尔王","变形记","天方夜谭","动物庄园","吉尔伽美什史诗","坎特伯雷故事集","简爱","老实人","战争与和平","简布罗迪小姐的青春","时代周刊","伊利亚特","一九八四","布拉兹库巴斯的死后回忆","地狱","圣诞颂歌","小王子","瓦解","达洛维夫人","奠酒人","荷马史诗","屠场","北京偶感","悲惨世界","呼啸山庄","午夜的孩子","安提戈涅","天堂","我们街区的孩子们","藻海无边","百年孤独","汤姆叔叔的小屋","一千零一夜","白牙","大师与玛格丽特","爱丽丝梦游仙境","美狄亚","黑暗的心","列王荣耀记","美国","奥兰多","尤利西斯","黑暗的左手","到灯塔去","三国演义"},
		// },
	}
	for _,t := range(tasks){
		style=t.Style
		for _,w := range(t.Words){
			prompt:=fmt.Sprintf(t.Formart, w)
			ChatGemini(prompt)
			time.Sleep(60*time.Second)
		}
	}
}

func ChatGemini(msg string){
	msgs :=[]models.Contents{}
	msgs = append(msgs, models.Contents{Role: "user", Parts:[]models.Parts{{Text:msg}}})
	
	rb:=models.GUserCtx{Contents:msgs}
	m, err := json.Marshal(&rb)
	if err!=nil{
		fmt.Println(err)
		ChatGPT(msg)
		return
	}
	ugi:=fmt.Sprint("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=", initializers.Config.APIKey["Gemini"])
	// ugi:=fmt.Sprint("https://deluxe-trifle-53ce77.netlify.app/v1/models/gemini-pro:generateContent?key=%s", initializers.Config.APIKey["Gemini"])
	r, err := http.NewRequest("POST",ugi, bytes.NewBuffer(m))
	if err != nil {
		fmt.Println(err)
		ChatGPT(msg)
		return
	}
	r.Header.Add("Content-Type", "application/json")
	// client:= &http.Client{Timeout: 15*time.Second}
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
		ChatGPT(msg)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err!=nil{
		ChatGPT(msg)
		return
	}
	
	answer := models.AnswerG{}
	err =json.Unmarshal(body, &answer)
	if err==nil && len(answer.Candidates)>0 && len(answer.Candidates[0].Content.Parts) >0{
		res:=strings.Replace(answer.Candidates[0].Content.Parts[0].Text,"*","",-1)
		// Log2file(res)
		logHoney(msg, res)

	}else{
		ChatGPT(msg)
	}
}

func ChatGPT(msg string){
	msgs :=[]models.OpenaiMsg{}
	msgs = append(msgs, models.OpenaiMsg{ Role: "user", Content: msg})

	rb:=models.OAChatGPTRequestBody{
		Model:"gpt-3.5-turbo-0125",
		Messages:msgs,
	}
	m, err := json.Marshal(&rb)
	if err!=nil{
		fmt.Println(err)
	}
	r, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(m))
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
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err!=nil{
		// panic(err)
		return
	}

	answer := models.Answer{}
	err =json.Unmarshal(body, &answer)
	if err==nil && len(answer.Choices)>0 {
		text := answer.Choices[0].Message.Content
		// Log2file(text)
		logHoney(msg,text)
	}
}

var honey models.Honey
var style string
func logHoney(p, c string){
	honey = models.Honey{Style:style, Prompt:p, Commit:c}
	initializers.DB.Create(&honey)
}
