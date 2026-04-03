package userhandler

import (
	"fmt"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type User_handler interface {
	handler.Handler
	Get() gin.HandlerFunc
	Post() gin.HandlerFunc
	Put() gin.HandlerFunc
	Delete() gin.HandlerFunc
	Patch() gin.HandlerFunc
	GetAllOnePage() gin.HandlerFunc
	GetMailStatus() gin.HandlerFunc
	SendMail() gin.HandlerFunc
	ChangeMailStatus() gin.HandlerFunc
	DeleteMailByMailId() gin.HandlerFunc
	DeleteMailAll() gin.HandlerFunc
	UserVagueSearch() gin.HandlerFunc
}
type User_handler_impl struct {
	s UserService.User_service
}

func (u *User_handler_impl) Set_service(svc service.Service) {
	u.s = svc.(UserService.User_service)
}

func (u *User_handler_impl) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		is_admin := c.GetBool("is_admin")
		user_id := c.GetInt("id")

		if !is_admin {
			e, err := u.s.Find_user_by_id(user_id)
			if err != global.ResponseSuccess {
				response.Fail(c, err)
				return
			}
			response.Success(c, e)
			return
		}
		var req UserSearchReqDto
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
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
		e := &User_entity.User{Id: -1, Name: data.Name, Password: data.Password, Is_admin: false}
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
		is_admin := c.GetBool("is_admin")
		user_id := c.GetInt("id")
		var req UserPutReqDto
		if !is_admin && user_id != req.Id {
			response.Fail(c, global.ResponseForbidden)
			return
		}
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}

		e := &User_entity.User{Id: req.Id, Name: req.Name, Password: req.Password}
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
		is_admin := c.GetBool("is_admin")
		user_id := c.GetInt("id")

		var id int
		var req UserDeleteReqDto
		if err := c.ShouldBindQuery(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		id = req.ID
		if !is_admin && id != 0 {
			response.Fail(c, global.ResponseForbidden)
			return
		}

		err := u.s.Delete_user(user_id, c.Request.Context())

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
		user_id := c.GetInt("id")
		var req UserPatchDto
		if err := c.ShouldBind(&req); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		if req.Name == "" && req.Password == "" {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		e, err := u.s.Find_user_by_id(user_id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		if req.Name != "" {
			if req.Name == e.Name {
				response.Fail(c, global.ResponseInvalidReqParams)
				return
			}
			e.Name = req.Name
		}
		if req.Password != "" {
			OK := u.s.Check_password(user_id, req.Password)
			if OK == global.ResponseSuccess {
				response.Fail(c, global.ResponseInvalidReqParams)
				return
			} else {
				e.Password = req.Password
			}

		} else {
			e.Password = ""
		}
		err = u.s.Update_user(e)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}

func (u *User_handler_impl) UserVagueSearch() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := UserVagueSearchReq{}
		if err := c.ShouldBind(&req); err != nil {
			fmt.Println(err)
			return
		}
		err, uList := u.s.UserSearch(req.VagueName)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, uList)
	}
}
