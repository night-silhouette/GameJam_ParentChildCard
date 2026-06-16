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
	DeployCard //-------------------------------------
	Judge
	MatchSuccess
	AnimationPlayEnd
	Combat
	CardCalc
	Debug
	Interrupt //-----------------------------------
	GetDisCard
	SkillCardNotify //-----------------------------
	GetEnergy
	SelectWeather
	GetChildCardList
	ActiveChildCard
	GetWeather
	ChildBelongChange //-----------------------------------
	GetUserId
	PositionChange  //-----------------------------------
	AnimationNotify //-----------------------------------
	HpChange        //-----------------------------------
	Ping
	ReConnect
	WeatherNotify  //-------------------------
	BuffCalcNotify //-------------------------
	GetRoundNum
	BuffChange
)

var ActionName = map[ActionCode]string{
	Fault:                 "错误",
	CancelMatch:           "取消匹配",
	GetSelfCardInHard:     "获取自己的卡牌信息",
	GetOpponentCardInHard: "获取对手的卡牌信息",
	GetBtCardInfo:         "获取场上的战斗信息",
	OverBattle:            "结束战斗",
	StartBattle:           "开始战斗",
	DeployCard:            "部署牌",
	Judge:                 "战斗回合判断",
	MatchSuccess:          "匹配成功",
	AnimationPlayEnd:      "动画结束",
	Combat:                "执行战斗行动",
	CardCalc:              "卡牌效果结算",
	Debug:                 "测试",
	Interrupt:             "中断选牌",
	GetDisCard:            "查看弃牌堆",
	GetEnergy:             "查看能量值",
	SelectWeather:         "天气选择",
	GetChildCardList:      "查看子牌堆",
	ActiveChildCard:       "子牌激活选择",
	GetWeather:            "获取天气",
	ChildBelongChange:     "捕获子牌",
	PositionChange:        "结算卡牌移动通知",
	AnimationNotify:       "行为动画通知",
	HpChange:              "hp变化通知",
	Ping:                  "ping",
	ReConnect:             "断线重连",
	WeatherNotify:         "天气结算通知",
	SkillCardNotify:       "法术牌结算通知",
	BuffCalcNotify:        "buff结算通知",
	GetUserId:             "获取双方id",
	GetRoundNum:           "获取回合数",
	BuffChange:            "buff发生改变",
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
