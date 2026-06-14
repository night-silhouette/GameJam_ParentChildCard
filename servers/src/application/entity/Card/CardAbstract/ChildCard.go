package CardAbstract

import "pcc_card/application/entity/protocol"

type ChildCard interface {
	Card
	Check(pc protocol.ProtocolCardWithCtx) ChildCheckFunc
	Trigger(pc protocol.ProtocolCardWithCtx, UserId int)
}

type ChildCheckFunc func(pc protocol.ProtocolCardWithCtx) (bool, int)

func (cFunc ChildCheckFunc) Exec(pc protocol.ProtocolCardWithCtx) (bool, int) {
	return cFunc(pc)
}
