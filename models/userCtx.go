package models

type Parts struct{
	Text	string	`json:"text"`
}
type Contents struct{
	Role	string	`json:"role"`
	Parts	[]Parts	`json:"parts"`
}
type GUserCtx struct{
	Contents	[]Contents	`json:"contents"`
}

type Candidates struct{
	Content	Contents	`json:"content"`
}

type AnswerG struct{
	Candidates	[]Candidates	`json:"candidates"`
}
// {
// 	"contents": [{
// 		"role":"user",
// 		"parts":[{"text": "Write the first line of a story about a magic backpack."}]
// 	},
// 	{
// 		"role": "model",
// 		"parts":[{"text": "In the bustling city of Meadow brook, lived a young girl named Sophie. She was a bright and curious soul with an imaginative mind."}]
// 	},
// 	{
// 		"role": "user",
// 		"parts":[{"text": "Can you set it in a quiet village in 1600s France?"}]
// 	}]
// }

// {
// 	"candidates": [
// 	  {
// 		"content": {
// 		  "parts": [
// 			{
// 			  "text": "我是 Gemini，谷歌开发的多模态 AI 语言模型。"
// 			}
// 		  ],
// 		  "role": "model"
// 		},
// 		"finishReason": "STOP",
// 		"index": 0,
// 		"safetyRatings": [
// 		  {
// 			"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT",
// 			"probability": "NEGLIGIBLE"
// 		  },
// 		  {
// 			"category": "HARM_CATEGORY_HATE_SPEECH",
// 			"probability": "NEGLIGIBLE"
// 		  },
// 		  {
// 			"category": "HARM_CATEGORY_HARASSMENT",
// 			"probability": "NEGLIGIBLE"
// 		  },
// 		  {
// 			"category": "HARM_CATEGORY_DANGEROUS_CONTENT",
// 			"probability": "NEGLIGIBLE"
// 		  }
// 		]
// 	  }
// 	],
// 	"promptFeedback": {
// 	  "safetyRatings": [
// 		{
// 		  "category": "HARM_CATEGORY_SEXUALLY_EXPLICIT",
// 		  "probability": "NEGLIGIBLE"
// 		},
// 		{
// 		  "category": "HARM_CATEGORY_HATE_SPEECH",
// 		  "probability": "NEGLIGIBLE"
// 		},
// 		{
// 		  "category": "HARM_CATEGORY_HARASSMENT",
// 		  "probability": "NEGLIGIBLE"
// 		},
// 		{
// 		  "category": "HARM_CATEGORY_DANGEROUS_CONTENT",
// 		  "probability": "NEGLIGIBLE"
// 		}
// 	  ]
// 	}
//   }