package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
)

type BuffBase struct {
	//几层
	Stacks int
	BuffId int
}

type Giant struct {
	BuffBase
}

func (b *Giant) GetBuffDto() BattleData.BuffDto {
	res := BattleData.BuffDto{}
	res.BuffId = b.BuffId
	res.BuffStacks = b.Stacks
	return res
}

func (b *Giant) BuffExecute(info *CardMeta.CardInfo, signalMap map[string]CardMeta.ControlSignal, dec *CardMeta.Decorator) {

}

func (b *Giant) RoundEnd() {

}

func NewGiant(Stacks int) *Giant {
	res := &Giant{}
	res.Stacks = Stacks
	res.BuffId = 0
	return res
}
