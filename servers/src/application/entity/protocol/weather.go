package protocol

type Weather int

const (
	ningjing Weather = iota
	shabao
	ganmu
	mingyang
	miwu
	dashu
	zhipou
	shuangjiang
	shutu
	fengdu
)

const WeatherCanSelectNum = 9

var WeatherFuncMap map[Weather]func(pc ProtocolCardWithCtx)

func InitWeatherFuncMap() {
	WeatherFuncMap = make(map[Weather]func(pc ProtocolCardWithCtx), WeatherCanSelectNum+1)
	WeatherFuncMap[ningjing] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[shabao] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[ganmu] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[mingyang] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[miwu] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[dashu] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[zhipou] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[shuangjiang] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[shutu] = func(pc ProtocolCardWithCtx) {}
	WeatherFuncMap[fengdu] = func(pc ProtocolCardWithCtx) {}
}
