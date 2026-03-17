package character_card

type Card29 struct {
	CharacterCardTemplate
}

func NewCard29() *Card29 {
	return &Card29{}
}

func (c *Card29) Attack() {

}
func (c *Card29) Hurt() {
}

func (c *Card29) GetID() int {
	return 29
}
