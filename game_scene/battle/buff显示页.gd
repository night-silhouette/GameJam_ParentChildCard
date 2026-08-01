extends Control
@export var card_node : Control

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	card_node.buff_change.connect(_buff_change)


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _buff_change(delta: float) -> void:
	for i  in card_node.buff_list:
	
