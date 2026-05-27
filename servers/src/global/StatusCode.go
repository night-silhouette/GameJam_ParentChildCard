package global

type ResponseStatusCode int

const (
	ResponseSuccess ResponseStatusCode = iota
	ResponseDataNotFound
	ResponseInternalServersError
	ResponseInvalidReqParams
	ResponseInvalidToken
	ResponseTokenExpired
	ResponseIncorrectTokenFormat
	ResponseDuplicateDataEntry
	ResponseRequiredParamsMissing
	ResponseDependentRecordsExist
	ResponseNotImplemented
	ResponseIncorrectPassword
	ResponseTokenMissing
	ResponseForbidden
	ResponseRepeatRequest
	ResponseUnknownError
	ResponseTokenHasUpdate
	ResponseBagsUnknownError

	BattleInvalidTiming
	BattleEffectStackOverflow
	BattleCardCategoryError
	BattleCardNotFound
	BattleNotInYourRound
	BattleHasCard
	BattleCardNumErr
)

var StatusMsg = map[ResponseStatusCode]string{
	ResponseSuccess:               "成功",
	ResponseDataNotFound:          "数据没找到",
	ResponseInternalServersError:  "服务器未知内部错误",
	ResponseInvalidReqParams:      "非法请求参数",
	ResponseInvalidToken:          "非法token",
	ResponseTokenExpired:          "token失效",
	ResponseIncorrectTokenFormat:  "token格式错误",
	ResponseTokenMissing:          "token缺失",
	ResponseDuplicateDataEntry:    "重复数据录入",
	ResponseRequiredParamsMissing: "参数值缺失",
	ResponseDependentRecordsExist: "存在关联数据",
	ResponseNotImplemented:        "接口未完善",
	ResponseIncorrectPassword:     "密码错误",
	ResponseForbidden:             "权限不足",
	ResponseRepeatRequest:         "重复请求",
	ResponseUnknownError:          "发生了一个未知错误，抱歉",
	ResponseTokenHasUpdate:        "token被更新,此token失效",
	BattleInvalidTiming:           "不在正确的战斗时机",
	BattleEffectStackOverflow:     "卡牌效果结算堆栈溢出",
	BattleCardCategoryError:       "卡牌种类有误",
	BattleCardNotFound:            "此牌没有被找到",
	BattleNotInYourRound:          "此时不在你的回合",
	BattleHasCard:                 "此位置有牌了",
	BattleCardNumErr:              "卡牌数量有误",
}
