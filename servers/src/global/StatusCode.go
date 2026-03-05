package global

type StatusCode int

const (
	StatusSuccess StatusCode = iota
	StatusDataNotFound
	StatusInternalServersError
	StatusInvalidReqParamsName
	StatusInvalidReqParamsClass
	StatusInvalidReqParams
	StatusInvalidToken
	StatusTokenExpired
	StatusIncorrectTokenFormat
)

var StatusMsg = map[StatusCode]string{
	StatusSuccess:               "成功",
	StatusDataNotFound:          "数据没找到",
	StatusInternalServersError:  "内部错误",
	StatusInvalidReqParamsName:  "非法请求参数名",
	StatusInvalidReqParamsClass: "非法请求参数类型",
	StatusInvalidReqParams:      "非法请求参数",
	StatusInvalidToken:          "非法token",
	StatusTokenExpired:          "token失效",
	StatusIncorrectTokenFormat:  "token格式错误",
}
