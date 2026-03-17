package CardImpl

type Card08 struct {
	BaseCard
}

func NewCard08() *Card08 {
	return &Card08{}
}

func (c *Card08) Attack() {

}
func (c *Card08) Hurt() {
}

func (c *Card08) GetID() int {
	return 8
}
