extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	set_all_controls_ignore(self)


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass


func set_all_controls_ignore(node: Node):
	for child in node.get_children():
		if child is Control:
			child.mouse_filter = Control.MOUSE_FILTER_IGNORE
		set_all_controls_ignore(child)
