package battlehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pcc_card/application/service"
	"pcc_card/application/service/battleservice"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BattleHandler interface {
	handler.Handler
	BattleWs() gin.HandlerFunc
}
type BattleHandlerImpl struct {
	s battleservice.BattleService
}

func (u *BattleHandlerImpl) Set_service(svc service.Service) {
	u.s = svc.(battleservice.BattleService)
}

func (u *BattleHandlerImpl) AddMatch(c *gin.Context, conn *websocket.Conn, goctx context.Context, res chan battleservice.PlayerChannel) {
	id := c.GetInt("id")
	if !u.s.IsHasID(id) {
		u.s.AddMatch(id)
		matchSignals := u.s.GetMatchSignals()
		myChan := make(chan battleservice.MatchResult, 1)
		matchSignals.Store(id, myChan)
		defer battleservice.MatchSignals.Delete(id)
		select {
		case result := <-myChan:
			Bt := battleservice.BC.GetBattleByUserID(id)
			playerChan := Bt.GetPlayerChanByUserID(id)
			res <- playerChan
			response.WsSuccess(conn, result)
			return
		case <-goctx.Done():
			battleservice.MatchPool.Delete(id)
			return
		}

	} else {
		response.WsFail(conn, global.ResponseRepeatRequest)
		return
	}
}

func (u *BattleHandlerImpl) BattleWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("升级失败:", err)
			return
		}
		goctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()
		defer conn.Close()
		//升级逻辑完成

		id := c.GetInt("id")
		transformAddMatchWithThis := make(chan battleservice.PlayerChannel)
		go u.AddMatch(c, conn, goctx, transformAddMatchWithThis)
		go u.ListenCancelMatch(conn, id)
		playerChan := <-transformAddMatchWithThis
		fmt.Println(playerChan)
	}

}

func (u *BattleHandlerImpl) ListenCancelMatch(conn *websocket.Conn, id int) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			response.WsFail(conn, global.ResponseUnknownError)
			return
		}
		decoder := json.NewDecoder(bytes.NewReader(p))
		decoder.DisallowUnknownFields() // 开启严苛模式

		var action Action
		err = decoder.Decode(&action)
		if err != nil {
			response.WsFail(conn, global.ResponseInvalidReqParams)
		}
		//action解析完成

		//取消匹配
		if action.CancelMatch {
			battleservice.MatchSignals.Delete(id)
			battleservice.MatchPool.Delete(id)
			response.WsSuccess(conn, "取消成功")
			time.Sleep(time.Millisecond * 200)
			conn.Close()
			return
		}
	}
}
