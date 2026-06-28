package CardAbstract

import "pcc_card/application/entity/BattleData"

type CtxRecord struct {
	LastHurt *BattleData.CardCalcValueDto //上一次收到伤害
}

func NewCtxRecord() *CtxRecord {
	res := CtxRecord{}
	return &res
}
