package battlehandler

import (
	"pcc_card/application/service"
	"pcc_card/application/service/battleservice"
	"pcc_card/global"
	"pcc_card/presentation/handler"
	"pcc_card/presentation/response"

	"github.com/gin-gonic/gin"
)

type BattleHandler interface {
	handler.Handler
	AddMatch() gin.HandlerFunc
}
type BattleHandlerImpl struct {
	s battleservice.BattleService
}

func (u *BattleHandlerImpl) Set_service(svc service.Service) {
	u.s = svc.(battleservice.BattleService)
}

func (u *BattleHandlerImpl) AddMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetInt("id")
		if !u.s.IsHasID(id) {
			u.s.AddMatch(id)
			response.Success(c, "进入匹配队列")
			matchSignals := u.s.GetMatchSignals()
			myChan := make(chan battleservice.MatchResult, 1)
			matchSignals.Store(id, myChan)
			defer battleservice.MatchSignals.Delete(id)
			select {
			case result := <-myChan:
				response.Success(c, result)
				return
			case <-c.Request.Context().Done():
				battleservice.MatchPool.Delete(id)
				return
			}

		} else {
			response.Fail(c, global.ResponseRepeatRequest)
			return
		}

	}

}

func (u *BattleHandlerImpl) LongPollingMatchState() {

}
