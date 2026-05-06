package userhandler

import (
	"pcc_card/global"
	"pcc_card/presentation/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type TimeReq struct {
	T1 int64 `json:"t1"`
}

func (u *User_handler_impl) TimeSync() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req TimeReq
		if err := c.ShouldBindWith(&req, binding.JSON); err != nil {
			response.Fail(c, global.ResponseInvalidReqParams)
		}
		tServer := time.Now().UnixMilli()
		response.Success(c, tServer)
	}
}

func (u *User_handler_impl) TimeDebug() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now().UnixMilli() + 1000*60*5
		response.Success(c, t)
	}
}
