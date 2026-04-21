extends Control
@export var object_pool : Node;
@export var card_vector: GridContainer;
	
# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	for i in range(10) :
		card_init(i,i)
		
func card_init(net_ID:int,ID:int):
	var card = object_pool.get_card()
	card.init_card(CardManager.querry_resoure_by_id(ID));
	card.net_ID = net_ID;
	card_vector.add_child(card);
	
func card_dead(crad):
	pass;
