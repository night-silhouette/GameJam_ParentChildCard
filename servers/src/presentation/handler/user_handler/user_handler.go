package user_handler

import (
	"pcc_card/application/entity"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

type User_handler interface {
	handler.Handler
}
type User_handler_impl struct {
	s UserService.User_service
}

func (u *User_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(UserService.User_service)
}

func (u *User_handler_impl) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserSearchReqDto
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}

		if req.ID != 0 {
			e, status := u.s.Find_user_by_id(req.ID)
			if status == global.ResponseSuccess {
				response.Success(c, e)
				return
			} else {
				response.Fail(c, status)
				return
			}
		} else if req.Name != "" {
			e, status := u.s.Find_user_by_name(req.Name)
			if status == global.ResponseSuccess {
				response.Success(c, e)
				return
			} else {
				response.Fail(c, status)
				return
			}
		} else {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}

	}
}

func (u *User_handler_impl) Post() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data UserPostDto
		if err := c.ShouldBind(&data); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		if data.Name == "" || data.Password == "" {
			response.Fail(c, global.ResponseRequiredParamsMissing)
			return
		}
		e := &entity.User{Id: -1, Name: data.Name, Password: data.Password}
		err := u.s.Create_user(e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
		return

	}
}

func (u *User_handler_impl) Put() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserPutReqDto
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		e := &entity.User{Id: req.Id, Name: req.Name, Password: req.Password}
		err := u.s.Update_user(e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
		return

	}
}

func (u *User_handler_impl) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserDeleteReqDto
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		e, _ := u.s.Find_user_by_id(req.ID)
		err := u.s.Delete_user(e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
		return
	}
}

func (u *User_handler_impl) Patch() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Fail(c, global.ResponseNotImplemented)
	}
}
