package protocol

import "pcc_card/application/entity/BattleData"

type Weather int

const (
	Ningjing Weather = iota
	Shabao           //沙暴
	Ganmu
	Mingyang
	Miwu
	Dashu
	Zhipou
	Shuangjiang
	Shutu
	Fengdu
)

// 每次有新天气,改一下这个值,否则不好做随机
const WeatherCanSelectNum = 9

var WeatherFuncMap map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed) //传入在战斗的卡的数组

func InitWeatherFuncMap() {
	WeatherFuncMap = make(map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed), WeatherCanSelectNum+1)
	WeatherFuncMap[Ningjing] = func(pc ProtocolCardWithCtx, card []BuffNeed) {} //无效果
	WeatherFuncMap[Shabao] = func(pc ProtocolCardWithCtx, card []BuffNeed) { //所有人收到一点真伤

		for _, c := range card {
			pc.ProtoColPush(NewCustom(func(pc ProtocolCardWithCtx) {
				pc.ProtoColAttackNoHurt(c.GetTempId(), 1, BattleData.TrueDamage)
			}))
		}

	}
	WeatherFuncMap[Ganmu] = func(pc ProtocolCardWithCtx, card []BuffNeed) { //出战牌回血
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewHeal(&TempId, 1, c.GetDec()))
		}
	}
	WeatherFuncMap[Mingyang] = func(pc ProtocolCardWithCtx, card []BuffNeed) { //
		UserIdList := pc.GetIds()
		for _, UserId := range UserIdList {
			pc.ProtoColPush(NewUpdateEnergy(UserId, 1))
		}

	}
	WeatherFuncMap[Miwu] = func(pc ProtocolCardWithCtx, card []BuffNeed) { //每个人有闪避
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(XuFeng, 1, 0.15, pc.CreateTempId())))
		}
	}
	WeatherFuncMap[Dashu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(HealingDecay, 1, 0.3, pc.CreateTempId())))
		}
	}
	WeatherFuncMap[Zhipou] = func(pc ProtocolCardWithCtx, card []BuffNeed) {

	}
	WeatherFuncMap[Shuangjiang] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(Vulnerability, 1, 0.4, pc.CreateTempId())))
		}
	}
	WeatherFuncMap[Shutu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(DamageImmunity, 1, 0.28, pc.CreateTempId())))
		}
	}
	WeatherFuncMap[Fengdu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		
	}
}

type ExecPosition int

const (
	RoundStart ExecPosition = iota
	RoundEnd
)

var WeatherExecPositionMap map[Weather]ExecPosition = map[Weather]ExecPosition{
	Ningjing:    RoundStart,
	Shabao:      RoundEnd,
	Ganmu:       RoundEnd,
	Mingyang:    RoundStart,
	Miwu:        RoundStart,
	Dashu:       RoundStart,
	Zhipou:      RoundStart,
	Shuangjiang: RoundStart,
	Shutu:       RoundStart,
	Fengdu:      RoundStart,
}
