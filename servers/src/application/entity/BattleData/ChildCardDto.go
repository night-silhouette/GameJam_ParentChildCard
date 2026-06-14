package BattleData

type ChildCardDto struct {
	CardDto
	ChildState ChildState `json:"child_state" mapstructure:"child_state"`
}

func NewChildCardDto(dto CardDto, ChildState ChildState) *ChildCardDto {
	res := &ChildCardDto{}
	res.CardDto = dto
	res.ChildState = ChildState
	return res
}

type ChildState int

const (
	Active ChildState = iota
	NotActive
	Died
	HasCatch
)
