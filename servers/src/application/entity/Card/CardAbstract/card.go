package CardAbstract

type Card interface {
	GetID() int
	GetInfo() map[string]any
}
