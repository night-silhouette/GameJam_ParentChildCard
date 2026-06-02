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

	Lev1Category1Num = 6
	Lev1Category2Num = 8
	Lev1Category3Num = 5
	Lev1Category4Num = 5

	SelectWeatherTime   = 8    //s
	BattleWaitTime      = 1    //s
	ActiveChildCardTime = 8    //s
	SelectCharacterTime = 10   //s
	SelectSkillCardTime = 1    //s
	JudgeWaitTime       = 1    //s
	CombatWaitTime      = 1000 //s
)
