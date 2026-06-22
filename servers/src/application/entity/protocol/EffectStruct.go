package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"time"
)

type Attack struct {
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
	Dec          *CardMeta.Decorator
	Category     BattleData.ValueChange
}

func NewAttack(UserId int, SendTempId int, TargetTempId int, AtkValue float64, Dec *CardMeta.Decorator, category BattleData.ValueChange) *Attack {
	res := Attack{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	res.Dec = Dec
	res.Category = category
	return &res
}

func (A *Attack) Execute(pc ProtocolCardWithCtx) {
	originValue := A.AtkValue
	FinalAtkValue := A.Dec.CalcAttack(originValue)
	pc.ProtoColCardBtAttack(A.SendTempId, A.UserId, A.TargetTempId, float64(FinalAtkValue), A.Category)
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type Hurt struct {
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
	Dec          *CardMeta.Decorator //这已经是被buff增益过的了
	Category     BattleData.ValueChange
}

func (A *Hurt) Execute(pc ProtocolCardWithCtx) {
	var FinalAtkValue int
	IfMiss := false
	if A.Category == BattleData.Damage { //判断是否是真伤,是的话,就不走装饰器
		FinalAtkValue, IfMiss = A.Dec.CalcHurt(A.AtkValue)
	} else if A.Category == BattleData.TrueDamage {
		FinalAtkValue = int(A.AtkValue)
	}
	pc.ProtoColReduceCardBtHp(A.SendTempId, A.TargetTempId, float64(FinalAtkValue))
	pc.ProtoNotifyValue(A.Category, -float64(FinalAtkValue), A.TargetTempId, IfMiss)
	if FinalAtkValue > 0 {
		pc.Broad(CardMeta.NewBroadInfo(CardMeta.Wound, A.SendTempId, A.TargetTempId)) //广播被击伤
	}
}

func NewHurt(UserId int, SendTempId int, TargetTempId int, AtkValue float64, Dec *CardMeta.Decorator, Category BattleData.ValueChange) *Hurt {
	res := Hurt{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	res.Dec = Dec
	res.Category = Category
	return &res
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type Heal struct {
	TargetTempId *int
	HealValue    float64
	Dec          *CardMeta.Decorator
}

func (H *Heal) Execute(pc ProtocolCardWithCtx) {
	var target int
	if H.TargetTempId != nil {
		target = *H.TargetTempId
	}
	originValue := H.HealValue
	FinalHeal := H.Dec.CalcHeal(originValue)
	pc.ProtoColHealCardBt(target, float64(FinalHeal))
	pc.ProtoNotifyValue(BattleData.Heal, H.HealValue, *H.TargetTempId, false)
}

func NewHeal(TargetTempId *int, HealValue float64, Dec *CardMeta.Decorator) *Heal {
	res := Heal{}

	res.TargetTempId = TargetTempId
	res.HealValue = HealValue
	res.Dec = Dec
	return &res
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type Interrupt struct {
	UserId           int
	Time             time.Duration
	TempIdList       []int
	SelectNum        int
	Res              chan []int
	CheckIsInterrupt *bool
	CallTempId       int
	InterruptType    BattleData.InterruptType
}

func (I *Interrupt) Execute(pc ProtocolCardWithCtx) {
	if !(*I.CheckIsInterrupt) {
		return
	}
	dto := BattleData.NewInterruptDto(I.Time, I.TempIdList, I.SelectNum, I.InterruptType, I.CallTempId)
	pc.ProtoColInterrupt(I.UserId, dto, I.Res, I.Time)
}

//----------------------------------------------------

type DisCard struct {
	UserId      int
	TempIdList  *[]int
	IsInterrupt *bool
}

func (D *DisCard) Execute(pc ProtocolCardWithCtx) {
	for _, tempId := range *D.TempIdList {
		pc.ProtoColMoveDisCardPool(D.UserId, tempId)
		pc.ProtoNotifyCardMove(BattleData.DisCardPool, tempId)
	}
	if !pc.CheckCard(D.UserId) { //检查出没有出战的牌了

		*D.IsInterrupt = true
	} else {
		*D.IsInterrupt = false
	}
}

func NewDisCard(UserId int, TempIdList *[]int, IsInterrupt *bool) *DisCard {
	res := DisCard{}
	res.UserId = UserId
	res.TempIdList = TempIdList
	res.IsInterrupt = IsInterrupt
	return &res
}

//----------------------------------------------------

type SetCardBt struct {
	TempIdList       *[]int
	UserId           int
	CheckIsInterrupt *bool
}

func (S *SetCardBt) Execute(pc ProtocolCardWithCtx) {
	if !(*S.CheckIsInterrupt) {
		return
	}
	tempId := (*S.TempIdList)[0]
	pc.ProtoColSetCardBt(S.UserId, tempId)

}

// NewSetCardBt 只上数组里的一张
func NewSetCardBt(UserId int, TempIdList *[]int, CheckIsInterrupt *bool) *SetCardBt {
	res := SetCardBt{}
	res.UserId = UserId
	res.TempIdList = TempIdList
	res.CheckIsInterrupt = CheckIsInterrupt
	return &res
}

//----------------------------------------------------

type GiveBuff struct {
	TempId *int
	Buff   Buff
}

func (G *GiveBuff) Execute(pc ProtocolCardWithCtx) {
	pc.GiveBuff(*G.TempId, &G.Buff)
	for _, UserId := range pc.GetIds() { //做buff改动的通知
		pc.ProtoSendAction(UserId, BattleDto.NewAction(BattleDto.BuffChange, BattleDto.Result, map[string]any{
			"data_all": pc.GetDataAll(UserId),
		}))
	}

}

func NewGiveBuff(TempId *int, Buff Buff) *GiveBuff {
	res := GiveBuff{}
	res.TempId = TempId
	res.Buff = Buff
	return &res
}

//----------------------------------------------------

type UpdateEnergy struct {
	UserId int
	Offset int
}

func (U *UpdateEnergy) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColUpdateEnergy(U.UserId, U.Offset)
	for _, UserId := range pc.GetIds() {
		pc.ProtoSendAction(UserId, BattleDto.NewAction(BattleDto.EnergyChange, BattleDto.Result, map[string]interface{}{
			"data_all": pc.GetDataAll(UserId),
		}))
	}
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

//----------------------------------------

type ChangeMaxHp struct {
	TargetTempId int
	MaxHp        float64
}

func (c *ChangeMaxHp) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColSetMaxHp(c.TargetTempId, c.MaxHp)
}

func NewChangeMaxHp(TargetTempId int, MaxHp float64) *ChangeMaxHp {
	res := ChangeMaxHp{}
	res.MaxHp = MaxHp
	res.TargetTempId = TargetTempId
	return &res
}
