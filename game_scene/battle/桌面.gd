extends TextureRect
@export var object_pool : Node
@export var card_manager: Node

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	$"母牌".object_pool = object_pool;
	$"母牌".card_manager = card_manager;
	$"子牌".object_pool = object_pool;
	$"子牌".card_manager = card_manager;
	$"敌方母牌".object_pool = object_pool;
	$"敌方母牌".card_manager = card_manager;
	$"敌方子牌".object_pool = object_pool;
	$"敌方子牌".card_manager = card_manager;
