package CardImpl

import (
	_ "embed"
	"pcc_card/application/entity/BattleData"
)

type BaseCard struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}

func (c *BaseCard) GetCardDto() BattleData.CardDto {
	res := BattleData.CardDto{}
	res.Id = c.ID
	if c.Info != nil {
		if c.Info["hp"] != nil {
			res.Hp = c.Info["hp"].(int)
		}
		if c.Info["damage"] != nil {
			res.Damage = c.Info["damage"].(int)
		}
	}
	return res
}

func (c *BaseCard) GetID() int {
	return -1
}

func (c *BaseCard) SetInfo(info map[string]any) {
	c.Info = info
}
func (c *BaseCard) GetInfo() map[string]any {
	return c.Info
}
