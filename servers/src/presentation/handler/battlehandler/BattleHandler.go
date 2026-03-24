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
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"pcc_card/presentation/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BattleHandler interface {
	handler.Handler
	BattleWs() gin.HandlerFunc
	DebugGetMachData() gin.HandlerFunc
	DebugBattleContainer() gin.HandlerFunc
}
type BattleHandlerImpl struct {
	s battleservice.BattleService
}

func (u *BattleHandlerImpl) Set_service(svc service.Service) {
	u.s = svc.(battleservice.BattleService)
}

func (u *BattleHandlerImpl) DebugGetMachData() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := battleservice.MatchPool.DebugGetMachData()
		response.Success(c, res)
	}
}

func (u *BattleHandlerImpl) DebugBattleContainer() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := battleservice.BC.GetBattleData()
		response.Success(c, res)
	}
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
		response.WsSuccess(conn, "ws连接成功")
		id := c.GetInt("id")
		transformAddMatchWithThis := make(chan battleservice.PlayerChannel)
		go u.AddMatch(c, conn, goctx, transformAddMatchWithThis)
		go u.ListenResquest(conn, id, goctx, transformAddMatchWithThis)

		playerChan := <-transformAddMatchWithThis
		go u.ListenResponse(conn, id, playerChan.ResponseChan, goctx)

	}

}

func (u *BattleHandlerImpl) ListenResquest(conn *websocket.Conn, id int, goctx context.Context, trans chan battleservice.PlayerChannel) {
	for {

		_, p, err := conn.ReadMessage()

		if err != nil {

			return
		}
		var playerC chan BattleDto.Action
		select {
		case playerChan := <-trans:
			playerC = playerChan.AcceptChan
		default:
		}

		decoder := json.NewDecoder(bytes.NewReader(p))
		decoder.DisallowUnknownFields() // 开启严苛模式

		var action BattleDto.Action
		err = decoder.Decode(&action)
		if err != nil {
			response.WsFail(conn, global.ResponseInvalidReqParams)
			continue
		}
		//action解析完成
		if action.ActionCode == BattleDto.CancelMatch {
			battleservice.MatchSignals.Delete(id)
			battleservice.MatchPool.Delete(id)
			response.WsSuccess(conn, "取消成功")
			time.Sleep(time.Millisecond * 200)
			conn.Close()
			return
		}
		select {
		case playerC <- action:
		case <-goctx.Done():
			return
		default:
			continue
		}

	}
}

func (u *BattleHandlerImpl) ListenResponse(conn *websocket.Conn, id int, playerC chan BattleDto.Action, goctx context.Context) {
	for {
		select {
		case Res, ok := <-playerC:
			if !ok {
				return
			}
			response.WsSuccess(conn, Res)
		case <-goctx.Done():
			return
		}

	}
}
