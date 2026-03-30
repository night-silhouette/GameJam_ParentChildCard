extends Node

# 游戏里的状态和凭证都存在这里
var current_token: String = ""
var my_player_id: int = 0
var is_admin: bool = false

func _ready():
	# 大管家搬个椅子坐在传达室，专听底层原始信号
	SignalBus.raw_api_responded.connect(_handle_raw_api_data)
	TokenManager.load_token();
	current_token =  TokenManager.get_token();	

# 这是大管家的核心大脑，完美对应你的后端文档
func _handle_raw_api_data(api_name: String, method: int, code: int, data: Variant, msg: String):
# 登录和token
	print("进入解包层")
	if code != 0 :
		print("code = ",code,NetError.get_message(code));
		return
		
	match api_name:

	# -------------------------------------
	# 1. 登录与 Token 相关 (/v1/token)
	# -------------------------------------
		"/v1/token/":
			if method == HTTPClient.METHOD_POST: 
			# POST 是登录
# 如果输出 String，说明你还没解析它，它只是个长得像字典的字符串
				if code == 0:
					TokenManager.save_token(str(data["token"])) # 存下Token
					current_token =  TokenManager.get_token();	
					SignalBus.login_success.emit()
				else:
					SignalBus.login_failed.emit(msg)

			elif method == HTTPClient.METHOD_GET: 
			# GET 是验证 Token 是否还活着
				if code == 0:
					
					SignalBus.token_validated_success.emit()
					print("token成功")
				else:
					current_token = ""
					SignalBus.network_disconnected.emit()
					print("toke过期 " + NetError.error_message[code]);

	# -------------------------------------
	# 2. 用户信息相关 (/v1/user)
	# -------------------------------------
		"/v1/user/":
			if method == HTTPClient.METHOD_GET: 
			# GET 是查资料
				if code == 0:
					my_player_id = int(data["id"])
					is_admin = bool(data["is_admin"])
					# 通知 UI 刷新名字
					SignalBus.user_info_fetched.emit(my_player_id, str(data["name"]), is_admin)

			elif method == HTTPClient.METHOD_POST:
				# POST 是注册
				if code == 0:
					SignalBus.user_registered_success.emit()
					print("注册成功")

				elif method == HTTPClient.METHOD_PUT:
				# PUT 是修改资料
					if code == 0:
						SignalBus.user_updated_success.emit()
						
