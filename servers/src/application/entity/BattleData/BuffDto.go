package BattleData

type BuffDto struct {
	BuffId     int `json:"buff_id" mapstructure:"buff_id"`
	BuffStacks int `json:"buff_stacks" mapstructure:"buff_stacks"`
}
