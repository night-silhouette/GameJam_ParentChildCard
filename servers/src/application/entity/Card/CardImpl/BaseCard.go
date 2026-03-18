package CardImpl

import (
	_ "embed"
	"pcc_card/application/entity/Card/CardAbstract"
)

type BaseCard struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}

func (c *BaseCard) GetID() int {
	return -1
}

func (c *BaseCard) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *BaseCard) SetInfo(info map[string]any) {
	c.Info = info
}
func (c *BaseCard) GetInfo() map[string]any {
	return c.Info
}
