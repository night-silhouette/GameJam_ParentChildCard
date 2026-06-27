package CardImpl

import (
	"context"
	_ "embed"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
)

type BaseCard struct {
	self CardAbstract.Card

	ID   int            `json:"id"`
	Info map[string]any `json:"-"`

	//动态变量
	CR                   *CardAbstract.CardRecord
	BtCtx                protocol.ProtocolCardWithCtx
	HpNow                float64
	AtkNow               float64
	TempId               int
	OwnerId              int
	BuffList             []*protocol.Buff
	Dec                  *CardMeta.Decorator
	SpecialCardStateChan chan *CardMeta.BroadInfo
	Form                 BattleData.Form
	changeJiangShi       bool
}

func (c *BaseCard) SetBtCtx(btCtx protocol.ProtocolCardWithCtx) {
	c.BtCtx = btCtx
}

func (c *BaseCard) GetHpNow() float64 {
	return c.HpNow
}
func (c *BaseCard) SetHpNow(hpNow float64) {
	OldHp := c.HpNow
	if OldHp <= hpNow {
		c.HpNow = hpNow
		return
	}
	c.HpNow = hpNow
	offset := OldHp - hpNow //这次扣掉的血
	c.CR.HurtedThisTurn += offset
}
func (c *BaseCard) GetAtkNow() float64 {
	return c.AtkNow
}
func (c *BaseCard) SetAtkNow(atkNow float64) {
	c.AtkNow = atkNow
}

func (c *BaseCard) GetTempId() int {
	return c.TempId
}
func (c *BaseCard) SetTempId(id int) {
	c.TempId = id

}
func (c *BaseCard) GetOwnerId() int {
	return c.OwnerId
}
func (c *BaseCard) SetOwnerId(id int) {
	c.OwnerId = id
}

func (c *BaseCard) GetID() int {
	return -1
}

func (c *BaseCard) SetInfo(info map[string]any) {
	c.Info = info
}
func (c *BaseCard) GetInfo() map[string]any {
	return c.Info
}

func (c *BaseCard) ReInitialize() {
	if val, ok := c.Info["hp"]; ok && val != nil {
		c.SetHpNow(c.Info["hp"].(float64))
	}
	if val, ok := c.Info["damage"]; ok && val != nil {
		c.SetAtkNow(c.Info["damage"].(float64))
	}
	if c.ID == 15 {
		c.SetHpNow(3)
	}
}

func (c *BaseCard) GetBuffList() *[]*protocol.Buff {
	return &c.BuffList
}

func (c *BaseCard) AppendBuff(b *protocol.Buff) {
	c.BuffList = append(c.BuffList, b)
}

// 这个函数是进过buff计算的
func (c *BaseCard) GetDec() *CardMeta.Decorator {
	NewDec := c.CalcDecByBuff(*c.Dec)
	return &NewDec
}

func (c *BaseCard) CalcDecByBuff(Dec CardMeta.Decorator) CardMeta.Decorator {
	for _, buff := range *c.GetBuffList() {
		switch buff.BuffId {
		case protocol.Vulnerability: //易伤
			Dec.HurtPordAdd(-buff.Value)
		case protocol.DamageImmunity: //免伤
			Dec.HurtPordAdd(buff.Value)
		case protocol.Block:
			Dec.HurtSumAdd(buff.Value)
		case protocol.HealingBoost: //治疗增强
			Dec.HealPordAdd(buff.Value)
		case protocol.HealingDecay: //治疗减弱
			Dec.HealPordAdd(-buff.Value)
		case protocol.BonusDamage:
			Dec.AttackSumAdd(buff.Value)
		case protocol.Powerful:
			Dec.AttackPordAdd(buff.Value)
		case protocol.Weakness:
			Dec.AttackPordAdd(-buff.Value)
		case protocol.XuFeng:
			Dec.EvadeAdd(buff.Value)
		}

	}
	return Dec
}

func (c *BaseCard) AddBuff(buff *protocol.Buff, pc protocol.ProtocolCardWithCtx) {
	c.AppendBuff(buff)
	protocol.BuffOnApplyFuncMap[buff.BuffId](pc, buff.Value, c) //执行挂载函数
}

