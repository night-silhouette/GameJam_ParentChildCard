package userhandler

type MailGetReq struct {
	Page int `form:"page" json:"page" binding:"required"`
}
type MailPostReq struct {
	AcceptId int    `form:"accept_id" json:"accept_id" binding:"required"`
	Body     string `form:"body" json:"body" binding:"required,max=500"`
	Category string `form:"category" json:"category" binding:"required"`
}

type MailStatusPostReq struct {
	MailIdList []int `form:"mail_id_list" json:"mail_id" binding:"required"`
	Status     int   `form:"status" json:"status" binding:"required"`
}

type MailDeleteReq struct {
	MailIdList []int `form:"mail_id_list" json:"mail_id_list" binding:"required"`
}
type MailFriendshipPostReq struct {
	MailId   int  `form:"mail_id" json:"mail_id" binding:"required"`
	IsFriend bool `form:"is_friend" json:"is_friend" binding:"required"`
}
