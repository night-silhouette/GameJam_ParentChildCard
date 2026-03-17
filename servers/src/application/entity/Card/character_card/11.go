package character_card

type Card11 struct {
	CharacterCardTemplate
}

func NewCard11() *Card11 {
	return &Card11{}
}

func (c *Card11) Attack() {

}
func (c *Card11) Hurt() {
}

func (c *Card11) GetID() int {
	return 11
}
