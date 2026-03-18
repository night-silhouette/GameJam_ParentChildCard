package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card33 struct {
	BaseCard
}

func NewCard33() *Card33 {
	return &Card33{}
}

func (c *Card33) GetID() int {
	return 33
}

func (c *Card33) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
