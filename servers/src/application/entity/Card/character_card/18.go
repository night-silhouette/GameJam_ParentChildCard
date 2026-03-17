package character_card

type Card18 struct {
	CharacterCardTemplate
}

func NewCard18() *Card18 {
	return &Card18{}
}

func (c *Card18) Attack() {

}
func (c *Card18) Hurt() {
}

func (c *Card18) GetID() int {
	return 18
}
