package CardAbstract

import "pcc_card/application/entity/BattleData"

type Card interface {
	GetID() int
	SetInfo(info map[string]any)
	GetInfo() map[string]any
	Clone() Card
}

func GetCardDto(c Card) BattleData.CardDto {
	res := BattleData.CardDto{}
	res.Id = c.GetID()
	info := c.GetInfo()
	if info != nil {
		if info["hp"] != nil {
			res.Hp = info["hp"].(float64)
		}
		if info["damage"] != nil {
			res.Damage = info["damage"].(float64)
		}
	}
	return res
}
