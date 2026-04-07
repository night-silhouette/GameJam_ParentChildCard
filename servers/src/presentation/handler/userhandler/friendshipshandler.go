package userhandler

import (
	"pcc_card/global"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (u *User_handler_impl) CreateFriendship() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req PostFriendshipsQeq
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		id := c.GetInt("id")
		err := u.s.AddFriendshipsRequest(c.Request.Context(), id, req.Id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}
func (u *User_handler_impl) GetFriendships() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		err, idMap := u.s.FindFriendships(c.Request.Context(), id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, idMap)
	}
}

func (u *User_handler_impl) DeleteFriendships() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		var req DeleteFriendshipsQeq
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
			return
		}
		err := u.s.DeleteFriendships(c.Request.Context(), id, req.Id)
		if err != global.ResponseSuccess {
			response.Fail(c, err)
			return
		}
		response.Success(c, "ok")
	}
}
