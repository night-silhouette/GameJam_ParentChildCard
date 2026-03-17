package CardImpl

type Card18 struct {
	BaseCard
}

func NewCard18() *Card18 {
	return &Card18{}
}

func (c *Card18) Attack() {

}
func (c *Card18) Hurt() {
}

func (c *Card18) GetID() int {
	return 18
}
