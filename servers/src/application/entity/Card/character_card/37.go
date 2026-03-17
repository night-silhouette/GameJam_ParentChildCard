package character_card

type Card37 struct {
	CharacterCardTemplate
}

func NewCard37() *Card37 {
	return &Card37{}
}

func (c *Card37) Attack() {

}
func (c *Card37) Hurt() {
}

func (c *Card37) GetID() int {
	return 37
}
