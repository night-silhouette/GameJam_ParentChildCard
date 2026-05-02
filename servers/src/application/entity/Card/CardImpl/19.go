package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card19 struct {
	BaseCard
}

func NewCard19() *Card19 {
	return &Card19{}
}

func (c *Card19) Attack(tempId int) {

}
func (c *Card19) Hurt(tempId int, HurtHp int) {
}

func (c *Card19) GetID() int {
	return 19
}

func (c *Card19) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card19) Skill(tempId int) {

}

func (c *Card19) Death(tempId int) {

}
