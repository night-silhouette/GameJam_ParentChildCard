package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
)

type Card06 struct {
	BaseCard
}

func NewCard06() *Card06 {
	return &Card06{}
}

func (c *Card06) Attack(tempId int) {

}
func (c *Card06) Hurt(tempId int, HurtHp float64) {
}

func (c *Card06) GetID() int {
	return 6
}

func (c *Card06) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card06) Skill(tempId int) {

}

func (c *Card06) Death(tempId int) {

}
