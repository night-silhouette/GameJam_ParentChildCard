extends Node

var myid;
func _ready():
	print("=== FULL API TEST START ===")


	test_request_delete_mail();
	
	
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
	SignalBus.request_login.emit("fht", "20051118fht")
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

func test_send_mail():
	print("send")
	
	SignalBus.request_send_mail.emit(14,"12345467")
	
func test_get_mail_number():
	
	SignalBus.request_get_mail_numberN.emit();
	
func test_get_mail():
	SignalBus.request_get_mail.emit(2);
	
func test_request_delete_mail():
	print("emit 前")
	var data: Array[int] = [47, 46];
	SignalBus.request_delete_mail.emit(data)
	print("emit 后")
# =========================
# 等待（防止请求挤一起）
# =========================

func wait():
	return get_tree().create_timer(1.5).timeout
