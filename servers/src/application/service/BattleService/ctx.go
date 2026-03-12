package BattleService

type Ctx struct {
	IDA    int
	IDB    int
	params map[string]any
}

func (c *Ctx) Set(key string, val any) {
	if c.params == nil {
		c.params = make(map[string]any)
	}
	c.params[key] = val
}

func (c *Ctx) Get(key string) (any, bool) {
	val, ok := c.params[key]
	return val, ok
}
