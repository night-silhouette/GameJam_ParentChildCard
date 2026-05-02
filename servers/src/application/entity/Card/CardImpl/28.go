package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card28 struct {
	BaseCard
}

func NewCard28() *Card28 {
	return &Card28{}
}

func (c *Card28) Attack(tempId int) {

}
func (c *Card28) Hurt(tempId int, HurtHp int) {
}

func (c *Card28) GetID() int {
	return 28
}

func (c *Card28) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card28) Skill(tempId int) {

}

func (c *Card28) Death(tempId int) {

}
