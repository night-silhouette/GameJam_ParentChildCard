package CardImpl

type Card15 struct {
	BaseCard
}

func NewCard15() *Card15 {
	return &Card15{}
}

func (c *Card15) Attack() {

}
func (c *Card15) Hurt() {
}

func (c *Card15) GetID() int {
	return 15
}
