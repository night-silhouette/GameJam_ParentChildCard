package character_card

type Card22 struct {
	CharacterCardTemplate
}

func NewCard22() *Card22 {
	return &Card22{}
}

func (c *Card22) Attack() {

}
func (c *Card22) Hurt() {
}

func (c *Card22) GetID() int {
	return 22
}
