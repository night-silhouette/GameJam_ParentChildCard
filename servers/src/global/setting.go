package global

const (
	TokenExpiredTime  = 30
	MatchLoopTime     = 500 //ms
	MatchTimeRadio    = 2
	MatchMaxWaitTime  = 3.0 //初定是18
	WsInterceptorTime = 350 //ms
)

var Isdebug string = "debug"

// 游戏参数
const (
	InitCardNum = 4

	Lev1Category1Num = 10
	Lev1Category2Num = 9
	Lev1Category3Num = 8
	Lev1Category4Num = 8

	SelectWeatherTime   = 10 //s
	BattleWaitTime      = 30 //s
	ActiveChildCardTime = 10 //s
	SelectCharacterTime = 10 //s
	SelectSkillCardTime = 15 //s
	JudgeWaitTime       = 7  //s
	CombatWaitTime      = 25 //s
)
