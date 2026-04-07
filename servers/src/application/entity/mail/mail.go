package mail

import (
	"pcc_card/global"
	"time"
)

var CategoryMap = map[string]int{
	"UserMail":           0,
	"FriendshipsRequest": 1,
}

type Mail struct {
	MailId   int       `json:"mail_id"`
	AcceptId int       `json:"accept_id"`
	SendId   int       `json:"send_id"`
	Body     string    `json:"body"`
	Category string    `json:"category"`
	Status   int       `json:"status"`
	CreateAt time.Time `json:"create_at"`
}

func NewMail(AcceptId int, SendId int, Body string, Category string) (*Mail, global.ResponseStatusCode) {
	if _, ok := CategoryMap[Category]; !ok {
		return nil, global.ResponseInvalidReqParams
	}

	mail := Mail{}
	mail.Status = 0
	mail.AcceptId = AcceptId
	mail.SendId = SendId
	mail.Body = Body
	mail.Category = Category
	return &mail, global.ResponseSuccess
}

type Filter struct {
	Id       string
	AcceptId string
	SendId   string
	Category string
	Status   string
}
