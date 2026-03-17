package character_card

type Card struct {
	CharacterCardTemplate
}

func NewCard() *Card {
	return &Card{}
}

func (c *Card) Attack() {

}
func (c *Card) Hurt() {
}
