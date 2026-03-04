package global

type StatusCode int

const (
	StatusSuccess StatusCode = iota
	StatusDataNotFound
	StatusInternalServersError
	StatusInvalidReqParams
)

var StatusMsg = map[StatusCode]string{
	StatusSuccess:              "成功",
	StatusDataNotFound:         "数据没找到",
	StatusInternalServersError: "内部错误",
	StatusInvalidReqParams:     "非法请求参数",
}
