package protocol

import (
	"pcc_card/application/entity/BattleData"
)

type Buff struct {
	TempId int //这是buff自己的id,而不是数据库上对应的code

	//几层
	Stacks int
	BuffId BuffId
	Value  float64 //数值
}

func (b *Buff) GetBuffDto() BattleData.BuffDto {
	res := BattleData.BuffDto{}
	res.BuffId = int(b.BuffId)
	res.BuffStacks = b.Stacks
	res.Value = b.Value
	return res
}

func NewBuffBase(buffId BuffId, buffStacks int, Value float64, TempId int) *Buff {
	return &Buff{
		BuffId: buffId,
		Stacks: buffStacks,
		Value:  Value,
		TempId: TempId,
	}
}
