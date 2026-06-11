extends GridContainer

## data_manager 引用
@export var card_manager: Node


func _ready() -> void:
	if card_manager == null:
		card_manager = get_node_or_null("../数据层/card_manager")
	
	if card_manager:
		card_manager.child_cards_ready.connect(_on_child_cards_ready)
		card_manager.interrupt_cards_ready.connect(_on_interrupt_cards_ready)


## 子牌堆数据就绪 → 填充 10 张 choose_card（match_code=3）
func _on_child_cards_ready() -> void:
	var children = get_children()
	if children.is_empty():
		return
	
	# 收集所有子牌 state zone 的卡（13~16）
	var child_cards: Array = []
	for z in range(Global.ZONE_CARD.CHILD_ACTIVE, Global.ZONE_CARD.CHILD_HAS_CATCH + 1):
		child_cards.append_array(card_manager.get_cards_by_zone(z))
	
	for i in range(children.size()):
		var card_node = children[i]
		if i < child_cards.size():
			var data = child_cards[i].duplicate()
			# 如果 data 没有 resouce，尝试从 card_manager 查询
			if data.get("resouce") == null:
				var card_id = int(data.get("id", -1))
				if card_id != -1:
					data["resouce"] = card_manager.querry_resoure_by_id(card_id)
			card_node.match_code = 3  # 子卡激活选择
			card_node.setup(data)
			card_node.visible = true
		else:
			card_node.visible = false


## 中断选牌数据就绪 → 填充对应 choose_card（match_code=99）
func _on_interrupt_cards_ready(card_list: Array, _select_num: int) -> void:
	var children = get_children()
	if children.is_empty():
		return
	
	for i in range(children.size()):
		var card_node = children[i]
		if i < card_list.size():
			var data = card_list[i]
			card_node.match_code = 99  # 中断选牌
			card_node.setup(data)
			card_node.visible = true
		else:
			card_node.visible = false
