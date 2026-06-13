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

	SelectWeatherTime   = 1  //s
	BattleWaitTime      = 1  //s
	ActiveChildCardTime = 1  //s
	SelectCharacterTime = 1  //s
	SelectSkillCardTime = 1  //s
	JudgeWaitTime       = 15 //s
	CombatWaitTime      = 25 //s
)
