package userhandler

type MailGetReq struct {
	Page int `form:"page" json:"page" binding:"required"`
}
type MailPostReq struct {
	AcceptId int    `form:"accept_id" json:"accept_id" binding:"required"`
	Body     string `form:"body" json:"body" binding:"required"`
}

type MailStatusPostReq struct {
	MailId []int `form:"mail_id" json:"mail_id" binding:"required"`
	Status int   `form:"status" json:"status" binding:"required"`
}

type MailDeleteReq struct {
	MailId []int `form:"mail_id" json:"mail_id" binding:"required"`
}
