package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
)

type Buff interface {
	BuffExecute(info *CardMeta.CardInfo, signalMap map[string]CardMeta.ControlSignal, dec *CardMeta.Decorator)
	GetBuffDto() BattleData.BuffDto
	RoundEnd()
}
