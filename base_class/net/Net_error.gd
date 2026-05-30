extends Node

# 从0开始的错误码
enum ResponseStatusCode {
	ResponseSuccess,
	ResponseDataNotFound,
	ResponseInternalServersError,
	ResponseInvalidReqParams,
	ResponseInvalidToken,
	ResponseTokenExpired,
	ResponseIncorrectTokenFormat,
	ResponseDuplicateDataEntry,
	ResponseRequiredParamsMissing,
	ResponseDependentRecordsExist,
	ResponseNotImplemented,
	ResponseIncorrectPassword,
	ResponseTokenMissing,
	ResponseForbidden,
	ResponseRepeatRequest,
	ResponseUnknownError,
	ResponseTokenHasUpdate,
	BattleInvalidTiming,
	BattleEffectStackOverflow,
	BattleCardCategoryError,
	BattleCardNotFound,
	BattleNotInYourRound,
	BattleHasCard,
	BattleCardNumErr
}

var error_message := {
	ResponseStatusCode.ResponseSuccess: "成功",

	ResponseStatusCode.ResponseDataNotFound: "数据不存在",

	ResponseStatusCode.ResponseInternalServersError: "服务器内部错误",

	ResponseStatusCode.ResponseInvalidReqParams: "请求参数无效",

	ResponseStatusCode.ResponseInvalidToken: "Token无效",

	ResponseStatusCode.ResponseTokenExpired: "Token已过期",

	ResponseStatusCode.ResponseIncorrectTokenFormat: "Token格式错误",

	ResponseStatusCode.ResponseDuplicateDataEntry: "数据重复",

	ResponseStatusCode.ResponseRequiredParamsMissing: "缺少必要参数",

	ResponseStatusCode.ResponseDependentRecordsExist: "存在依赖数据，无法操作",

	ResponseStatusCode.ResponseNotImplemented: "功能未实现",

	ResponseStatusCode.ResponseIncorrectPassword: "密码错误",

	ResponseStatusCode.ResponseTokenMissing: "缺少Token",

	ResponseStatusCode.ResponseForbidden: "没有权限",
	ResponseStatusCode.ResponseRepeatRequest: "重复请求",
	ResponseStatusCode.ResponseUnknownError: "未知错误",
	ResponseStatusCode.ResponseTokenHasUpdate: "Token已被更新",
	ResponseStatusCode.BattleInvalidTiming: "不在正确的战斗时机",
	ResponseStatusCode.BattleEffectStackOverflow:     "卡牌效果结算堆栈溢出",
	ResponseStatusCode.BattleCardCategoryError:       "卡牌种类有误",
	ResponseStatusCode.BattleCardNotFound:            "此牌没有被找到",
	ResponseStatusCode.BattleNotInYourRound:          "此时不在你的回合",
	ResponseStatusCode.BattleHasCard:                 "此位置有牌了",
	ResponseStatusCode.BattleCardNumErr:              "卡牌数量有误"

}

func get_message(code:int) -> String:
	return error_message.get(code, "未知错误")
	
func is_success(code:int) -> bool:
	return code == ResponseStatusCode.ResponseSuccess
