package CardImpl

type Card16 struct {
	BaseCard
}

func NewCard16() *Card16 {
	return &Card16{}
}

func (c *Card16) Attack() {

}
func (c *Card16) Hurt() {
}

func (c *Card16) GetID() int {
	return 16
}
