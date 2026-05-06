package protocolCardWithCtx

type Effect interface {
	Execute()
}

type Death struct {
	DeathId int
}

func (Death *Death) Execute() {}
