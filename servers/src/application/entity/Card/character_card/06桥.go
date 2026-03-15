package character_card

import (
	_ "embed"
	"encoding/json"
)

type Card06 struct {
	CharacterCardTemplate
}

//go:embed 06桥.json
var Data06 []byte

func NewCard06() *Card06 {
	Card := Card06{}
	json.Unmarshal(Data06, &Card)
	return &Card
}

func (c *Card06) Attack() {

}
func (c *Card06) Hurt() {

}
