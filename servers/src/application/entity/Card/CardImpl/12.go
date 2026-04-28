package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card12 struct {
	BaseCard
}

func NewCard12() *Card12 {
	return &Card12{}
}

func (c *Card12) Attack(w BattleData.Where) {

}
func (c *Card12) Hurt() {
}

func (c *Card12) GetID() int {
	return 12
}

func (c *Card12) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card12) Skill() {

}
