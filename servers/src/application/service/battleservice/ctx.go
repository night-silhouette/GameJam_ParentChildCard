package battleservice

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"sync"
)

type Ctx struct {
	mu    sync.RWMutex
	IDA   int
	IDB   int
	DataA *PlayerData
	DataB *PlayerData
}

type PlayerData struct {
	ID           int
	CardInHand   *[]CardAbstract.ParentCard
	SkillCardUse *CardAbstract.SkillCard
	ParentCardBT *CardAbstract.ParentCard
	ChildCardBT  *CardAbstract.ChildCard
}
