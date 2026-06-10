package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
)

type BuffBase struct {
	//几层
	Stacks int
	BuffId BuffId
	Value  float64 //数值
}

func (b *BuffBase) GetBuffDto() BattleData.BuffDto {
	res := BattleData.BuffDto{}
	res.BuffId = int(b.BuffId)
	res.BuffStacks = b.Stacks
	res.Value = b.Value
	return res
}

func NewBuffBase(buffId BuffId, buffStacks int, Value float64) *BuffBase {
	return &BuffBase{
		BuffId: buffId,
		Stacks: buffStacks,
		Value:  Value,
	}
}

func (b *BuffBase) RoundEnd(pc ProtocolCardWithCtx) {
	b.Stacks -= 1
	BuffRoundEndFuncMap[b.BuffId](pc, b.Value)
}

func AddBuff(card CardAbstract.Card, buff Buff) {
	dec := card.GetDec()
	card.AppendBuff(buff)

}
