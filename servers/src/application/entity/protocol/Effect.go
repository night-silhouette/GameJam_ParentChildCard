package protocol

type Effect interface {
	Execute(pc ProtocolCardWithCtx)
}
