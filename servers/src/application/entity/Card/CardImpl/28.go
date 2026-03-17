package CardImpl

type Card28 struct {
	BaseCard
}

func NewCard28() *Card28 {
	return &Card28{}
}

func (c *Card28) Attack() {

}
func (c *Card28) Hurt() {
}

func (c *Card28) GetID() int {
	return 28
}
