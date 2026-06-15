package protocol

import "pcc_card/application/entity/BattleData"

type Weather int

const (
	ningjing Weather = iota
	shabao           //沙暴
	ganmu
	mingyang
	miwu
	dashu
	zhipou
	shuangjiang
	shutu
	fengdu
)

// 每次有新天气,改一下这个值,否则不好做随机
const WeatherCanSelectNum = 9

var WeatherFuncMap map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed) //传入在战斗的卡的数组

func InitWeatherFuncMap() {
	WeatherFuncMap = make(map[Weather]func(pc ProtocolCardWithCtx, card []BuffNeed), WeatherCanSelectNum+1)
	WeatherFuncMap[ningjing] = func(pc ProtocolCardWithCtx, card []BuffNeed) {} //无效果
	WeatherFuncMap[shabao] = func(pc ProtocolCardWithCtx, card []BuffNeed) {

		for _, c := range card {
			pc.ProtoColPush(NewCustom(func(pc ProtocolCardWithCtx) {
				pc.ProtoColAttackNoHurt(c.GetTempId(), 1, BattleData.TrueDamage)
			}))
		}

	}
	WeatherFuncMap[ganmu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {
		for _, c := range card {
			TempId := c.GetTempId()
			pc.ProtoColPush(NewHeal(&TempId, 1, c.GetDec()))
		}
	}
	WeatherFuncMap[mingyang] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherFuncMap[miwu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherFuncMap[dashu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherFuncMap[zhipou] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherFuncMap[shuangjiang] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherFuncMap[shutu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
	WeatherFuncMap[fengdu] = func(pc ProtocolCardWithCtx, card []BuffNeed) {}
}
