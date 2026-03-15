extends Node
const BASE_URL = "http://你的服务器IP或域名" 

# 【重点关注：Token 的临时住所】

var current_token: String = "" 

# ==========================================
# 第一层：底层基石 (万能快递员)
# ==========================================
func _call_api(api_name: String, method: int, body_data: Dictionary = {}, use_token: bool = true):
	# 1. 招募一个临时的快递员
	var http = HTTPRequest.new()
	add_child(http) # 让他归 Net 管
	

	http.request_completed.connect(_on_request_completed.bind(http, api_name))
	
	# 3. 准备（Headers）
	var headers = ["Content-Type: application/json"]
	
	# 这里就是决定要不要使用token
	if use_token:
		# 如果当前没有 Token（可能还没登录），在控制台提个醒，防呆设计
		if current_token == "":
			push_warning("警告：请求 " + api_name + " 需要 Token，但当前没找到 Token！")
			
		# 把 Token 拼接到规定的 Authorization 格式里，塞进信封
		var auth_string = "Authorization: " + current_token
		headers.append(auth_string)

	
	# 4. 准备信件正文（Body）
	var body_string = ""
	# 如果字典里有东西，而且不是 GET 请求（GET不用带body），就把它转成 JSON 罐头
	if not body_data.is_empty() and method != HTTPClient.METHOD_GET:
		body_string = JSON.stringify(body_data)
		
	# 5. 快递员出发！
	var final_url = BASE_URL + api_name
	var err = http.request(final_url, headers, method, body_string)
	
	# 防御性编程：如果连门都没出（比如 URL 格式写错了），直接销毁快递员，省得卡在那儿
	if err != OK:
		push_error("请求发送失败，API: " + api_name)
		http.queue_free()

# ==========================================
# 统一的回执处理中心
# ==========================================
func _on_request_completed(result, response_code, headers, body, http_node, api_name):
	# 【最重要的一步】：快递员使命达成了，请他回老家（清理内存）
	http_node.queue_free()

	# 1. 邮差半路翻车了？(断网了)
	if result != HTTPRequest.RESULT_SUCCESS:
		push_error("网络连接失败！API: " + api_name)
		SignalBus.api_responded.emit(api_name, -1, null, "网络连接失败")
		return

	# 2. 对方俱乐部不给进门？(比如 404，500)
	if response_code != 200:
		push_error("服务器异常！HTTP 状态码: " + str(response_code))
		SignalBus.api_responded.emit(api_name, -1, null, "服务器异常 (" + str(response_code) + ")")
		return

	# 3. 顺利拆包裹，拿到真正的服务器回复
	var json_text = body.get_string_from_utf8()
	var res = JSON.parse_string(json_text)

	# 4. 判断包裹是不是空包或者假包
	if typeof(res) != TYPE_DICTIONARY or not res.has("Code"):
		push_error("服务器返回的格式不对！解析失败。")
		SignalBus.api_responded.emit(api_name, -1, null, "数据格式错误")
		return

	# 5. 一切安好，大喇叭广播（发出信号，带上真正的业务数据）
	var data = null
	if res.has("Data"):
		data = res["Data"]
		
	SignalBus.api_responded.emit(api_name, res["Code"], data, res["Msg"])
