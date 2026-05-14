extends TextureRect
@export var object_pool : Node
@export var card_manager: Node

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	$"卡牌区".object_pool = object_pool;
	$"卡牌区".card_manager = card_manager
