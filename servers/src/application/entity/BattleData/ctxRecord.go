package BattleData

type HurtRecord struct {
	CardCalcValueDto
	TempId int
}

func NewHurtRecord(dto CardCalcValueDto, TempId int) *HurtRecord {
	return &HurtRecord{
		CardCalcValueDto: dto,
		TempId:           TempId,
	}
}

type CtxRecord struct {
	LastHurt *HurtRecord //上一次收到伤害//value是正的,表示扣了多少血
}

func NewCtxRecord() *CtxRecord {
	res := CtxRecord{}
	return &res
}
