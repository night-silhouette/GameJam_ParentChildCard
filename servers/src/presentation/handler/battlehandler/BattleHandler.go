package battlehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/service"
	"pcc_card/application/service/UserService"
	"pcc_card/application/service/battleservice"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"pcc_card/presentation/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type BattleHandler interface {
	handler.Handler
	BattleWs() gin.HandlerFunc
	DebugGetMachData() gin.HandlerFunc
	DebugBattleContainer() gin.HandlerFunc
	WsReconnect() gin.HandlerFunc
}
type BattleHandlerImpl struct {
	s       battleservice.BattleService
	User_s  UserService.User_service
	writeMu sync.Mutex
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

func (u *BattleHandlerImpl) AddMatch(c *gin.Context, conn *websocket.Conn, goctx context.Context, res chan battleservice.PlayerChannel, data BattleData.EnterBtData) {
	id := c.GetInt("id")

	if !u.s.IsHasID(id) {
		u.s.AddMatch(id, data)
		matchSignals := u.s.GetMatchSignals()
		myChan := make(chan battleservice.MatchResult, 1)
		matchSignals.Store(id, myChan)
		defer battleservice.MatchSignals.Delete(id)
		select {
		case <-myChan:
			Bt := battleservice.BC.GetBattleByUserID(id)
			playerChan := Bt.GetPlayerChanByUserID(id)
			res <- playerChan
			res <- playerChan //因为最初这里的子线程信息传递设计比较乱，造成了有极强的耦合。由于分发playerChan的管道，外面有两个接受者，所以
			//这里要塞入两份指针，一份给req，还有一份给response

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

func (u *BattleHandlerImpl) WsUpdate(c *gin.Context) *websocket.Conn {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("升级失败:", err)
		return nil
	}
	return conn
}

func (u *BattleHandlerImpl) BattleWs() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := u.WsUpdate(c)
		if conn == nil {
			return
		}
		//升级逻辑完成

		id := c.GetInt("id")
		HandlerCtx, cancel := context.WithCancel(c.Request.Context())
		//----------------------生命周期管理----------------------
		defer func() {

			u.NotifyOpOffline(id) //掉线断链的通知

			Util.CreateTimer(time.Second*2, func() {
				cancel()
				conn.Close() //两秒之后再断，这样前端可以收到死掉的信息
			})

		}()
		//----------------------生命周期管理----------------------

		//检查是不是加入战斗了
		_, HaveBt := u.s.CheckUserIdIsBattle(HandlerCtx, id)
		fmt.Println("IsHaveBt", HaveBt)
		if HaveBt == global.ResponseSuccess { //找到了,也就是有这个战斗
			response.WsFailWithMsg(conn, global.ResponseRepeatRequest, "已经在战斗了,断线重连用另外那个接口")
			return
		}
		//检查是不是加入战斗了

		//-------------------对带入的card和gold信息解析---------------------
		BtData := c.Query("btData")
		decodedJson, err := url.QueryUnescape(BtData)
		if err != nil {
			fmt.Println(err)
			response.WsFail(conn, global.ResponseInvalidReqParams)
			return
		}
		var data BattleData.EnterBtData

		err = json.Unmarshal([]byte(decodedJson), &data)
		if err != nil {
			fmt.Println(err)
			response.WsFail(conn, global.ResponseInvalidReqParams)
			return
		}
		fmt.Println(data)
		//检验data合法性
		ok := u.User_s.CheckBtDataIsValid(HandlerCtx, id, data.CardList, data.Gold)
		if ok != global.ResponseSuccess {
			fmt.Println(ok)
			response.WsFail(conn, ok)
			return
		}
		//-------------------对带入的card和gold信息解析---------------------

		fmt.Println("数据合法")

		transformAddMatchWithThis := make(chan battleservice.PlayerChannel, 2)
		go u.AddMatch(c, conn, HandlerCtx, transformAddMatchWithThis, data)
		go u.ListenRequest(conn, id, HandlerCtx, transformAddMatchWithThis, cancel)

		select {
		case <-HandlerCtx.Done():
			return
		case playerChan := <-transformAddMatchWithThis:
			go u.ListenResponse(conn, playerChan.ResponseChan, HandlerCtx)
		}

		select { //阻塞，不让handler直接结束
		case <-HandlerCtx.Done():
			return
		}

	}

}

func (u *BattleHandlerImpl) ListenRequest(conn *websocket.Conn, id int, goctx context.Context, trans chan battleservice.PlayerChannel, cancelFunc context.CancelFunc) {
	var playerC chan BattleDto.Action
	//拦截器
	//Interceptor := Util.NewInterceptor(global.WsInterceptorTime * time.Millisecond)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			cancelFunc()
			return
		}
		//todo 拦截器
		//if Interceptor.ShouldBlock(p) {
		//	continue
		//}

		select {
		case playerChan := <-trans:
			fmt.Println(playerChan)
			playerC = playerChan.AcceptChan
		default:
		} //每一次去监听match有不有把playerC给这

		flag := u.ListenRequestConn(p, conn, cancelFunc, id, goctx, playerC)
		if !flag {
			return
		}
	}
}

