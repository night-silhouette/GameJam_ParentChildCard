extends Node



func call_api(api_name: String, method: int, body_data: Dictionary = {}, use_token: bool = true):
	var http = HTTPRequest.new()
	add_child(http)

	http.request_completed.connect(on_request_completed.bind(http, api_name, method))

	var headers = ["Content-Type: application/json"]

	if use_token:
		if Packethandler.current_token == "":
			pass
		var auth_string = "Authorization: " + Packethandler.current_token
		headers.append(auth_string)

	var body_string = ""
	if not body_data.is_empty():
		body_string = JSON.stringify(body_data)

	var final_url = Global.BASE_URL + api_name
	var err = http.request(final_url, headers, method, body_string)

	if err != OK:
		http.queue_free()
		SignalBus.raw_api_responded.emit(api_name, method, -1, null, "请求发送失败")


func on_request_completed(result, response_code, headers, body, http_node, api_name, method):
	http_node.queue_free()
	
	##print("状态",response_code);
	if result != HTTPRequest.RESULT_SUCCESS or response_code != 200:
		
		SignalBus.raw_api_responded.emit(api_name, method, -1, null, "网络异常")
		
		return
	
	##print("网络成功返回")
	var res = JSON.parse_string(body.get_string_from_utf8())
	
	var code = res["code"]
	var data = res["data"] if res.has("data") else null
	var msg = res["msg"]
	#快递员干完活了，直接扔给传达室，下班！
	SignalBus.raw_api_responded.emit(api_name, method, code, data, msg)
	

	
