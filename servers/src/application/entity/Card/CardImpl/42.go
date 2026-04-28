package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card42 struct {
	BaseCard
}

func NewCard42() *Card42 {
	return &Card42{}
}

func (c *Card42) Attack(w BattleData.Where) {

}
func (c *Card42) Hurt() {
}

func (c *Card42) GetID() int {
	return 42
}

func (c *Card42) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card42) Skill() {

}
