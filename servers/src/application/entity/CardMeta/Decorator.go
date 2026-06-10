package CardMeta

type Decorator struct {
	AttackSum  float64
	AttackPord float64 //加算的
	HurtSum    float64
	HurtPord   float64 //(还收到多少伤害)
	HealPord   float64 //
}

func (c *Decorator) CalcHeal(OriginValue float64) int {
	return int(OriginValue * c.HealPord)
}

func (c *Decorator) CalcAttack(OriginValue float64) int {
	return int((OriginValue + c.AttackSum) * c.AttackPord)
}

func (c *Decorator) CalcHurt(AttackValue float64) int {
	return int(AttackValue*c.HurtPord + c.HurtSum)
}

func (c *Decorator) AttackPordAdd(AttackPord float64) {
	c.AttackPord += AttackPord
}
func (c *Decorator) HurtPordAdd(HurtPord float64) {
	c.HurtPord *= 1 - HurtPord
}

func NewDecorator() *Decorator {
	res := Decorator{
		AttackSum:  0,
		AttackPord: 1,
		HurtSum:    0,
		HurtPord:   1,
		HealPord:   1,
	}
	return &res
}
