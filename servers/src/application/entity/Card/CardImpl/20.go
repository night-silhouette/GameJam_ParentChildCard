package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card20 struct {
	BaseCard
}

func NewCard20() *Card20 {
	return &Card20{}
}

func (c *Card20) Attack(tempId int) {

}
func (c *Card20) Hurt(tempId int) {
}

func (c *Card20) GetID() int {
	return 20
}

func (c *Card20) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card20) Skill(tempId int) {

}
