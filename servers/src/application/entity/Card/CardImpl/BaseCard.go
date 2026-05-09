package CardImpl

import (
	_ "embed"
	"pcc_card/application/entity/protocol"
)

type BaseCard struct {
	ID            int            `json:"id"`
	Info          map[string]any `json:"-"`
	StateCodeChan chan protocol.Effect
	//动态变量
	BtCtx   protocol.ProtocolCardWithCtx
	HpNow   float64
	AtkNow  float64
	TempId  int
	OwnerId int
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
