package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card17 struct {
	BaseCard
}

func NewCard17() *Card17 {
	return &Card17{}
}

func (c *Card17) Attack(w BattleData.Where) {

}
func (c *Card17) Hurt() {
}

func (c *Card17) GetID() int {
	return 17
}

func (c *Card17) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card17) Skill() {

}
