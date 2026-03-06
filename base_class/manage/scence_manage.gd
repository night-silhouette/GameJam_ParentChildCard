extends Node

var current_scene : Node;
var root_node : Node;

func _ready() -> void:
	SignalBus.change_scence.connect(fchange_scence);

func register_root(node: Node):
	root_node = node
	
func fchange_scence(state) :
	var	next_path : String;
	match state:
		"start":
			next_path = "res://game_scene/main/start_scence.tscn"
		"tologin":
			next_path = "res://game_scene/login/login_scence.tscn"
			
	goto_scene(next_path);
	
func goto_scene(path: String):
	if current_scene:
		current_scene.queue_free()

	var scene = load(path)
	current_scene = scene.instantiate()
	root_node.add_child(current_scene)
