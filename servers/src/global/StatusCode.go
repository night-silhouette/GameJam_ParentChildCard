package global

type ResponseStatusCode int

const (
	ResponseSuccess ResponseStatusCode = iota
	ResponseDataNotFound
	ResponseInternalServersError
	ResponseInvalidReqParamsName
	ResponseInvalidReqParamsClass
	ResponseInvalidReqParams
	ResponseInvalidToken
	ResponseTokenExpired
	ResponseIncorrectTokenFormat
	ResponseDuplicateDataEntry
	ResponseRequiredParamsMissing
	ResponseDependentRecordsExist
)

var StatusMsg = map[ResponseStatusCode]string{
	ResponseSuccess:               "成功",
	ResponseDataNotFound:          "数据没找到",
	ResponseInternalServersError:  "服务器未知内部错误",
	ResponseInvalidReqParamsName:  "非法请求参数名",
	ResponseInvalidReqParamsClass: "非法请求参数类型",
	ResponseInvalidReqParams:      "非法请求参数",
	ResponseInvalidToken:          "非法token",
	ResponseTokenExpired:          "token失效",
	ResponseIncorrectTokenFormat:  "token格式错误",
	ResponseDuplicateDataEntry:    "重复数据录入",
	ResponseRequiredParamsMissing: "参数值缺失",
	ResponseDependentRecordsExist: "存在关联数据",
}
