package CardMeta

import "math/rand"

type Decorator struct {
	AttackSum  float64
	AttackPord float64 //加算的
	HurtSum    float64
	HurtPord   float64 //(还收到多少伤害)
	HealPord   float64 //
	Evade      float64 //闪避率//小数
}

func (c *Decorator) CalcHeal(OriginValue float64) int {
	return int(OriginValue * c.HealPord)
}

func (c *Decorator) CalcAttack(OriginValue float64) int {
	return int((OriginValue + c.AttackSum) * c.AttackPord)
}

// 输出:数值,是否被闪避
func (c *Decorator) CalcHurt(AttackValue float64) (int, bool) {
	flag := false
	res := int(AttackValue*c.HurtPord + c.HurtSum)
	if float64(rand.Intn(100)) < c.Evade*100 {
		res = 0
		flag = true
	}
	return res, flag
}

// 增加攻击比例
func (c *Decorator) AttackPordAdd(AttackPord float64) {
	c.AttackPord += AttackPord
}

// 增加免伤率,输入的是免伤值,可正可负
func (c *Decorator) HurtPordAdd(HurtPord float64) {
	c.HurtPord *= 1 - HurtPord
}

func (c *Decorator) HealPordAdd(HealPord float64) {
	c.HealPord += HealPord
}

func (c *Decorator) AttackSumAdd(AttackSum float64) {
	c.AttackSum += AttackSum
}

// 增加格挡值
func (c *Decorator) HurtSumAdd(HurtSum float64) {
	c.HurtSum += HurtSum
}

func (c *Decorator) EvadeAdd(EvadeValue float64) {
	c.Evade += EvadeValue
}

func NewDecorator() *Decorator {
	res := Decorator{
		AttackSum:  0,
		AttackPord: 1,
		HurtSum:    0,
		HurtPord:   1,
		HealPord:   1,
		Evade:      0,
	}
	return &res
}
