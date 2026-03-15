# PacketHandler.gd (建议作为单例 Autoload)

enum ResponseCode {
	SUCCESS = 0,
	DATA_NOT_FOUND = 1,
	INTERNAL_ERROR = 2,
	INVALID_PARAMS = 3,
	INVALID_TOKEN = 4,
	TOKEN_EXPIRED = 5,
	# ... 对应你的文档
}

func _on_response(_result, response_code, _headers, body, http_node):
	http_node.queue_free() # 释放搬运工
	
	var raw_data = body.get_string_from_utf8()
	var res = JSON.parse_string(raw_data)
	
	# 1. 检查 HTTP 层错误
	if response_code >= 400:
		SignalBus.network_error.emit("网络异常: " + str(response_code))
		return

	# 2. 解析你的标准 Res 结构
	var code = res.get("Code", -1)
	var data = res.get("Data", {})
	var msg = res.get("Msg", "")

	if code == ResponseCode.SUCCESS:
		_distribute_data(data)
	elif code == ResponseCode.INVALID_TOKEN or code == ResponseCode.TOKEN_EXPIRED:
		SignalBus.auth_failed.emit() # 弹出登录框
	else:
		SignalBus.logic_error.emit(msg) # 提示错误信息

func _distribute_data(data):
	# 这里是延展性的核心：根据 Data 里的内容决定发射什么信号
	if data is Dictionary:
		if data.has("token"):
			NetClient.token = data.token
			SignalBus.login_success.emit()
		if data.has("is_admin"):
			SignalBus.user_info_received.emit(data)
	elif data == "pong":
		print("服务器还活着")
