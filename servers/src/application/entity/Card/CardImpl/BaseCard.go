package CardImpl

import (
	_ "embed"
	"pcc_card/application/entity/Card/CardAbstract"
)

type BaseCard struct {
	ID            int            `json:"id"`
	Info          map[string]any `json:"-"`
	StateCodeChan chan CardAbstract.StateCode
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

func (c *BaseCard) GetStateCodeChan() chan CardAbstract.StateCode {
	return c.StateCodeChan
}
func (c *BaseCard) SetStateCodeChan(ch chan CardAbstract.StateCode) {
	c.StateCodeChan = ch
}
