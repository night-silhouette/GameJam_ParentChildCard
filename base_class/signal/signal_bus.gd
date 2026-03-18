extends Node
signal change_scence(message : String);
signal change_ui(message : String);


##请求结束信号
signal raw_api_responded(api_name: String, method: int, code: int, data: Variant, msg: String)


##token和login的接受后处理信号
signal login_success()
signal login_failed(msg: String)
signal token_validated_success() # 验证旧Token有效
signal user_info_fetched(id: int, user_name: String, is_admin: bool)
signal user_registered_success()
signal user_updated_success()
signal network_disconnected()

# =========================
# Token / 登录
# =========================

## 登录，获取 token
signal request_login(username: String, password: String)

## 校验当前 token 是否有效
signal request_validate_token()

## 获取当前用户信息（普通用户用这个）
signal request_get_user_self()

## 按 ID 查询用户（管理员）
signal request_get_user_by_id(user_id: int)

## 按用户名查询用户（管理员）
signal request_get_user_by_name(name: String)

## 注册新用户
signal request_register_user(name: String, password: String)

## 修改用户信息（自己 or 管理员）
signal request_update_user(id: int, name: String, password: String)

## 删除自己账号
signal request_delete_user_self()

## 按 ID 删除用户（管理员）
signal request_delete_user_by_id(id: int)
