extends Node

var myid;
func _ready():
	print("=== FULL API TEST START ===")

	# 监听所有关键结果
	_connect_signals()

	await get_tree().create_timer(1.0).timeout

	await test_register()
	await test_login()
	await test_validate_token()
	await test_get_self()
	await test_update_user()

	print("=== TEST FLOW END ===")


# =========================
# 信号监听（核心）
# =========================

func _connect_signals():
	SignalBus.login_success.connect(func():
		print("[OK] login success")
	)

	SignalBus.login_failed.connect(func(msg):
		print("[FAIL] login failed:", msg)
	)

	SignalBus.token_validated_success.connect(func():
		print("[OK] token valid")
	)

	SignalBus.network_disconnected.connect(func():
		print("[FAIL] token invalid / disconnected")
	)

	SignalBus.user_info_fetched.connect(func(id, name, is_admin):
		print("[OK] user info:", id, name, is_admin)
		myid = id;
	)

	SignalBus.user_registered_success.connect(func():
		print("[OK] register success")
	)

	SignalBus.user_updated_success.connect(func():
		print("[OK] update success")
	)


# =========================
# 测试流程（带等待）
# =========================

func test_register():
	print("\n[TEST] register")
	SignalBus.request_register_user.emit("test_user1", "123456")
	await wait()


func test_login():
	print("\n[TEST] login")
	SignalBus.request_login.emit("test_user1", "123456")
	await wait()


func test_validate_token():
	print("\n[TEST] validate token")
	SignalBus.request_validate_token.emit()
	await wait()


func test_get_self():
	print("\n[TEST] get self")
	SignalBus.request_get_user_self.emit()
	await wait()


func test_update_user():
	print("\n[TEST] update user")

	# ⚠️ 这里用你状态管理里的 id
	var id = myid;
	SignalBus.request_update_user.emit(id, "new_name", "654321")

	await wait()


# =========================
# 等待（防止请求挤一起）
# =========================

func wait():
	return get_tree().create_timer(1.5).timeout
