package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card40 struct {
	BaseCard
}

func NewCard40() *Card40 {
	return &Card40{}
}

func (c *Card40) Attack(tempId int) {

}
func (c *Card40) Hurt(tempId int, HurtHp int) {
}

func (c *Card40) GetID() int {
	return 40
}

func (c *Card40) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card40) Skill(tempId int) {

}

func (c *Card40) Death(tempId int) {

}
