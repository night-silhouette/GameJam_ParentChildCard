package CardImpl

type Card21 struct {
	BaseCard
}

func NewCard21() *Card21 {
	return &Card21{}
}

func (c *Card21) Attack() {

}
func (c *Card21) Hurt() {
}

func (c *Card21) GetID() int {
	return 21
}
