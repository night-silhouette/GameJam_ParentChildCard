package character_card

type Card16 struct {
	CharacterCardTemplate
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
