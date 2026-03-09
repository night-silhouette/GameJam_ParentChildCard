package token_handler

import (
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

type Token_handler interface {
	handler.Handler
}
type Token_handler_impl struct {
	s UserService.User_service
}

func (u *Token_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(UserService.User_service)
}
func (u *Token_handler_impl) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TokenGetDto
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		_, err := u.s.Is_valid_token(req.Token)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, global.ResponseSuccess)
		return
	}
}

func (u *Token_handler_impl) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto TokenPostDto
		if err := c.ShouldBind(&dto); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		e, err := u.s.Find_user_by_name(dto.Name)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		id := e.Id
		err = u.s.Check_password(id, dto.Password)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		token := u.s.Release_token(id)
		response.Success(c, gin.H{"token": token})
	}
}
func (u *Token_handler_impl) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Fail(c, global.ResponseNotImplemented)
		return
	}
}
func (u *Token_handler_impl) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Fail(c, global.ResponseNotImplemented)
		return
	}
}
func (u *Token_handler_impl) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Fail(c, global.ResponseNotImplemented)
		return
	}
}
