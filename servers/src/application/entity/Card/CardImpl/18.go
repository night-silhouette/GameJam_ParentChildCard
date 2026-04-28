package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card18 struct {
	BaseCard
}

func NewCard18() *Card18 {
	return &Card18{}
}

func (c *Card18) Attack(w BattleData.Where) {

}
func (c *Card18) Hurt() {
}

func (c *Card18) GetID() int {
	return 18
}

func (c *Card18) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card18) Skill() {

}
