package userhandler

import (
	"pcc_card/global"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (u *User_handler_impl) GetAllOnePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req MailGetReq

		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		id := c.GetInt("id")
		// 传入 c.Request.Context()
		res, err := u.s.GetAllOnePage(c.Request.Context(), id, req.Page)
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
		// 传入 c.Request.Context()
		res, err := u.s.GetMailStatus(c.Request.Context(), id)
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
		// 传入 c.Request.Context()
		err := u.s.SendMail(c.Request.Context(), id, req.Body, req.AcceptId)
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
			// 传入 c.Request.Context()
			err := u.s.ChangeMailStatus(c.Request.Context(), id, MailId, req.Status)
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
		err := u.s.DeleteMailByMailId(c.Request.Context(), req.MailId, id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}

func (u *User_handler_impl) DeleteMailAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		// 传入 c.Request.Context()
		err := u.s.DeleteMailAll(c.Request.Context(), id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}
