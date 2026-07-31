package CardMeta

type ControlSignal int

const (
	Wound ControlSignal = iota //被击伤
)

type BroadInfo struct {
	ControlSignal ControlSignal
	Caller        int //引发者
	Publisher     int //发布者
}

// 没用的字段就传-1吧
func NewBroadInfo(signal ControlSignal, Caller int, Publisher int) *BroadInfo {
	res := BroadInfo{}
	res.ControlSignal = signal
	res.Caller = Caller
	res.Publisher = Publisher
	return &res
}
