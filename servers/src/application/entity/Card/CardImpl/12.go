package CardImpl

type Card12 struct {
	BaseCard
}

func NewCard12() *Card12 {
	return &Card12{}
}

func (c *Card12) Attack() {

}
func (c *Card12) Hurt() {
}

func (c *Card12) GetID() int {
	return 12
}
