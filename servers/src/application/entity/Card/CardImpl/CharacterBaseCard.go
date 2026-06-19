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
	offset := int(c.GetInfo()["skillCharge"].(float64))
	if !c.BtCtx.ProtoColCanUpdateEnergy(c.OwnerId, -offset) {
		return false
	}
	if c.CheckIsHaveBuff(protocol.Binding) {
		return false
	}

	c.Notify(BattleData.AnAttack, -1, c.GetTempId(), TargetId)

	c.EffectAttack(TargetId, c.AtkNow, BattleData.Damage)
	c.EffectUpdateEnergy(-offset) //反向压入,先扣能量,再伤害
	return true
}

func (c *CharacterBaseCard) Hurt(AttackId int, HurtHp float64, category BattleData.ValueChange) {
	c.Notify(BattleData.AnHurt, -1, AttackId, c.GetTempId())
	c.EffectHurt(AttackId, HurtHp, category)
}

// 父类的skill函数,消耗了能量,通知前端,true表示,能量已经扣了
func (c *CharacterBaseCard) Skill(TargetId int) bool {
	offset := int(c.GetInfo()["skillCharge"].(float64))
	if !c.BtCtx.ProtoColCanUpdateEnergy(c.OwnerId, -offset) {
		return false
	}
	if c.CheckIsHaveBuff(protocol.Binding) {
		return false
	}
	c.Notify(BattleData.AnSkill, -1, c.GetTempId(), TargetId)
	c.EffectUpdateEnergy(-offset)
	return true
}

func (c *CharacterBaseCard) Death(AttackId int) {

	if c.BtCtx.GetWeather() == protocol.Fengdu { //如果天气是这个,就变僵尸
		c.ChangeForm(BattleData.JiangShi)
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

//todo
//---------二次分装---------

func (c *CharacterBaseCard) EffectAttack(targetTempId int, AtkHp float64, category BattleData.ValueChange) {
	c.BtCtx.ProtoColPush(protocol.NewAttack(c.OwnerId, c.TempId, targetTempId, AtkHp, c.GetDec(), category))
}
func (c *CharacterBaseCard) EffectHurt(AttackId int, AtkHp float64, Category BattleData.ValueChange) {
	c.BtCtx.ProtoColPush(protocol.NewHurt(c.OwnerId, AttackId, c.TempId, AtkHp, c.GetDec(), Category))
}

func (c *CharacterBaseCard) EffectHeal(targetTempId int, HealHp float64) {
	c.BtCtx.ProtoColPush(protocol.NewHeal(&targetTempId, HealHp, c.GetDec()))
}

func (c *CharacterBaseCard) EffectUpdateEnergy(offset int) {
	c.BtCtx.ProtoColPush(protocol.NewUpdateEnergy(c.OwnerId, offset))
}

// 通知标明行为发起者和受到者
func (c *CharacterBaseCard) Notify(Beh BattleData.AnimationBehavior, UserID int, CallerId int, AcceptorId int) { //都是tempid
	c.BtCtx.Notify(BattleData.NewAnimationDto(CallerId, AcceptorId, Beh), UserID)
}
func (c *CharacterBaseCard) Interrupt(res *[]int, time time.Duration, TempIdList []int, SelectNum int, CheckIsInterrupt *bool, interruptType BattleData.InterruptType) { //res一定要塞到effect函数里处理
	resChan := make(chan []int)
	c.BtCtx.ProtoColPush(&protocol.Interrupt{
		UserId:           c.OwnerId,
		Time:             time,
		TempIdList:       TempIdList,
		SelectNum:        SelectNum,
		Res:              resChan,
		CheckIsInterrupt: CheckIsInterrupt,
		CallTempId:       c.TempId,
		InterruptType:    interruptType,
	})
	go func() {
		val := <-resChan
		*res = val
		c.BtCtx.ProtoColCancelInterrupt()
	}()
}

func (c *CharacterBaseCard) DisCard(TempIdList *[]int, IsInterrupt *bool) {
	c.BtCtx.ProtoColPush(protocol.NewDisCard(c.OwnerId, TempIdList, IsInterrupt))
}

func (c *CharacterBaseCard) SetCardBt(TempIdList *[]int, IsInterrupt *bool) {
	c.BtCtx.ProtoColPush(protocol.NewSetCardBt(c.OwnerId, TempIdList, IsInterrupt))
}

func (c *CharacterBaseCard) GiveBuff(TempId *int, b protocol.Buff) {
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(TempId, b))
}

func (c *CharacterBaseCard) ReMoveBuffByTempId(BuffTempId int) {
	if c.BuffList == nil || len(c.BuffList) == 0 {
		return
	}
	oldList := c.BuffList
	newList := make([]*protocol.Buff, 0, len(oldList))
	for _, buff := range oldList {
		if buff == nil {
			continue
		}
		if buff.TempId == BuffTempId {
			// 如果你需要在 Buff 销毁时通知前端播放特效（比如消失动画），可以在这里做：
			continue // 跳过它，不把它装进新切片，达到删除的效果
		}

		// 没匹配上的正常保留
		newList = append(newList, buff)
	}

	// 3.  将过滤后的全新切片指针重新赋给卡牌
	c.BuffList = newList
}

func (c *CharacterBaseCard) NewCustom(ExecFunc func(pc protocol.ProtocolCardWithCtx)) {
	c.BtCtx.ProtoColPush(protocol.NewCustom(ExecFunc))
}

func (c *CharacterBaseCard) ChangeMaxHp(TargetTempId int, MaxHp float64) {
	c.BtCtx.ProtoColPush(protocol.NewChangeMaxHp(TargetTempId, MaxHp))
}

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
