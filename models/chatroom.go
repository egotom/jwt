package models

type ChatRoomUser struct {
	Wxid			string	`gorm:"index;type:varchar(50);not null;unique" json:"wxid"`
	CustomAccount	string 	`gorm:"type:varchar(100);json:"customAccount"`
	Name 			string 	`gorm:"type:varchar(100);json:"name"`
	Nickname 		string 	`gorm:"type:varchar(100);json:"nickname"`
	InviteAt	    string	`gorm:"text" json:"invite_at"`
	Room    		string	`gorm:"text" json:"room"`
	IsMenber	    bool	`gorm:"default:false" json:"is_menber"`
	IsFriend	    bool	`gorm:"default:false" json:"is_friend"`
	Model
}