// 阻塞的,读conn,解析action,传给playerc(playerc是传进来的管道)
func (u *BattleHandlerImpl) ListenRequestConn(p []byte, conn *websocket.Conn, cancelFunc context.CancelFunc, UserId int, goctx context.Context, playerC chan BattleDto.Action) bool {
	decoder := json.NewDecoder(bytes.NewReader(p))
	decoder.DisallowUnknownFields() // 开启严苛模式
	var action BattleDto.Action
	err := decoder.Decode(&action)
	if err != nil {
		u.writeMu.Lock()
		response.WsFailWithErr(conn, global.ResponseInvalidReqParams, err)
		u.writeMu.Unlock()
		return true
	}
	//action解析完成
	if action.ActionCode == BattleDto.CancelMatch && action.Predicates == BattleDto.Notify {
		fmt.Println("取消匹配")
		battleservice.MatchSignals.Delete(UserId)
		battleservice.MatchPool.Delete(UserId)
		response.WsSuccess(conn, "取消成功")
		time.Sleep(time.Millisecond * 200)
		conn.Close()
		return false
	}

	select {
	case playerC <- action:
	case <-goctx.Done():
		return false
	default:

		if playerC == nil {
			u.writeMu.Lock()
			response.WsFailWithMsg(conn, global.BattleInvalidTiming, "正在匹配中")
			u.writeMu.Unlock()
		}
	}
	return true
}

func (u *BattleHandlerImpl) ListenResponse(conn *websocket.Conn, playerC chan BattleDto.Action, goctx context.Context) {
	for {
		select {
		case Res, _ := <-playerC:

			if Res.ActionCode == BattleDto.Fault {
				code := Res.ActionData.(global.ResponseStatusCode)
				u.writeMu.Lock()
				response.WsFail(conn, code)
				u.writeMu.Unlock()
				continue
			} //监听内部错误
			u.writeMu.Lock()
			response.WsSuccess(conn, Res) //直接返回action
			u.writeMu.Unlock()
		case <-goctx.Done():
			return
		}

	}
}

// 有一方异常退出的时候,给另一方消息通知他离线了(有BT的nil解释)//传入自己的userid
func (u *BattleHandlerImpl) NotifyOpOffline(UserId int) {
	Bt := battleservice.BC.GetBattleByUserID(UserId)
	if Bt == nil {
		return
	}
	ctx := Bt.Ctx
	OpId := ctx.GetOpponentId(UserId)
	Bt.SM.SendActionById(OpId, BattleDto.NewAction(BattleDto.OpOffline, BattleDto.Notify, ""))
}

func (u *BattleHandlerImpl) NotifyOnOnline(UserId int) {
	Bt := battleservice.BC.GetBattleByUserID(UserId)
	if Bt == nil {
		return
	}
	ctx := Bt.Ctx
	OpId := ctx.GetOpponentId(UserId)
	Bt.SM.SendActionById(OpId, BattleDto.NewAction(BattleDto.OpOnline, BattleDto.Notify, ""))
}

// 通知双方,传入任意一个id就可以
func (u *BattleHandlerImpl) NotifyOverBattle(UserId int) {
	Bt := battleservice.BC.GetBattleByUserID(UserId)
	if Bt == nil {
		return
	}
	ctx := Bt.Ctx
	OpId := ctx.GetOpponentId(UserId)
	Bt.SM.SendActionById(OpId, BattleDto.NewAction(BattleDto.OverBattle, BattleDto.Notify, ""))
	Bt.SM.SendActionById(OpId, BattleDto.NewAction(BattleDto.OverBattle, BattleDto.Notify, ""))
}

func (u *BattleHandlerImpl) WsReconnect() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := u.WsUpdate(c)
		if conn == nil {
			return
		}
		UserId := c.GetInt("id")
		HandlerCtx, HandlerCancel := context.WithCancel(c.Request.Context())
		defer func() {
			u.NotifyOpOffline(UserId) //掉线断链的通知
			Util.CreateTimer(time.Second*2, func() {
				HandlerCancel()
				conn.Close() //两秒之后再断，这样前端可以收到死掉的信息
			})
		}()

		_, HaveBt := u.s.CheckUserIdIsBattle(HandlerCtx, UserId)
		if HaveBt != global.ResponseSuccess {
			fmt.Println("没在战斗啊,别重连")
			return
		}
		Bt := battleservice.BC.GetBattleByUserID(UserId)
		if Bt == nil {
			fmt.Println("重连找不到战斗battle")
			return
		}
		Nt := Bt.GetPlayerChanByUserID(UserId)

		//接受
		go u.ListenResponse(conn, Nt.ResponseChan, HandlerCtx)
		go func() {
			for {
				_, p, err := conn.ReadMessage()
				if err != nil {
					HandlerCancel()
					return
				}
				//todo 拦截器
				//if Interceptor.ShouldBlock(p) {
				//	continue
				//}
				u.ListenRequestConn(p, conn, HandlerCancel, UserId, HandlerCtx, Nt.AcceptChan)
			}
		}()

		u.NotifyOnOnline(UserId)

		select { //阻塞，不让handler直接结束
		case <-HandlerCtx.Done():
			return
		}

	}
}
