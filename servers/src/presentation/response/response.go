package response

import (
	"net/http"
	"pcc_card/global"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type response struct {
	Code global.ResponseStatusCode `json:"code"`
	Msg  string                    `json:"msg"`
	Data any                       `json:"data"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response{
		Code: global.ResponseSuccess,
		Msg:  global.StatusMsg[global.ResponseSuccess],
		Data: data,
	})
}

func Fail(c *gin.Context, code global.ResponseStatusCode) {
	c.JSON(http.StatusOK, response{
		Code: code,
		Msg:  global.StatusMsg[code],
		Data: gin.H{},
	})
}

func WsSuccess(conn *websocket.Conn, data any) {
	conn.WriteJSON(response{
		Code: global.ResponseSuccess,
		Msg:  global.StatusMsg[global.ResponseSuccess],
		Data: data,
	})
}

func WsFail(conn *websocket.Conn, code global.ResponseStatusCode) {
	conn.WriteJSON(response{
		Code: code,
		Msg:  global.StatusMsg[code],
		Data: gin.H{},
	})
}
