package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
)

type Card1002 struct {
	BaseCard
}

func NewCard1002() *Card1002 {
	return &Card1002{}
}

func (c *Card1002) PlayMagic() {

}

func (c *Card1002) GetID() int {
	return 1002
}

func (c *Card1002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
