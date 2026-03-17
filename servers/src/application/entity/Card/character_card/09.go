package character_card

type Card09 struct {
	CharacterCardTemplate
}

func NewCard09() *Card09 {
	return &Card09{}
}

func (c *Card09) Attack() {

}
func (c *Card09) Hurt() {
}

func (c *Card09) GetID() int {
	return 9
}
