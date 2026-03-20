package tokenhandler

import (
	"net/http"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

type Token_handler interface {
	handler.Handler
	Middleware_token_check() gin.HandlerFunc
	Post() gin.HandlerFunc
	Get() gin.HandlerFunc
	Put() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
}
type Token_handler_impl struct {
	s UserService.User_service
}

func (u *Token_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(UserService.User_service)
}
func (u *Token_handler_impl) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, "验证成功")
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
		token, err := u.s.Release_token(id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
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

func (u *Token_handler_impl) Middleware_token_check() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		if path == "/v1/user/" && method == http.MethodPost {
			c.Next()
			return
		}
		if path == "/ping" {
			c.Next()
			return
		}
		if path == "/v1/token/" && method == http.MethodPost {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			response.Fail(c, global.ResponseTokenMissing)
			c.Abort()
			return
		}
		id, is_admin, err := u.s.Is_valid_token(token)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			c.Abort()
			return
		}
		c.Set("id", id)
		c.Set("is_admin", is_admin)

		c.Next()
	}
}
