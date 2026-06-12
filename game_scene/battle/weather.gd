extends Container

var index_judge = -1;
@onready var gou1 = $"1/勾"
@onready var gou2 = $"2/勾"
@onready var gou3 = $"3/勾"
@onready var gou4 = $"4/勾"

func _ready() -> void:
	for child in get_children():
		if child is  TextureButton:
			child.toggled.connect(_on_button_toggled.bind(child))
	init_all_unpressed()

func init_all_unpressed() -> void:
	for child in get_children():
		if child is TextureButton:
			child.set_block_signals(true)
			child.button_pressed = false
			_set_gou_visible(child, false)
			child.set_block_signals(false)
	index_judge = -1

func _on_button_toggled(is_pressed: bool, button_node: TextureButton) -> void:
	if is_pressed:
		button_node.button_group.allow_unpress = true
		_update_gou_visibility(button_node)
		_switch_game_logic(button_node.name)
	else:
		_set_gou_visible(button_node, false)
		index_judge = -1

func _set_gou_visible(button_node: TextureButton, visible: bool) -> void:
	match button_node.name:
		"1":
			gou1.visible = visible
		"2":
			gou2.visible = visible
		"3":
			gou3.visible = visible
		"4":
			gou4.visible = visible

func _update_gou_visibility(selected_button: TextureButton) -> void:
	for child in get_children():
		if child is TextureButton:
			if child.button_pressed:
				_set_gou_visible(child, true)
			else:
				_set_gou_visible(child, false)

func _switch_game_logic(button_name: String) -> void:
	match button_name:
		"1":
			index_judge = 0
		"2":
			index_judge = 1
		"3":
			index_judge = 2
		"4":
			index_judge = 3
