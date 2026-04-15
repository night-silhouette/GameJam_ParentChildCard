extends Node

# 预加载卡牌场景
const CARD_SCENE = preload("res://base_class/card/card.tscn")

# 存储“闲置”卡牌的数组（仓库）
var _pool = []

# 初始预热：在游戏开始或加载时调用
func _ready():
	# 比如预先生成 40 张牌
	pre_warm_pool(40)


func create_new_card():
	var card = CARD_SCENE.instantiate()
	
	card.visible = false
	return card

func pre_warm_pool(amount: int):
	for i in range(amount):
		var card = create_new_card();
		_pool.append(card)
		add_child(card)

func get_card():
	var card
	if _pool.size() > 0:
		card = _pool.pop_back() # 从仓库拿最后一张

	card.visible = true
	return card

func return_card(card):
	if card.get_parent():
		card.get_parent().remove_child(card)
	
	card.visible = false
	card.net_ID = null;
	
	if not card.is_inside_tree():
		add_child(card)
	_pool.append(card)
