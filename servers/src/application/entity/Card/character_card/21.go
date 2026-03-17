package character_card

type Card21 struct {
	CharacterCardTemplate
}

func NewCard21() *Card21 {
	return &Card21{}
}

func (c *Card21) Attack() {

}
func (c *Card21) Hurt() {
}
