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
	pc.ProtoNotifyValue()
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

type Heal struct {
	UserId       int
	TargetTempId *int
	HealValue    float64
}

func (H *Heal) Execute(pc ProtocolCardWithCtx) {
	var target int
	if H.TargetTempId != nil {
		target = *H.TargetTempId
	}
	pc.ProtoColHealCardBt(H.UserId, target, H.HealValue)
	pc.ProtoNotifyValue()
}

func NewHeal(UserId int, TargetTempId *int, HealValue float64) *Heal {
	res := Heal{}
	res.UserId = UserId
	res.TargetTempId = TargetTempId
	res.HealValue = HealValue
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
		pc.ProtoNotifyCardMove()
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

//----------------------------------------------------

type UpdateEnergy struct {
	UserId int
	Offset int
}

func (U *UpdateEnergy) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColUpdateEnergy(U.UserId, U.Offset)
}

func NewUpdateEnergy(UserId int, Offset int) *UpdateEnergy {
	res := UpdateEnergy{}
	res.UserId = UserId
	res.Offset = Offset
	return &res
}

//----------------------------------------------------

type Custom struct {
	ExecFunc func(pc ProtocolCardWithCtx)
}

func (c *Custom) Execute(pc ProtocolCardWithCtx) {
	c.ExecFunc(pc)
}

func NewCustom(ExecFunc func(pc ProtocolCardWithCtx)) *Custom {
	res := Custom{}
	res.ExecFunc = ExecFunc
	return &res
}
