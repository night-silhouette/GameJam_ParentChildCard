package CardImpl

type Card17 struct {
	BaseCard
}

func NewCard17() *Card17 {
	return &Card17{}
}

func (c *Card17) Attack() {

}
func (c *Card17) Hurt() {
}

func (c *Card17) GetID() int {
	return 17
}
