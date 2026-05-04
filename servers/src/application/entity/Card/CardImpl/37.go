package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card37 struct {
	CharacterBaseCard
}

func NewCard37() *Card37 {
	return &Card37{}
}

func (c *Card37) GetID() int {
	return 37
}

func (c *Card37) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

