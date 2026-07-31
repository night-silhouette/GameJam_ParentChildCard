package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/global"
	"time"
)

type Effect interface {
	Execute(pc ProtocolCardWithCtx)
}

type EffectBase struct {
	CheckIsInterrupt bool
	InterruptConfig  *InterruptConfig
}

type InterruptConfig struct {
	UserId        int
	Time          time.Duration
	TempIdList    []int
	SelectNum     int
	Res           chan []int
	CallTempId    int
	InterruptType BattleData.InterruptType
}

func (e *EffectBase) Execute(pc ProtocolCardWithCtx) {}

// 唤起中断
func (e *EffectBase) CallInterrupt(pc ProtocolCardWithCtx, InterruptCallback func(res []int, pc ProtocolCardWithCtx)) {
	dto := BattleData.NewInterruptDto(e.InterruptConfig.Time, e.InterruptConfig.TempIdList, e.InterruptConfig.SelectNum, e.InterruptConfig.InterruptType, e.InterruptConfig.CallTempId)
	pc.ProtoColInterrupt(e.InterruptConfig.UserId, dto, e.InterruptConfig.Res, e.InterruptConfig.Time)

	go func() {
		val := <-e.InterruptConfig.Res
		pc.ProtoColCancelInterrupt()
		InterruptCallback(val, pc)
	}()
}
func NewInterruptConfig(UserId int, TempIdList []int, SelectNum int, CallTempId int, InterruptType BattleData.InterruptType) InterruptConfig {
	config := InterruptConfig{}
	config.Res = make(chan []int)                //初始化res 管道
	config.Time = global.Interrupt * time.Second //统一中断的时间
	config.TempIdList = TempIdList
	config.UserId = UserId
	config.SelectNum = SelectNum
	config.CallTempId = CallTempId
	config.InterruptType = InterruptType
	return config
}
