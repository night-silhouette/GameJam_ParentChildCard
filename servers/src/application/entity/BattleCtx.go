package entity

type BattleCtx struct {
	params map[string]any
}

func (c *BattleCtx) Set(key string, val any) {
	if c.params == nil {
		c.params = make(map[string]any)
	}
	c.params[key] = val
}

func (c *BattleCtx) Get(key string) (any, bool) {
	val, ok := c.params[key]
	return val, ok
}
