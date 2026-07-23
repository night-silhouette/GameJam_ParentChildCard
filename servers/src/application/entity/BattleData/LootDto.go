package BattleData

type LootDto struct {
	Data   []int `json:"data" mapstructure:"data"`
	LootID int   `json:"loot_id" mapstructure:"loot_id"`
}
