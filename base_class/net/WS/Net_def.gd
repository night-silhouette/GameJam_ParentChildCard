# NetDef.gd (建议作为独立脚本或 Autoload)
extends Node

enum Predicate {
	EMPTY = 0,
	NOTIFY = 1,
	QUERY = 2,
	RESULT = 3,
	FINISH = 4,
	SUCCEED = 5
}

enum Action {
	FAULT = 0,
	CANCEL_MATCH = 1,
	GET_SELF_CARDS = 2,
	GET_OPPONENT_CARDS = 3,
	GET_BT_INFO = 4,
	OVER_BATTLE = 5,
	START_BATTLE = 6,
	DEPLOY_CARD = 7,
	JUDGE = 8,
	MATCH_SUCCESS = 9,
	ANIMATION_END = 10,
	COMBAT = 11,
	CardCalc = 12,
	Debug = 13,
	Interrupt = 14,
	GetDisCard = 15,
	SkillCardCalc = 16,
	GetEnergy = 17,
	SelectWeather = 18,
	GetChildCardList = 19,
	ActiveChildCard = 20,
	GetWeather = 21,           
	CatchChild = 22,      
	CardMove =  23,       
	AnimationNotify = 24,
	ValueNotify = 25,           
}

const ACTION_NAME = {
	Action.CANCEL_MATCH: "取消匹配",
	Action.GET_SELF_CARDS: "获取自己的卡牌信息",
	Action.GET_OPPONENT_CARDS: "获取对手的卡牌信息",
	Action.GET_BT_INFO: "获取场上的战斗信息",
	Action.OVER_BATTLE: "结束战斗",
	Action.START_BATTLE: "开始战斗",
	Action.DEPLOY_CARD: "部署一张牌",
	Action.JUDGE: "战斗回合判断",
	Action.MATCH_SUCCESS: "匹配成功",
	Action.ANIMATION_END: "动画结束",
	Action.COMBAT: "执行战斗行动",
	Action.CardCalc: "卡牌效果结算",
	Action.Debug: "测试",
	Action.Interrupt: "中断选牌",
	Action.GetDisCard: "查看弃牌堆",
	Action.SkillCardCalc: "法术牌的计算",
	Action.GetEnergy: "查看能量值",
	Action.SelectWeather: "选择天气",
	Action.GetChildCardList: "查看子牌堆",
	Action.ActiveChildCard: "激活子卡牌",
	Action.GetWeather:            "获取天气",
	Action.CatchChild:            "捕获子牌",
	Action.CardMove:              "结算卡牌移动通知",
	Action.AnimationNotify:       "行为动画通知",
	Action.ValueNotify:           "数值变化通知",
	
}

# 辅助函数：快速获取名字打印日志
static func get_action_name(code: int) -> String:
	return ACTION_NAME.get(code, "未知动作(Code:%d)" % code)
	
func get_predicate_name(value: int) -> String:
	# 检查索引是否在有效范围内
	if value >= 0 and value < Predicate.size():
		return Predicate.keys()[value]
	return "Unknown"
