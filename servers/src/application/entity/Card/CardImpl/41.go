package CardImpl

type Card41 struct {
	BaseCard
}

func NewCard41() *Card41 {
	return &Card41{}
}

func (c *Card41) Attack() {

}
func (c *Card41) Hurt() {
}

func (c *Card41) GetID() int {
	return 41
}
