extends Node
const BASE_URL = "http://120.26.145.68:10086" 


func call_api(api_name: String, method: int, body_data: Dictionary = {}, use_token: bool = true):
	# 1. 招募一个临时的快递员
	print("api调用成功")
	var http = HTTPRequest.new()
	add_child(http) # 让他归 Net 管
	
	http.request_completed.connect(on_request_completed.bind(http, api_name, method))
	
	# 3. 准备（Headers）
	var headers = ["Content-Type: application/json"]
	
	# 这里就是决定要不要使用token
	if use_token:
		
		# 如果当前没有 Token（可能还没登录），在控制台提个醒，防呆设计
		if Packethandler.current_token == "":
			print("警告：请求 " + api_name + " 需要 Token，但当前没找到 Token！")
		# 把 Token 拼接到规定的 Authorization 格式里，塞进信封
		var auth_string = "Authorization: " + Packethandler.current_token
		headers.append(auth_string)
	
	# 4. （Body）
	var body_string = ""
	# 如果字典里有东西，而且不是 GET 请求（GET不用带body），就把它转成 JSON 罐头
	if not body_data.is_empty() :
		body_string = JSON.stringify(body_data)
	#print(body_string);
	# 5. 快递员出发！
	var final_url = BASE_URL + api_name
	var err = http.request(final_url, headers, method, body_string)
	
	# 防御性编程：如果连门都没出（比如 URL 格式写错了），直接销毁快递员，省得卡在那儿
	if err != OK:
		push_error("请求发送失败，API: " + api_name)
		http.queue_free()


func on_request_completed(result, response_code, headers, body, http_node, api_name, method):
	http_node.queue_free()
	
	print("状态",response_code);
	if result != HTTPRequest.RESULT_SUCCESS or response_code != 200:
		
		SignalBus.raw_api_responded.emit(api_name, method, -1, null, "网络异常")
		
		return
	
	print("网络成功返回")
	var res = JSON.parse_string(body.get_string_from_utf8())
	
	var code = res["code"]
	var data = res["data"] if res.has("data") else null
	var msg = res["msg"]
	#快递员干完活了，直接扔给传达室，下班！
	SignalBus.raw_api_responded.emit(api_name, method, code, data, msg)
	

	
