package BattleData

// ChildCardCatchDto 好丑陋,但就这样吧,从野生子牌堆catch下来的就Origin=-1
type ChildCardCatchDto struct {
	Origin  int      `json:"origin" mapstructure:"origin"`
	Object  int      `json:"object" mapstructure:"object"`
	DataAll *DataAll `json:"data_all" mapstructure:"data_all"`
}
