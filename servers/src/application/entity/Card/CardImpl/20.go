package CardImpl

type Card20 struct {
	BaseCard
}

func NewCard20() *Card20 {
	return &Card20{}
}

func (c *Card20) Attack() {

}
func (c *Card20) Hurt() {
}

func (c *Card20) GetID() int {
	return 20
}
