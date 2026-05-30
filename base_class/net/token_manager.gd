extends Node

var token: String = ""

# 调试开关（true=正常存储，false=只用内存）

const SAVE_PATH = "user://token.save"


# =========================
# 保存 Token
# =========================
func save_token(new_token: String):
	token = new_token

	if not Global.token_save:
		return

	var data = {
		"token": new_token
	}

	var file = FileAccess.open(SAVE_PATH, FileAccess.WRITE)
	file.store_string(JSON.stringify(data))
	file.close()

	print("Token saved")


# =========================
# 读取 Token
# =========================
func load_token():
	if  not Global.token_save:
		return

	if not FileAccess.file_exists(SAVE_PATH):
		return

	var file = FileAccess.open(SAVE_PATH, FileAccess.READ)
	if not file:
		print("错误：无法打开文件，错误码：", FileAccess.get_open_error())
		return

	var text = file.get_as_text()
	file.close()

	if text.is_empty():
		print("错误：文件存在，但内容是空的！")
		return

	var data = JSON.parse_string(text)

	if data == null:
		pass;
		print("错误：JSON 解析失败，文本内容可能损坏：", text)	
	elif typeof(data) != TYPE_DICTIONARY:
		pass;
		print("错误：解析成功但不是字典，实际类型是：", typeof(data))
	else:
		token = data.get("token", "")
		print("成功加载 Token")


# =========================
# 获取 Token
# =========================
func get_token() -> String:
	return token


# =========================
# 清除 Token（退出登录用）
# =========================
func clear_token():
	token = ""

	if not Global.token_save:
		#print("调试模式：只清除内存 Token")
		return

	if FileAccess.file_exists(SAVE_PATH):
		DirAccess.remove_absolute(SAVE_PATH)

	print("Token cleared")
