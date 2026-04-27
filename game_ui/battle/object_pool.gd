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
	if _pool.size() > 0:
		var card = _pool.pop_back()
		
		# 关键步骤：从 object_pool 的子节点列表中移除
		# 这样它就变成了一个“自由”节点，可以被 add_child 到其他地方
		if card.get_parent():
			card.get_parent().remove_child(card)
			
		card.visible = true
		return card
	else:
		# 如果池子空了，临时创建一个（防止报错）
		return create_new_card()
		
func return_card(card):
	if card.get_parent():
		card.get_parent().remove_child(card)
	

	if not card.is_inside_tree():
		add_child(card)
	_pool.append(card)
