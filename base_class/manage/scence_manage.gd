extends Node

var current_scene : Node
var root_node : Node

func _ready() -> void:
	SignalBus.change_scence.connect(fchange_scence)

func register_root(node: Node):
	root_node = node


func fchange_scence(state):

	var next_path : String

	match state:
		"start":
			next_path = "res://game_scene/main/start_scence.tscn"

		"tologin":
			next_path = "res://game_scene/login/login_scence.tscn"

	await goto_scene(next_path)


func goto_scene(path: String) -> void:

	if root_node == null:
		push_error("root_node not registered")
		return

	# 1 异步加载请求
	ResourceLoader.load_threaded_request(path)

	# 2 等待加载完成
	while ResourceLoader.load_threaded_get_status(path) == ResourceLoader.THREAD_LOAD_IN_PROGRESS:
		await get_tree().process_frame

	# 3 获取资源
	var scene_res = ResourceLoader.load_threaded_get(path)

	if scene_res == null:
		push_error("scene load failed")
		return

	# 4 删除旧场景
	if current_scene:
		current_scene.queue_free()

	# 5 实例化
	current_scene = scene_res.instantiate()

	# 6 添加
	root_node.add_child(current_scene)
