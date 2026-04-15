extends Control
@export var object_pool : Node;

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
		pass;
func card_init(net_ID:int,ID:int,card):
	card.init_card(CardManager.querry_resoure_by_id(ID));
	card.net_ID = net_ID;
