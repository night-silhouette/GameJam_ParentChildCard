package BattleData

type GoodsDto struct {
	Price    int     `json:"price" mapstructure:"price"`
	CardId   int     `json:"card_id" mapstructure:"card_id"`
	GoodsId  int     `json:"goods_id" mapstructure:"goods_id"`
	Discount float64 `json:"discount" mapstructure:"discount"`
}

func NewGoodsDto(Price int, CardId int, GoodsId int, Discount float64) *GoodsDto {
	res := GoodsDto{}
	res.Discount = Discount
	res.Price = Price
	res.CardId = CardId
	res.GoodsId = GoodsId
	return &res
}
