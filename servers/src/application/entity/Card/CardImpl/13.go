package CardImpl

type Card13 struct {
	BaseCard
}

func NewCard13() *Card13 {
	return &Card13{}
}

func (c *Card13) Attack() {

}
func (c *Card13) Hurt() {
}

func (c *Card13) GetID() int {
	return 13
}
