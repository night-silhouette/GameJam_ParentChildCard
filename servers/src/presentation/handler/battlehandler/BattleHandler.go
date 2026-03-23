package battlehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
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
		CancelMatchContext, CancelMatchCancel := context.WithCancel(context.Background())
		go u.ListenCancelMatch(conn, id, CancelMatchContext)

		playerChan := <-transformAddMatchWithThis
		CancelMatchCancel()

		go func() { //listen客户端
			for {
				_, p, err := conn.ReadMessage()
				if err != nil {
					return
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
				select {
				case playerChan.AcceptChan <- action:
				case <-goctx.Done():
					return
				}
			}
		}()
		go func() {
			for {
				select {
				case Res, ok := <-playerChan.ResponseChan:
					if !ok {
						return
					}
					response.WsSuccess(conn, Res)
				case <-goctx.Done():
					return
				}

			}
		}()

	}

}

func (u *BattleHandlerImpl) ListenCancelMatch(conn *websocket.Conn, id int, Goctx context.Context) {
	defer fmt.Println("【系统】匹配取消监听协程已安全退出，归还连接控制权")
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
		_, p, err := conn.ReadMessage()
		if Goctx.Err() != nil {
			// 临走前把闹钟关掉（恢复成永不超时）
			conn.SetReadDeadline(time.Time{})
			return
		}
		if err != nil {
			// 如果只是因为闹钟响了（超时），那就 continue 回到循环顶端检查 Context
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			// 如果是真正的连接断开，那就彻底退出
			return
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

		//取消匹配
		if action.ActionCode == BattleDto.CancelMatch {
			battleservice.MatchSignals.Delete(id)
			battleservice.MatchPool.Delete(id)
			response.WsSuccess(conn, "取消成功")
			time.Sleep(time.Millisecond * 200)
			conn.Close()
			return
		}
	}
}
