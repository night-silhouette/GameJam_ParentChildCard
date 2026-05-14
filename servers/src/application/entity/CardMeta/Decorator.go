package CardMeta

type Decorator struct {
	Additive   float64
	Multiplier float64
}

func (c *Decorator) Calc(Value float64) float64 {
	return (Value + c.Additive) * c.Multiplier
}

func (c *Decorator) Add(NewD Decorator) {
	c.Additive = c.Additive + NewD.Additive
	c.Multiplier = c.Multiplier * NewD.Multiplier
}

func NewDecorator() *Decorator {
	res := Decorator{}
	res.Multiplier = 1.0
	res.Additive = 1.0
	return &res
}
