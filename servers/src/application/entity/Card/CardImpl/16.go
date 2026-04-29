package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card16 struct {
	BaseCard
}

func NewCard16() *Card16 {
	return &Card16{}
}

func (c *Card16) Attack(tempId int) {

}
func (c *Card16) Hurt(tempId int) {
}

func (c *Card16) GetID() int {
	return 16
}

func (c *Card16) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card16) Skill(tempId int) {

}
