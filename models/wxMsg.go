package models

type WxMsg struct{
	MsgId			int64	`json:"msgId"`
    MsgSequence		int64	`json:"msgSequence"`
    Type			int64	`json:"type"`
    Content			string 	`json:"content"`
    Signature		string	`json:"signature"`
    FromUser		string	`json:"fromUser"`
    ToUser			string	`json:"toUser"`
    DisplayFullContent	string	`json:"displayFullContent"`
    CreateTime		int64 	`json:"createTime"`
}