package character_card

type Card38 struct {
	CharacterCardTemplate
}

func NewCard38() *Card38 {
	return &Card38{}
}

func (c *Card38) Attack() {

}
func (c *Card38) Hurt() {
}

func (c *Card38) GetID() int {
	return 38
}
