package character_card

type Card06 struct {
	CharacterCardTemplate
}

func NewCard06() *Card06 {
	return &Card06{}
}

func (c *Card06) Attack() {

}
func (c *Card06) Hurt() {
}

func (c *Card06) GetID() int {
	return 6
}
