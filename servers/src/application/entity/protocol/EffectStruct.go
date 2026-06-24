package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
)

type Attack struct {
	EffectBase
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
	Dec          *CardMeta.Decorator
	Category     BattleData.ValueChange
}

// 中断的话,TargetTempId传什么都可以
func NewAttack(UserId int, SendTempId int, TargetTempId int, AtkValue float64, Dec *CardMeta.Decorator, category BattleData.ValueChange, isInterrupt bool, config *InterruptConfig) *Attack {
	res := Attack{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	res.Dec = Dec
	res.Category = category
	res.CheckIsInterrupt = isInterrupt
	res.InterruptConfig = config
	return &res
}

func (A *Attack) Execute(pc ProtocolCardWithCtx) {
	if !A.CheckIsInterrupt {
		A.ShareLogic(pc, A.TargetTempId)
	} else {
		A.CallInterrupt(pc, func(res []int, pc ProtocolCardWithCtx) {
			for _, ObjId := range res {
				A.ShareLogic(pc, ObjId)
			}
		})
	}

}
func (A *Attack) ShareLogic(pc ProtocolCardWithCtx, Obj int) {
	originValue := A.AtkValue
	FinalAtkValue := A.Dec.CalcAttack(originValue)
	pc.ProtoColCardBtAttack(A.SendTempId, A.UserId, Obj, float64(FinalAtkValue), A.Category)
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type Hurt struct {
	EffectBase
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
	EffectBase
	TargetTempId *int
	HealValue    float64
	Dec          *CardMeta.Decorator
}

func (H *Heal) Execute(pc ProtocolCardWithCtx) {
	if !H.CheckIsInterrupt {
		H.ShareLogic(pc, *H.TargetTempId)
	} else {
		H.CallInterrupt(pc, func(res []int, pc ProtocolCardWithCtx) {
			for _, ObjId := range res {
				H.ShareLogic(pc, ObjId)
			}
		})
	}

}

func (H *Heal) ShareLogic(pc ProtocolCardWithCtx, ObjId int) {
	originValue := H.HealValue
	FinalHeal := H.Dec.CalcHeal(originValue)
	pc.ProtoColHealCardBt(ObjId, float64(FinalHeal))
	pc.ProtoNotifyValue(BattleData.Heal, H.HealValue, ObjId, false)
}

func NewHeal(TargetTempId *int, HealValue float64, Dec *CardMeta.Decorator, isInterrupt bool, config *InterruptConfig) *Heal {
	res := Heal{}
	res.TargetTempId = TargetTempId
	res.HealValue = HealValue
	res.Dec = Dec
	res.CheckIsInterrupt = isInterrupt
	res.InterruptConfig = config
	return &res
}

//-----------------------------------------------------------------------------------------------------------------------------------------

type DisCard struct {
	EffectBase
	UserId     int
	TempIdList []int
}

func (D *DisCard) Execute(pc ProtocolCardWithCtx) {
	if !D.CheckIsInterrupt {
		D.ShareLogic(pc, D.TempIdList)
	} else {
		D.CallInterrupt(pc, func(res []int, pc ProtocolCardWithCtx) {
			D.ShareLogic(pc, res)
		})
	}

}
func (D *DisCard) ShareLogic(pc ProtocolCardWithCtx, ObjId []int) {
	for _, tempId := range ObjId {
		pc.ProtoColMoveDisCardPool(D.UserId, tempId)
		pc.ProtoNotifyCardMove(BattleData.DisCardPool, tempId)
	}
}

func NewDisCard(UserId int, TempIdList []int, isInterrupt bool, config *InterruptConfig) *DisCard {
	res := DisCard{}
	res.UserId = UserId
	res.TempIdList = TempIdList
	res.CheckIsInterrupt = isInterrupt
	res.InterruptConfig = config
	return &res
}

//----------------------------------------------------

type SetCardBt struct {
	EffectBase
	TargetId int
	UserId   int
}

func (S *SetCardBt) Execute(pc ProtocolCardWithCtx) {
	if !S.CheckIsInterrupt {
		S.ShareLogic(pc, S.TargetId)
	} else {
		S.CallInterrupt(pc, func(res []int, pc ProtocolCardWithCtx) {
			for _, ObjId := range res {
				S.ShareLogic(pc, ObjId)
			}
		})
	}

}

func (S *SetCardBt) ShareLogic(pc ProtocolCardWithCtx, ObjId int) {
	pc.ProtoColSetCardBt(S.UserId, ObjId)
}

// NewSetCardBt 只上数组里的一张
func NewSetCardBt(UserId int, TargetId int, CheckIsInterrupt bool, config *InterruptConfig) *SetCardBt {
	res := SetCardBt{}
	res.UserId = UserId
	res.TargetId = TargetId
	res.CheckIsInterrupt = CheckIsInterrupt
	res.InterruptConfig = config
	return &res
}

//----------------------------------------------------

type GiveBuff struct {
	EffectBase
	TempId *int
	Buff   Buff
}

func (G *GiveBuff) Execute(pc ProtocolCardWithCtx) {
	if !G.CheckIsInterrupt {
		G.ShareLogic(pc, *G.TempId)
	} else {
		G.CallInterrupt(pc, func(res []int, pc ProtocolCardWithCtx) {
			for _, ObjId := range res {
				G.ShareLogic(pc, ObjId)
			}
		})
	}
}

func (G *GiveBuff) ShareLogic(pc ProtocolCardWithCtx, ObjId int) {
	pc.GiveBuff(ObjId, &G.Buff)
	for _, UserId := range pc.GetIds() { //做buff改动的通知
		pc.ProtoSendAction(UserId, BattleDto.NewAction(BattleDto.BuffChange, BattleDto.Result, map[string]any{
			"data_all": pc.GetDataAll(UserId),
		}))
	}
}

func NewGiveBuff(TempId *int, Buff Buff, isInterrupt bool, config *InterruptConfig) *GiveBuff {
	res := GiveBuff{}
	res.TempId = TempId
	res.Buff = Buff
	res.InterruptConfig = config
	res.CheckIsInterrupt = isInterrupt
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
