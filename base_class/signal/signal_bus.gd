extends Node

# ==========================================
# 1. 游戏内/系统通用信号
# 用于处理本地场景切换、UI 刷新或底层网络状态
# ==========================================

## 切换场景。用法：SignalBus.change_scence.emit("res://levels/main.tscn")
signal change_scence(message : String)

## 切换 UI 面板。用法：SignalBus.change_ui.emit("login_panel")
signal change_ui(message : String)

## 网络连接断开时触发（如物理断网、服务器宕机）
signal network_disconnected()


# ==========================================
# 2. 网络请求信号 (Outgoing - 发送给服务器)
# 命名规范：以 request_ 开头。通常在 UI 脚本中 emit
# ==========================================

## [配对 A] 登录请求
signal request_login(username: String, password: String)

## [配对 B] 校验 Token 是否过期
signal request_validate_token()

## [配对 C] 获取当前登录用户的详细信息
signal request_get_user_self()

## [配对 C] 根据 ID 查询用户信息（管理权限）
signal request_get_user_by_id(user_id: int)

## [配对 C] 根据用户名查询用户信息（管理权限）
signal request_get_user_by_name(name: String)

## [配对 D] 注册新账号
signal request_register_user(name: String, password: String)

## [配对 E] 更新用户信息
signal request_update_user(id: int, name: String, password: String)

## [配对 F] 删除当前登录的账号
signal request_delete_user_self()

## [配对 G] 根据 ID 删除指定用户（管理权限）
signal request_delete_user_by_id(id: int)


# ==========================================
# 3. 网络接收信号 (Incoming - 服务器返回结果)
# 命名规范：以 success/failed/fetched 结尾
# ==========================================

## [通用] 原始 API 响应信号。包含所有请求的底层数据，用于日志审计或全局错误处理
signal raw_api_responded(api_name: String, method: int, code: int, data: Variant, msg: String)

## [配对 A] 登录成功：由后端处理逻辑验证后触发
signal login_success()

## [配对 A] 登录失败：返回失败原因 msg
signal login_failed(msg: String)

## [配对 B] Token 验证通过：说明用户不需要重新登录
signal token_validated_success()

## [配对 C] 用户信息获取成功：返回具体的 ID、用户名及权限等级
signal user_info_fetched(id: int, user_name: String, is_admin: bool)

## [配对 D] 注册流程完成
signal user_registered_success()

## [配对 E] 用户资料修改成功
signal user_updated_success()
