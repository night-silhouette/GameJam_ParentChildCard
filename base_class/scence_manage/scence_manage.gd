extends Node

var current_scene : Node;
var root_node : Node;

func register_root(node: Node):
	root_node = node

func goto_scene(path: String):
	if current_scene:
		current_scene.queue_free()

	var scene = load(path)
	current_scene = scene.instantiate()
	root_node.add_child(current_scene)
