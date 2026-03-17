package character_card

type Card43 struct {
	CharacterCardTemplate
}

func NewCard43() *Card43 {
	return &Card43{}
}

func (c *Card43) Attack() {

}
func (c *Card43) Hurt() {
}

func (c *Card43) GetID() int {
	return 43
}
