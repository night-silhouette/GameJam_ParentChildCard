package userhandler

import (
	"pcc_card/global"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

func (u *User_handler_impl) GetAllOnePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MailGetReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		id := c.GetInt("id")
		res, err := u.s.GetAllOnePage(id, req.Page)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, res)
	}
}

func (u *User_handler_impl) GetMailStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		res, err := u.s.GetMailStatus(id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, res)
	}
}

func (u *User_handler_impl) SendMail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MailPostReq
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		id := c.GetInt("id")
		err := u.s.SendMail(id, req.Body, req.AcceptId)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "发送成功")
	}
}

func (u *User_handler_impl) ChangeMailStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		var req MailStatusPostReq
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		for MailId := range req.MailId {
			err := u.s.ChangeMailStatus(id, MailId, req.Status)
			if err != global.ResponseSuccess {
				response.Fail(c, err)
				return
			}
		}
		response.Success(c, "ok")

	}
}

func (u *User_handler_impl) DeleteMailByMailId() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		var req MailDeleteReq
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		for MailId := range req.MailId {
			err := u.s.DeleteMailByMailId(id, MailId)
			if err != global.ResponseSuccess {
				response.Fail(c, err)
				return
			}
		}
		response.Success(c, "clz我猜你做到这要4月10号啦")
	}
}

func (u *User_handler_impl) DeleteMailAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		err := u.s.DeleteMailAll(id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}
