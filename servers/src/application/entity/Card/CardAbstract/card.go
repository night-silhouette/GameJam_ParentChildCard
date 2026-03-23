package CardAbstract

import "pcc_card/application/entity/BattleData"

type Card interface {
	GetID() int
	SetInfo(info map[string]any)
	GetInfo() map[string]any
	Clone() Card
	GetCardDto() BattleData.CardDto
}
