package CardImpl

type Card36 struct {
	BaseCard
}

func NewCard36() *Card36 {
	return &Card36{}
}

func (c *Card36) Attack() {

}
func (c *Card36) Hurt() {
}

func (c *Card36) GetID() int {
	return 36
}
