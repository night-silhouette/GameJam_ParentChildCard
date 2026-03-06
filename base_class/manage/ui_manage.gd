extends Node

var current_ui : Node;
var root_node : Node;

func _ready() -> void:
	# 假设 SignalBus 中定义了 ui_change 信号
	SignalBus.change_ui.connect(fui_change);

func register_root(node: Node):
	root_node = node
	
func fui_change(state) :
	print("jin")
	var	next_path : String;
	match state:
		"tologin":
			next_path = "res://game_ui/login/login_ui.tscn"
			print("成功")
			
	goto_ui(next_path);
	
func goto_ui(path: String):
	if current_ui:
		current_ui.queue_free()

	var ui = load(path)
	current_ui = ui.instantiate()
	root_node.add_child(current_ui)
