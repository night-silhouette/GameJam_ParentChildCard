package BattleDto

import "pcc_card/global"

type Action struct {
	ActionCode ActionCode `json:"action_code"`
	ActionName string     `json:"action_name"`
	ActionData any        `json:"action_data"`
	Predicates Predicates `json:"predicates"`
}

type ActionCode int

const (
	Fault ActionCode = iota
	CancelMatch
	GetSelfCardInHard
	GetOpponentCardInHard
	GetBtCardInfo
	OverBattle
	StartBattle
	DeployCard
	Judge
	MatchSuccess
	AnimationPlayEnd
	Combat
	CardCalc
	Debug
	Interrupt
	GetDisCard
	SkillCardCalc
	GetEnergy
	SelectWeather
	GetChildCardList
	ActiveChildCard
	GetWeather
	CatchChild
)

var ActionName = map[ActionCode]string{
	Fault:                 "错误",
	CancelMatch:           "取消匹配",
	GetSelfCardInHard:     "获取自己的卡牌信息",
	GetOpponentCardInHard: "获取对手的卡牌信息",
	GetBtCardInfo:         "获取场上的战斗信息",
	OverBattle:            "结束战斗",
	StartBattle:           "开始战斗",
	DeployCard:            "部署一张牌",
	Judge:                 "战斗回合判断",
	MatchSuccess:          "匹配成功",
	AnimationPlayEnd:      "动画结束",
	Combat:                "执行战斗行动",
	CardCalc:              "卡牌效果结算",
	Debug:                 "测试",
	Interrupt:             "中断选牌",
	GetDisCard:            "查看弃牌堆",
	SkillCardCalc:         "法术牌计算",
	GetEnergy:             "查看能量值",
	SelectWeather:         "天气选择",
	GetChildCardList:      "查看子牌堆",
	ActiveChildCard:       "子牌激活选择",
	GetWeather:            "获取天气",
	CatchChild:            "捕获子牌",
}

type Predicates int

const (
	Empty Predicates = iota
	Notify
	Query
	Result
	Finish
	Succeed
)

func NewAction(actionCode ActionCode, Predicates Predicates, ActionData any) Action {
	res := Action{}
	res.ActionCode = actionCode
	res.ActionName = ActionName[actionCode]
	res.ActionData = ActionData
	res.Predicates = Predicates
	return res
}

func NewErrAction(code global.ResponseStatusCode) Action {
	res := Action{}
	res.ActionCode = Fault
	res.Predicates = Empty
	res.ActionData = code
	return res
}
