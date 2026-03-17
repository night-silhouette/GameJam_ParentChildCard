package CardImpl

type Card19 struct {
	BaseCard
}

func NewCard19() *Card19 {
	return &Card19{}
}

func (c *Card19) Attack() {

}
func (c *Card19) Hurt() {
}

func (c *Card19) GetID() int {
	return 19
}
