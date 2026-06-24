package protocol

import "pcc_card/application/entity/BattleData"

func InitWeather() {
	InitWeatherOnaApplyMap()
	InitWeatherFuncMap()
}

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

var WeatherOnaApplyMap map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed)

func InitWeatherOnaApplyMap() {
	WeatherOnaApplyMap = make(map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed))
	WeatherOnaApplyMap[Ningjing] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Shabao] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Ganmu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Miwu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Dashu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Zhipou] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Shuangjiang] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Shutu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherOnaApplyMap[Fengdu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}

}

var WeatherFuncMap map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed) //传入在战斗的卡的数组

func InitWeatherFuncMap() {
	WeatherFuncMap = make(map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed), WeatherCanSelectNum+1)
	WeatherFuncMap[Ningjing] = func(pc ProtocolCardWithCtx, card []BuffNeed) {} //无效果
	WeatherFuncMap[Shabao] = func(pc ProtocolCardWithCtx, card []BuffNeed) {    //所有人收到一点真伤

		for _, c := range card {
			pc.ProtoColPush(NewCustom(func(pc ProtocolCardWithCtx) {
				pc.ProtoColAttackNoHurt(c.GetTempId(), 1, BattleData.TrueDamage)
			}))
		}

	}
	WeatherFuncMap[Ganmu] = func(pc ProtocolCardWithCtx, card []BuffNeed) { //出战牌回血
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewHeal(&TempId, 1, c.GetDec(), false, &InterruptConfig{}))
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
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(XuFeng, 1, 0.15, pc.CreateTempId()), false, &InterruptConfig{}))
		}
	}
	WeatherFuncMap[Dashu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(HealingDecay, 1, 0.3, pc.CreateTempId()), false, &InterruptConfig{}))
		}
	}
	WeatherFuncMap[Zhipou] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		if !pc.GetWinnerIsAction() {
			for _, c := range card {
				if c.GetOwnerId() != pc.GetWinnerId() { //跳掉输的人的牌
					continue
				}

				pc.ProtoColPush(NewCustom(func(pc ProtocolCardWithCtx) {
					pc.ProtoColAttackNoHurt(c.GetTempId(), 2, BattleData.Damage)
				}))
			}
		}
	}
	WeatherFuncMap[Shuangjiang] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(Vulnerability, 1, 0.4, pc.CreateTempId()), false, &InterruptConfig{}))
		}
	}
	WeatherFuncMap[Shutu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewGiveBuff(&TempId, *NewBuffBase(DamageImmunity, 1, 0.28, pc.CreateTempId()), false, &InterruptConfig{}))
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
