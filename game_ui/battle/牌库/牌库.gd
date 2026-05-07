extends Control
@export var object_pool : Node;
@export var card_vector: GridContainer;
@export var card_manager: Node;
@export var animation_player : AnimationPlayer
# Called when the node enters the scene tree for the first time.
var cards: Array   # 你的全部牌
var page_size: int = 10
var start_index: int = 0

func _ready() -> void:
	card_manager.show_card.connect(refresh_ui)
	
func refresh_ui():
	cards = get_cards_by_zone(0)
	var page_cards = get_current_page()

	# 先清空原有卡牌节点
	for child in card_vector.get_children():
		child.free_card();
		object_pool.return_card(child);

	# 再生成当前页
	for card in page_cards:
		var child= object_pool.get_card()
		child.update_card(card);
		card_vector.add_child(child);

func get_current_page() -> Array:
	var end_index = min(start_index + page_size, cards.size())
	return cards.slice(start_index, end_index)

func get_cards_by_zone(zone: int) -> Array:
	var result = []
	for c in card_manager.card_list:
		if c.zone == zone:
			result.append(c)
	return result
	
func next_page():
	if start_index + page_size < cards.size():
		start_index += page_size
	
func prev_page():
	start_index = max(start_index - page_size, 0)


func _on_右切换_button_down() -> void:
	next_page();
	refresh_ui();
	


func _on_左切换_button_down() -> void:
	prev_page();
	refresh_ui();
