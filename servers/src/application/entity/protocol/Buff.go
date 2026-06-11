package protocol

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
	Vulnerability
	Block

	HealingBoost
	HealingDecay

	Wither
	Binding
	Retaliate
	Confine
	Giant
	Disarm
)

var BuffRoundEndFuncMap map[BuffId]func(pc ProtocolCardWithCtx, value float64)

func InitBuffRoundEndFuncMap() {
	BuffRoundEndFuncMap = make(map[BuffId]func(pc ProtocolCardWithCtx, value float64))

	BuffRoundEndFuncMap[BonusDamage] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Powerful] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Weakness] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[DamageImmunity] = func(pc ProtocolCardWithCtx, value float64) {

	}
	BuffRoundEndFuncMap[Vulnerability] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Block] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[HealingBoost] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[HealingDecay] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Wither] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Binding] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Retaliate] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Confine] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Giant] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffRoundEndFuncMap[Disarm] = func(pc ProtocolCardWithCtx, value float64) {}
}

var BuffOnApplyFuncMap map[BuffId]func(pc ProtocolCardWithCtx, value float64)

func InitBuffOnApplyFuncMap() {
	BuffOnApplyFuncMap = make(map[BuffId]func(pc ProtocolCardWithCtx, value float64))

	BuffOnApplyFuncMap[BonusDamage] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Powerful] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Weakness] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[DamageImmunity] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Vulnerability] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Block] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[HealingBoost] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[HealingDecay] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Wither] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Binding] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Retaliate] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Confine] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Giant] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnApplyFuncMap[Disarm] = func(pc ProtocolCardWithCtx, value float64) {}
}

var BuffOnRemoveFuncMap map[BuffId]func(pc ProtocolCardWithCtx, value float64)

func InitBuffOnRemoveFuncMap() {
	BuffOnRemoveFuncMap = make(map[BuffId]func(pc ProtocolCardWithCtx, value float64))

	BuffOnRemoveFuncMap[BonusDamage] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Powerful] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Weakness] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[DamageImmunity] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Vulnerability] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Block] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[HealingBoost] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[HealingDecay] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Wither] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Binding] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Retaliate] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Confine] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Giant] = func(pc ProtocolCardWithCtx, value float64) {}
	BuffOnRemoveFuncMap[Disarm] = func(pc ProtocolCardWithCtx, value float64) {}
}
