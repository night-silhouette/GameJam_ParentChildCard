package CardImpl

import (
	_ "embed"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
)

type BaseCard struct {
	ID            int            `json:"id"`
	Info          map[string]any `json:"-"`
	StateCodeChan chan protocol.Effect
	//动态变量
	BtCtx            protocol.ProtocolCardWithCtx
	HpNow            float64
	AtkNow           float64
	TempId           int
	OwnerId          int
	BuffList         []*protocol.Buff
	Dec              *CardMeta.Decorator
	ControlSignalMap map[string]CardMeta.ControlSignal
	Form             BattleData.Form
}

func (c *BaseCard) SetBtCtx(btCtx protocol.ProtocolCardWithCtx) {
	c.BtCtx = btCtx
}

func (c *BaseCard) GetHpNow() float64 {
	return c.HpNow
}
func (c *BaseCard) SetHpNow(hpNow float64) {
	c.HpNow = hpNow
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

func (c *BaseCard) GetStateCodeChan() chan protocol.Effect {
	return c.StateCodeChan
}
func (c *BaseCard) SetStateCodeChan(ch chan protocol.Effect) {
	c.StateCodeChan = ch
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
func (c *BaseCard) InitBuffList() {
	c.BuffList = make([]*protocol.Buff, 0, 8)
}

func (c *BaseCard) SetDec(Dec *CardMeta.Decorator) {
	c.Dec = Dec
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

func (c *BaseCard) InitControlSignalMap() {
	c.ControlSignalMap = make(map[string]CardMeta.ControlSignal)
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
