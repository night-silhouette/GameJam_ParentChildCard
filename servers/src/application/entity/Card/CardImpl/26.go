package CardImpl

type Card26 struct {
	BaseCard
}

func NewCard26() *Card26 {
	return &Card26{}
}

func (c *Card26) Attack() {

}
func (c *Card26) Hurt() {
}

func (c *Card26) GetID() int {
	return 26
}
