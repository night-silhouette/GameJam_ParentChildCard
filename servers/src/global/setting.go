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

	WaitAnimationRecall = 30 //等待动画结束

	SelectWeatherTime   = 10   //s
	BattleWaitTime      = 10   //s
	ActiveChildCardTime = 10   //s
	Interrupt           = 10   //s 中断
	SelectSkillCardTime = 10   //s
	JudgeWaitTime       = 1000 //s
	CombatWaitTime      = 100  //s
)
