package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card38 struct {
	BaseCard
}

func NewCard38() *Card38 {
	return &Card38{}
}

func (c *Card38) Attack(tempId int) {

}
func (c *Card38) Hurt(tempId int, HurtHp float64) {
}

func (c *Card38) GetID() int {
	return 38
}

func (c *Card38) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card38) Skill(tempId int) {

}

func (c *Card38) Death(tempId int) {

}