func (c *BaseCard) BuffRoundEnd(pc protocol.ProtocolCardWithCtx) {
	for _, b := range c.BuffList { //循环结算当回合所有的buff
		protocol.BuffRoundEndFuncMap[b.BuffId](pc, b.Value, c) //每回合执行回合结束函数
		b.Stacks -= 1                                          //层数减一
		if b.Stacks == 0 {
			protocol.BuffOnRemoveFuncMap[b.BuffId](pc, b.Value, c)
			var targetIdx int
			for i, v := range c.BuffList {
				if v.TempId == b.TempId {
					targetIdx = i
					break
				}
			}
			// 循环结束后，如果找到了，再执行安全删除
			if targetIdx != -1 {
				c.BuffList = append(c.BuffList[:targetIdx], c.BuffList[targetIdx+1:]...)
			}
		} //buff结束了,执行buff结束函数
	}
}
func (c *BaseCard) GetForm() BattleData.Form {
	return c.Form
}

func (c *BaseCard) SetForm(form BattleData.Form) {
	c.Form = form
}

func (c *BaseCard) IntSpecialCardStateChan(goctx context.Context) {

	go func() {
		for {
			select {
			case v := <-c.SpecialCardStateChan:
				c.self.BroadCallBack(v)
			case <-goctx.Done():
				return
			}
		}
	}()
}

func (c *BaseCard) BroadCallBack(v *CardMeta.BroadInfo) {

}

// 所有的卡的一些重要的初始化在这
func (c *BaseCard) ShareInit(goctx context.Context, ctx protocol.ProtocolCardWithCtx) {
	c.self = c
	c.Dec = CardMeta.NewDecorator()
	c.BuffList = make([]*protocol.Buff, 0, 8)
	c.SetForm(BattleData.NormalForm)
	c.SetBtCtx(ctx)
	c.changeJiangShi = false
	c.CR = CardAbstract.NewCardRecord()
	var ch chan *CardMeta.BroadInfo
	ch = make(chan *CardMeta.BroadInfo, 4)
	c.SpecialCardStateChan = ch
	go c.IntSpecialCardStateChan(goctx)
}

func (c *BaseCard) PutBroadInfo(v *CardMeta.BroadInfo) {
	c.SpecialCardStateChan <- v
}

//---------二次分装---------

func (c *BaseCard) EffectAttack(targetTempId int, AtkHp float64, category BattleData.ValueChange) {
	c.BtCtx.ProtoColPush(protocol.NewAttack(c.OwnerId, c.TempId, targetTempId, AtkHp, c.GetDec(), category, false, &protocol.InterruptConfig{}))
}
func (c *BaseCard) EffectHurt(AttackId int, AtkHp float64, Category BattleData.ValueChange) {
	c.BtCtx.ProtoColPush(protocol.NewHurt(c.OwnerId, AttackId, c.TempId, AtkHp, c.GetDec(), Category))
}

func (c *BaseCard) EffectHeal(targetTempId int, HealHp float64) {
	c.BtCtx.ProtoColPush(protocol.NewHeal(&targetTempId, HealHp, c.GetDec(), false, &protocol.InterruptConfig{}))
}

func (c *BaseCard) EffectUpdateEnergy(offset int) {
	c.BtCtx.ProtoColPush(protocol.NewUpdateEnergy(c.OwnerId, offset))
}

// 通知标明行为发起者和受到者
func (c *BaseCard) Notify(Beh BattleData.AnimationBehavior, UserID int, CallerId int, AcceptorId int) { //都是tempid
	c.BtCtx.Notify(BattleData.NewAnimationDto(CallerId, AcceptorId, Beh), UserID)
}

func (c *BaseCard) DisCard(TempIdList []int) {
	c.BtCtx.ProtoColPush(protocol.NewDisCard(c.OwnerId, TempIdList, false, &protocol.InterruptConfig{}))
}

func (c *BaseCard) SetCardBt(TempId int) {
	c.BtCtx.ProtoColPush(protocol.NewSetCardBt(c.OwnerId, TempId, false, &protocol.InterruptConfig{}))
}

func (c *BaseCard) GiveBuff(TempId *int, b protocol.Buff) {
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(TempId, b, false, &protocol.InterruptConfig{}))
}

func (c *BaseCard) ReMoveBuffByTempId(BuffTempId int) {
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

func (c *BaseCard) NewCustom(ExecFunc func(pc protocol.ProtocolCardWithCtx)) {
	c.BtCtx.ProtoColPush(protocol.NewCustom(ExecFunc))
}

func (c *BaseCard) ChangeMaxHp(TargetTempId int, MaxHp float64) {
	c.BtCtx.ProtoColPush(protocol.NewChangeMaxHp(TargetTempId, MaxHp))
}

func (c *BaseCard) NextRound() {}

func (c *BaseCard) RoundEnd() {
	c.CR.RoundEnd()
}

func (c *BaseCard) GetCR() *CardAbstract.CardRecord {
	return c.CR
}
