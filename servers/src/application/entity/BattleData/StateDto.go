package BattleData

type StateCode int

const (
	InitState StateCode = iota
	ActiveChildCard
	SelectWeather
	SelectSkillCard
	Judge
	Combat
	CardCalc
)

// 根据字符串名字映射statecode
var StateMap = map[string]StateCode{
	"ShuffleDeal":     InitState,
	"SelectSkillCard": SelectSkillCard,
	"Judge":           Judge,
	"Combat":          Combat,
	"CardCalc":        CardCalc,
	"SelectWeather":   SelectWeather,
	"ActiveChildCard": ActiveChildCard,
}
