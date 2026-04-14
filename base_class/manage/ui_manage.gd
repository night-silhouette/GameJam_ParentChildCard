extends Node

var current_ui : Node
var root_node : Node

func _ready() -> void:
	# 监听UI切换信号
	SignalBus.change_ui.connect(fui_change)

func register_root(node: Node):
	root_node = node


func fui_change(state):


	var next_path : String

	match state:
		"tologin":
			next_path = "res://game_ui/login/login_ui.tscn"
		"tomenu":
			next_path = "res://game_ui/menu/menu_ui.tscn"
		"tobattle":
			next_path = "res://game_ui/battle/battle_ui.tscn"
	await goto_ui(next_path)


func goto_ui(path: String):

	if root_node == null:
		push_error("UI root_node not registered")
		return

	# 1 开始后台加载
	ResourceLoader.load_threaded_request(path)

	# 2 等待加载完成
	while ResourceLoader.load_threaded_get_status(path) == ResourceLoader.THREAD_LOAD_IN_PROGRESS:
		await get_tree().process_frame

	# 3 获取加载结果
	var ui_res = ResourceLoader.load_threaded_get(path)

	if ui_res == null:
		push_error("UI load failed: " + path)
		return

	# 4 实例化
	var new_ui = ui_res.instantiate()

	# 5 先加新的UI
	root_node.add_child(new_ui)

	# 6 再删除旧UI（防止一帧空UI）
	if current_ui:
		current_ui.queue_free()

	current_ui = new_ui
