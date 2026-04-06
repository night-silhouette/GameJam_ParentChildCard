package mail

import "time"

var CategoryMap = map[string]int{
	"UserMail": 0,
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

func NewMail(AcceptId int, SendId int, Body string, Category string) *Mail {
	mail := Mail{}
	mail.Status = 0
	mail.AcceptId = AcceptId
	mail.SendId = SendId
	mail.Body = Body
	mail.Category = Category
	return &mail
}

type Filter struct {
	Id       string
	AcceptId string
	SendId   string
	Category string
	Status   string
}
