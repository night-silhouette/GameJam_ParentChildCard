package CardImpl

type Card27 struct {
	BaseCard
}

func NewCard27() *Card27 {
	return &Card27{}
}

func (c *Card27) Attack() {

}
func (c *Card27) Hurt() {
}

func (c *Card27) GetID() int {
	return 27
}
