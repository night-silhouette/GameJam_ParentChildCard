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
		fmt.Println(res)
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
			res <- playerChan //因为最初这里的子线程信息传递设计比较乱，造成了有极强的耦合。由于分发playerChan的管道，外面有两个接受者，所以
			//这里要塞入两份指针，一份给req，还有一份给response

			response.WsSuccess(conn, BattleDto.NewAction(BattleDto.StartBattle, BattleDto.Notify, result))
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
		transformAddMatchWithThis := make(chan battleservice.PlayerChannel, 2)
		go u.AddMatch(c, conn, goctx, transformAddMatchWithThis)
		go u.ListenResquest(conn, id, goctx, transformAddMatchWithThis, cancel)

		playerChan := <-transformAddMatchWithThis
		var OverGameChan chan bool = make(chan bool, 1)
		go u.ListenResponse(conn, id, playerChan.ResponseChan, goctx, OverGameChan)
		ret, _ := <-OverGameChan
		if ret {
			return
		}
	}

}

func (u *BattleHandlerImpl) ListenResquest(conn *websocket.Conn, id int, goctx context.Context, trans chan battleservice.PlayerChannel, cancelFunc context.CancelFunc) {
	var playerC chan BattleDto.Action
	go func() {
		select {
		case <-goctx.Done():
			return
		case playerChan := <-trans:
			playerC = playerChan.AcceptChan
		}
	}()
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			cancelFunc()
			return
		}

		decoder := json.NewDecoder(bytes.NewReader(p))
		decoder.DisallowUnknownFields() // 开启严苛模式

		var action BattleDto.Action
		err = decoder.Decode(&action)
		if err != nil {
			response.WsFailWithErr(conn, global.ResponseInvalidReqParams, err)
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
			if playerC == nil {
				response.WsFailWithMsg(conn, global.BattleInvalidTiming, "正在匹配中")
			}
		}

	}
}

func (u *BattleHandlerImpl) ListenResponse(conn *websocket.Conn, id int, playerC chan BattleDto.Action, goctx context.Context, OverGamechan chan bool) {
	for {
		select {
		case Res, _ := <-playerC:

			if Res.ActionCode == BattleDto.OverBattle {
				OverGamechan <- true
			}
			response.WsSuccess(conn, Res)
		case <-goctx.Done():
			return
		}

	}
}
