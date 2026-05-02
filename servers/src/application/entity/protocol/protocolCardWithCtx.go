package protocol

type ProtocolCardWithCtx interface {
	ProtoColPush(e Effect)
	ProtoColSetCardBtHp(UserId int, tempId int, NowHp float64)
}
