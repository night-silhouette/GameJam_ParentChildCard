package CardAbstract

type Card interface {
	GetID() int
	SetInfo(info map[string]any)
	GetInfo() map[string]any
	Clone() Card
}
