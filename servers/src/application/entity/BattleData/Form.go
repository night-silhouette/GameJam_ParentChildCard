package BattleData

type Form int

const (
	NormalForm Form = iota
	JiangShi
	Die
	E
)

type FormValue struct {
	Hp       float64
	Damage   float64
	Priority int
}

func NewFormValue(Hp float64, Damage float64, Priority int) *FormValue {
	res := &FormValue{}
	res.Hp = Hp
	res.Damage = Damage
	res.Priority = Priority
	return res
}

var FormValuesMap = map[Form]*FormValue{
	JiangShi: NewFormValue(1, 1, 0),
	Die:      NewFormValue(4, 3, 2),
	E:        NewFormValue(2, 2, 1),
}
