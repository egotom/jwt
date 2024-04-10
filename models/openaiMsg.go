package models

type OpenaiMsg struct {
	Role		string		`json:"role"`
    Content		string		`json:"content"`
}
type OAChatGPTRequestBody struct {
    Model            string  `json:"model"`
    Messages         []OpenaiMsg  `json:"messages"`
}

type TCosumer struct{
    Prompt_tokens       int64   `json:"prompt_tokens"`
    Completion_tokens   int64  `json:"completion_tokens"`
    Total_tokens        int64  `json:"total_tokens"`
}

type Choice struct {
    Index           int64       `json:"index"`
    Message         OpenaiMsg   `json:"message"` 
    Finish_reason   string      `json:"finish_reason"`
}

type Answer struct {
    Id      string      `json:"id"`
    Choices []Choice    `json:"choices"`
    Usage   TCosumer    `json:"usage"`
}