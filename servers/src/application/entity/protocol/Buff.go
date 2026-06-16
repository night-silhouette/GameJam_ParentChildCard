package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
)

func InitBuff() {
	InitBuffRoundEndFuncMap()
	InitBuffOnRemoveFuncMap()
	InitBuffOnApplyFuncMap()

}

type BuffId int

const (
	BonusDamage BuffId = iota
	Powerful
	Weakness

	DamageImmunity
	Vulnerability //免伤
	Block

	HealingBoost
	HealingDecay

	Wither
	Binding
	Retaliate
	Confine
	Giant
	Disarm
	XuFeng
)

var BuffRoundEndFuncMap map[BuffId]func(pc ProtocolCardWithCtx, value float64, card BuffNeed)

// ------------------回合结束---------------------
func InitBuffRoundEndFuncMap() {
	BuffRoundEndFuncMap = make(map[BuffId]func(pc ProtocolCardWithCtx, value float64, card BuffNeed))

	BuffRoundEndFuncMap[BonusDamage] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Powerful] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Weakness] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[DamageImmunity] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Vulnerability] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Block] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[HealingBoost] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[HealingDecay] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Wither] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {
		pc.ProtoColPush(NewCustom(func(pc ProtocolCardWithCtx) {
			pc.ProtoColAttackNoHurt(card.GetTempId(), int(value), BattleData.Damage)
		}))
	}
	BuffRoundEndFuncMap[Binding] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Retaliate] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Confine] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Giant] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[Disarm] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffRoundEndFuncMap[XuFeng] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
}

//------------------回合结束---------------------

var BuffOnApplyFuncMap map[BuffId]func(pc ProtocolCardWithCtx, value float64, card BuffNeed)

// ------------------buff上去的时候---------------------
func InitBuffOnApplyFuncMap() {
	BuffOnApplyFuncMap = make(map[BuffId]func(pc ProtocolCardWithCtx, value float64, card BuffNeed))

	BuffOnApplyFuncMap[BonusDamage] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Powerful] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Weakness] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[DamageImmunity] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Vulnerability] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {
	}
	BuffOnApplyFuncMap[Block] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[HealingBoost] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[HealingDecay] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Wither] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Binding] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Retaliate] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Confine] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Giant] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[Disarm] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnApplyFuncMap[XuFeng] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
}

//------------------buff上去的时候---------------------

var BuffOnRemoveFuncMap map[BuffId]func(pc ProtocolCardWithCtx, value float64, card BuffNeed)

// ------------------buff清除的时候---------------------
func InitBuffOnRemoveFuncMap() {
	BuffOnRemoveFuncMap = make(map[BuffId]func(pc ProtocolCardWithCtx, value float64, card BuffNeed))

	BuffOnRemoveFuncMap[BonusDamage] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Powerful] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Weakness] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[DamageImmunity] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Vulnerability] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Block] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[HealingBoost] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[HealingDecay] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Wither] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Binding] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Retaliate] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Confine] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Giant] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[Disarm] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
	BuffOnRemoveFuncMap[XuFeng] = func(pc ProtocolCardWithCtx, value float64, card BuffNeed) {}
}

//------------------buff清除的时候---------------------

type BuffNeed interface {
	SetDec(Dec *CardMeta.Decorator)
	GetDec() *CardMeta.Decorator
	GetTempId() int
}
