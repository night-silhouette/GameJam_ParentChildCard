package CardImpl

type Card42 struct {
	BaseCard
}

func NewCard42() *Card42 {
	return &Card42{}
}

func (c *Card42) Attack() {

}
func (c *Card42) Hurt() {
}

func (c *Card42) GetID() int {
	return 42
}
