package handler

import (
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/global"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

type User_handler interface {
	Handler
}
type User_handler_impl struct {
	s UserService.User_service
}

func (u *User_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(UserService.User_service)
}

type UserSearchReq struct {
	ID   int    `form:"id"`
	Name string `form:"name"`
}

func (u *User_handler_impl) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserSearchReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.StatusInvalidReqParamsClass)
		}
		if req.ID != 0 {
			e, status := u.s.Find_user_by_id(req.ID)
			if status == global.StatusSuccess {
				response.Success(c, e)
			} else {
				response.Fail(c, status)
			}
		} else if req.Name != "" {
			e, status := u.s.Find_user_by_name(req.Name)
			if status == global.StatusSuccess {
				response.Success(c, e)
			} else {
				response.Fail(c, status)
			}
		} else {
			response.Fail(c, global.StatusInvalidReqParamsName)
		}

	}
}

func (u *User_handler_impl) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, "pong")
	}
}

func (u *User_handler_impl) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, "pong")
	}
}

func (u *User_handler_impl) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, "pong")
	}
}

func (u *User_handler_impl) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Success(c, "pong")
	}
}
