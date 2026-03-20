extends Node

var token: String = ""

const SAVE_PATH = "user://token.save"



# =========================
# 保存 Token
# =========================
func save_token(new_token: String):
	token = new_token
	print(token);
	var data = {
		"token": new_token
	}

	var file = FileAccess.open(SAVE_PATH, FileAccess.WRITE)
	var text = file.get_as_text()
	file.store_string(JSON.stringify(data))
	file.close()

	print("Token saved")


# =========================
# 读取 Token
# =========================

func load_token():
	if not FileAccess.file_exists(SAVE_PATH):
		print("错误：本地根本没有这个文件")
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
		print("错误：JSON 解析失败，文本内容可能损坏：", text)	
	elif typeof(data) != TYPE_DICTIONARY:
		print("错误：解析成功但不是字典，实际类型是：", typeof(data))
	else:
		token = data.get("token", "")
		print("成功加载 Token:", token)

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

	if FileAccess.file_exists(SAVE_PATH):
		DirAccess.remove_absolute(SAVE_PATH)

	print("Token cleared")
