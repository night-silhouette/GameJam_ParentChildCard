package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card08 struct {
	BaseCard
}

func NewCard08() *Card08 {
	return &Card08{}
}

func (c *Card08) Attack(w BattleData.Where) {

}
func (c *Card08) Hurt() {
}

func (c *Card08) GetID() int {
	return 8
}

func (c *Card08) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card08) Skill() {

}
