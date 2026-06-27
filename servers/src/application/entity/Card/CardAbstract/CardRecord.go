package CardAbstract

type CardRecord struct {
	HurtedThisTurn float64
}

func NewCardRecord() *CardRecord {
	return &CardRecord{}
}

func (c *CardRecord) RoundEnd() {
	c.HurtedThisTurn = 0
}
