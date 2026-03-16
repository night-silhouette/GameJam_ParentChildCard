extends Node
signal change_scence(message : String);
signal change_ui(message : String);


##请求结束信号
signal raw_api_responded(api_name: String, method: int, code: int, data: Variant, msg: String)



signal login_success()
signal login_failed(msg: String)
signal token_validated_success() # 验证旧Token有效
signal user_info_fetched(id: int, user_name: String, is_admin: bool)
signal user_registered_success()
signal user_updated_success()
signal network_disconnected()
