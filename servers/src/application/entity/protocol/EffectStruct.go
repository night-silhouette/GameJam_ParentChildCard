package protocol

import (
	"pcc_card/application/entity/BattleData"
	"time"
)

type Attack struct {
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
}

func NewAttack(UserId int, SendTempId int, TargetTempId int, AtkValue float64) *Attack {
	res := Attack{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	return &res
}

func (A *Attack) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColCardBtAttack(A.SendTempId, A.UserId, A.TargetTempId, A.AtkValue)
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type Hurt struct {
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
}

func (A *Hurt) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColReduceCardBtHp(A.SendTempId, A.UserId, A.TargetTempId, A.AtkValue)
}

func NewHurt(UserId int, SendTempId int, TargetTempId int, AtkValue float64) *Hurt {
	res := Hurt{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	return &res
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type Interrupt struct {
	UserId     int
	Time       time.Duration
	TempIdList []int
	SelectNum  int
	Res        chan []int
}

func (I *Interrupt) Execute(pc ProtocolCardWithCtx) {
	dto := BattleData.NewInterruptDto(I.Time, I.TempIdList, I.SelectNum)
	pc.ProtoColInterrupt(I.UserId, dto, I.Res, I.Time)
}

//----------------------------------------------------

type DisCard struct {
	UserId     int
	TempIdList *[]int
}

func (D *DisCard) Execute(pc ProtocolCardWithCtx) {
	for _, tempId := range *D.TempIdList {
		pc.ProtoColMoveDisCardPool(D.UserId, tempId)
		pc.Notify(BattleData.NewAnimationDto(tempId, BattleData.AnDisCard, pc.GetBtCardInfo(D.UserId)), D.UserId)

	}
}

func NewDisCard(UserId int, TempIdList *[]int) *DisCard {
	res := DisCard{}
	res.UserId = UserId
	res.TempIdList = TempIdList
	return &res
}

//----------------------------------------------------

type SetCardBt struct {
	TempIdList *[]int
	UserId     int
}

func (S *SetCardBt) Execute(pc ProtocolCardWithCtx) {
	tempId := (*S.TempIdList)[0]
	pc.ProtoColSetCardBt(S.UserId, tempId)

}

// NewSetCardBt 只上数组里的一张
func NewSetCardBt(UserId int, TempIdList *[]int) *SetCardBt {
	res := SetCardBt{}
	res.UserId = UserId
	res.TempIdList = TempIdList

	return &res
}

//----------------------------------------------------

type GiveBuff struct {
	BuffListP *[]Buff
	TempId    int
	Buff      Buff
}

func (G *GiveBuff) Execute(pc ProtocolCardWithCtx) {
	*G.BuffListP = append(*G.BuffListP, G.Buff)
}

func NewGiveBuff(TempId int, Buff Buff, BuffListP *[]Buff) *GiveBuff {
	res := GiveBuff{}
	res.TempId = TempId
	res.Buff = Buff
	res.BuffListP = BuffListP
	return &res
}
