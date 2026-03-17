package CardImpl

type Card31 struct {
	BaseCard
}

func NewCard31() *Card31 {
	return &Card31{}
}

func (c *Card31) Attack() {

}
func (c *Card31) Hurt() {
}

func (c *Card31) GetID() int {
	return 31
}
