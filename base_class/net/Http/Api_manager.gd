extends Node

func _ready():
	# Token
	SignalBus.request_login.connect(_on_request_login)
	SignalBus.request_validate_token.connect(_on_request_validate_token)

	# 用户
	SignalBus.request_get_user_self.connect(_on_get_user_self)
	SignalBus.request_get_user_by_id.connect(_on_get_user_by_id)
	SignalBus.request_get_user_by_name.connect(_on_get_user_by_name)

	SignalBus.request_register_user.connect(_on_register_user)
	SignalBus.request_update_user.connect(_on_update_user)
	SignalBus.request_delete_user_self.connect(_on_delete_user_self)
	SignalBus.request_delete_user_by_id.connect(_on_delete_user_by_id)
	
	SignalBus.request_send_mail.connect(_request_send_mail)
	
	SignalBus.request_get_mail_numberN.connect(_request_get_mail_numberN);
	SignalBus.request_get_mail.connect(_request_get_mail);
	
	SignalBus.request_delete_mail.connect(_request_delete_mail);
	
	
# =========================
# Token
# =========================

func _on_request_login(username: String, password: String):
	var body = {
		"name": username,
		"password": password
	}
	NetworkClient.call_api("/v1/token/", HTTPClient.METHOD_POST, body, false)


func _on_request_validate_token():
	NetworkClient.call_api("/v1/token/", HTTPClient.METHOD_GET)


# =========================
# 用户 GET
# =========================

func _on_get_user_self():
	NetworkClient.call_api("/v1/user/", HTTPClient.METHOD_GET)


func _on_get_user_by_id(user_id: int):
	var url = "/v1/user?id=" + str(user_id)
	NetworkClient.call_api(url, HTTPClient.METHOD_GET)


func _on_get_user_by_name(name: String):
	var url = "/v1/user?name=" + name
	NetworkClient.call_api(url, HTTPClient.METHOD_GET)


# =========================
# 用户 POST（注册）
# =========================

func _on_register_user(name: String, password: String):
	var body = {
		"name": name,
		"password": password
	}
	print("zhuc")
	NetworkClient.call_api("/v1/user/", HTTPClient.METHOD_POST, body, false)


# =========================
# 用户 PUT（修改）
# =========================

func _on_update_user(id: int, name: String, password: String):
	var body = {
		"id": id,
		"name": name,
		"password": password
	}
	NetworkClient.call_api("/v1/user/", HTTPClient.METHOD_PUT, body)


# =========================
# 用户 DELETE
# =========================

func _on_delete_user_self():
	NetworkClient.call_api("/v1/user/", HTTPClient.METHOD_DELETE)


func _on_delete_user_by_id(id: int):
	var url = "/v1/user?id=" + str(id)
	NetworkClient.call_api(url, HTTPClient.METHOD_DELETE)

	
func _request_send_mail(id:int,txt:String):
	var url ="/v1/mail/";
	var body = {
		"accept_id" :id,
		"body" : txt
	}
	NetworkClient.call_api(url,HTTPClient.METHOD_POST,body);
	
func _request_get_mail_numberN():
	var url="/v1/mail/status/";
	NetworkClient.call_api(url,HTTPClient.METHOD_GET);
	
	
func _request_get_mail(page : int):
	var url = "/v1/mail/";
	var body ={
		"page" : page
	}
	print(body)
	NetworkClient.call_api(url,HTTPClient.METHOD_GET,body);
	
func _request_delete_mail(id:Array[int]):
	print("xin");
	var url = "/v1/mail/";
	var body ={
		"mail_id" : id
	}
	NetworkClient.call_api(url,HTTPClient.METHOD_DELETE,body);
