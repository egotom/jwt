package models

import "encoding/xml"

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

type Msg49 struct{
    XMLName xml.Name `xml:"msg"`
    Appmsg struct{
        Title   string  `xml:"title"`
        Type    int64   `xml:"type"`
        Refermsg struct{
            Type    int64   `xml:"type"`
            Content string  `xml:"content"`
        }   `xml:"refermsg"`
    }   `xml:"appmsg"`
}