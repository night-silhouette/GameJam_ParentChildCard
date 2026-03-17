package CardImpl

type Card30 struct {
	BaseCard
}

func NewCard30() *Card30 {
	return &Card30{}
}

func (c *Card30) Attack() {

}
func (c *Card30) Hurt() {
}

func (c *Card30) GetID() int {
	return 30
}
