package CardAbstract

import "pcc_card/application/entity/protocol"

type ChildCard interface {
	Check(pc protocol.ProtocolCardWithCtx) ChildCheckFunc
	Trigger(pc protocol.ProtocolCardWithCtx)
}

type ChildCheckFunc func(pc protocol.ProtocolCardWithCtx) bool
