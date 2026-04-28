package protocolCardWithCtx

import "pcc_card/application/entity/Effect"

type ProtocolCardWithCtx interface {
	ProtoColPush(e Effect.Effect)
}
