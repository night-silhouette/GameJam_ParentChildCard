package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
	"pcc_card/global"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"time"
)

type CharacterBaseCard struct {
	BaseCard
	Card CardAbstract.Card
}

func (c *CharacterBaseCard) Attack(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId) //检查指向对象是否有守护
	return c.ShareAttack(FinalId)
}

func (c *CharacterBaseCard) Retaliate(AttackId int, HurtHp float64) {
	RetaliateTotal := 0.0
	for _, buff := range c.BuffList {
		if buff.BuffId == protocol.Retaliate {
			RetaliateTotal += buff.Value
		}
	}
	if RetaliateTotal != 0.0 {
		c.Notify(BattleData.AnRetaliate, -1, AttackId, c.GetTempId())
		c.EffectAttack(AttackId, float64(int(HurtHp*RetaliateTotal)), BattleData.TrueDamage)
	}
}

func (c *CharacterBaseCard) NoSourceHurt(HurtHp float64, category BattleData.ValueChange) {
	c.EffectHurt(-1, HurtHp, category)
}

func (c *CharacterBaseCard) Hurt(AttackId int, HurtHp float64, category BattleData.ValueChange) {
	//-------------反击逻辑------------
	c.Retaliate(AttackId, HurtHp)
	//-------------反击逻辑------------

	c.Notify(BattleData.AnHurt, -1, AttackId, c.GetTempId())
	c.EffectHurt(AttackId, HurtHp, category)
}

// 父类的skill函数,消耗了能量,通知前端,true表示,能量已经扣了
func (c *CharacterBaseCard) Skill(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId)
	return c.ShareSkill(FinalId)
}

// 如果被无主伤害杀死,那杀死者的id为-1
func (c *CharacterBaseCard) Death(AttackId int) {
	if c.BtCtx.GetWeather() == protocol.Fengdu && !c.changeJiangShi { //如果天气是这个,就变僵尸
		c.ChangeForm(BattleData.JiangShi)
		c.changeJiangShi = true
		return //如果变僵尸了,就不用死了
	}
	c.Notify(BattleData.AnDeath, -1, AttackId, c.GetTempId())
	CheckIsInterrupt := true //声明这个指针bool,接受等待弃牌的结果,弃牌完,如果没有出战牌的话,就要中断,默认值随便搞的
	SelectCharacterCard := make([]int, 0)
	c.SetCardBt(&SelectCharacterCard, &CheckIsInterrupt)
	c.Interrupt(&SelectCharacterCard, global.Interrupt*time.Second, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 1, &CheckIsInterrupt, BattleData.Deploy)
	c.DisCard(&[]int{c.GetTempId()}, &CheckIsInterrupt) //反向压入
}

func (c *CharacterBaseCard) BtCry() {

}

func (c *CharacterBaseCard) RoundEnd() {

}

func (c *CharacterBaseCard) NextRound() {

}

//todo

// 把卡设置成新的form的状态,包含通知
func (c *CharacterBaseCard) ChangeForm(Form BattleData.Form) bool {
	v := BattleData.FormValuesMap[Form]
	c.SetHpNow(v.Hp)
	c.SetAtkNow(v.Damage) //取出value

	//-------------------通知form数值变化-------------------
	for _, UserId := range c.BtCtx.GetIds() {
		data := BattleData.NewFormChangeDto(Form, c.GetTempId(), *c.BtCtx.GetDataAll(UserId))
		c.BtCtx.ProtoSendAction(UserId, BattleDto.NewAction(BattleDto.FormChange, BattleDto.Result, data))
	}
	return true
	//-------------------通知form数值变化-------------------
}

func (c *CharacterBaseCard) CheckIsHaveBuff(BuffId protocol.BuffId) bool {
	for _, b := range c.BuffList {
		if b.BuffId == BuffId {
			return true
		}
	}
	return false

}

func (c *CharacterBaseCard) CheckGuard(TargetId int) int {
	if TargetId == c.GetTempId() {
		return TargetId
	}
	var FinalId int
	IsGuard, Guard := c.BtCtx.CheckBuff(TargetId, protocol.Guard) //检查指向对象是否有守护
	if IsGuard {
		FinalId = int(Guard.Value)
	} else {
		FinalId = TargetId
	}
	return FinalId
}

func (c *CharacterBaseCard) ShareAttack(TargetId int) bool {
	offset := int(c.GetInfo()["skillCharge"].(float64))
	if !c.BtCtx.ProtoColCanUpdateEnergy(c.OwnerId, -offset) {
		return false
	}
	if c.CheckIsHaveBuff(protocol.Binding) { //能量吞掉再滚出去
		c.EffectUpdateEnergy(-offset)
		return false
	}

	c.Notify(BattleData.AnAttack, -1, c.GetTempId(), TargetId)
	c.EffectAttack(TargetId, c.AtkNow, BattleData.Damage)
	c.EffectUpdateEnergy(-offset) //反向压入,先扣能量,再伤害
	return true
}

func (c *CharacterBaseCard) ShareSkill(TargetId int) bool {
	offset := int(c.GetInfo()["skillCharge"].(float64))
	if !c.BtCtx.ProtoColCanUpdateEnergy(c.OwnerId, -offset) {
		return false
	}
	if c.CheckIsHaveBuff(protocol.Binding) {
		c.EffectUpdateEnergy(-offset)
		return false
	}
	c.Notify(BattleData.AnSkill, -1, c.GetTempId(), TargetId)
	c.EffectUpdateEnergy(-offset)
	return true
}
