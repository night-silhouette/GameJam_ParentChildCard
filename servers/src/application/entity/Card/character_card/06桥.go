package character_card

import (
	_ "embed"
	"pcc_card/application/service/battleservice"
)

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

func (c *Card06) CreateCardData{
	c.ID=6
	c.Info["hp"]=2
	c.Info["damage"]=0
}