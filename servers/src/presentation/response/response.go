package response

import (
	"net/http"
	"pcc_card/global"

	"github.com/gin-gonic/gin"
)

type response struct {
	Code global.StatusCode `json:"code"`
	Msg  string            `json:"msg"`
	Data any               `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{
		Code: global.StatusSuccess,
		Msg:  global.StatusMsg[global.StatusSuccess],
		Data: data,
	})
}

func Fail(c *gin.Context, code global.StatusCode) {
	c.JSON(http.StatusOK, response{
		Code: code,
		Msg:  global.StatusMsg[code],
		Data: nil,
	})
}
