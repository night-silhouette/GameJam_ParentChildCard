extends Control

@onready var grid: GridContainer = $GridContainer
var card_list
var is_refersh = true
var refresh_gold : int
func _ready() -> void:
	SignalBus.shop_get.connect(_shop_get)
	SignalBus.shop_post.connect(_shop_post)
	SignalBus.refresh_glod_get.connect(_refresh_glod_get)
	SignalBus.refresh_shop.connect(_refresh_shop)
	SignalBus.buy_card.connect(_buy_card)
	InventoryManager.gold_updated.connect(_on_gold_changed)
	
	SignalBus.request_shop_get.emit()
	SignalBus.request_refresh_glod_get.emit()
	SignalBus.request_get_self_gold.emit()
	await  SignalBus.shop_get
	SignalBus.request_shop_get.emit()
	$GridContainer.visible = true
	
func import_shop_data(data_list: Array) -> void:
	var children = grid.get_children()
	for i in range(children.size()):
		var shop_card_node = children[i]
		if i < data_list.size() and shop_card_node.has_method("set_up"):
			shop_card_node.set_up(data_list[i])
			# 复活
			shop_card_node.visible = true
			shop_card_node.modulate.a = 1.0
			shop_card_node.mouse_filter = Control.MOUSE_FILTER_STOP
			shop_card_node.process_mode = Node.PROCESS_MODE_INHERIT
		else:
			# 超出数据范围的假死
			shop_card_node.visible = false
			shop_card_node.modulate.a = 0.0
			shop_card_node.mouse_filter = Control.MOUSE_FILTER_IGNORE
			shop_card_node.process_mode = Node.PROCESS_MODE_DISABLED

func filter_by_goods(data_list: Array) -> void:
	var valid_ids: Array = []
	
	for data in data_list:
		valid_ids.append(int(data.get("goods_id", -1)))
	print("valid_ids: ",valid_ids)
	
	for child in grid.get_children():
		if not "goods_id" in child:
			continue
		if child.goods_id in valid_ids:
			continue
		# 不匹配的假死
		child.visible = false
		child.modulate.a = 0.0
		child.mouse_filter = Control.MOUSE_FILTER_IGNORE
		child.process_mode = Node.PROCESS_MODE_DISABLED


func _shop_get(data_list):
	card_list = data_list
	if is_refersh:
		import_shop_data(card_list)
		is_refersh = false
	else :
		filter_by_goods(data_list)
		
func _refresh_glod_get(gold):
	refresh_gold = gold
	$"转换/coin/num".text = str(refresh_gold)


func _on_gold_changed(value: int):
	$"金币显示/num".text = str(value)
	
func _shop_post():
	SignalBus.request_shop_get.emit()
	SignalBus.request_refresh_glod_get.emit()
	SignalBus.request_get_self_gold.emit()
	
func _refresh_shop():
	SignalBus.request_shop_get.emit()
	SignalBus.request_refresh_glod_get.emit()
	SignalBus.request_get_self_gold.emit()

func _buy_card(goods_id):
	for child in grid.get_children():
		if child.get("goods_id") == goods_id:
			if InventoryManager.gold < child.price:
				InventoryManager.notice_msg = "金币不足，无法购买"
				return
			break
	SignalBus.request_shop_post.emit(goods_id)
	
func _on_转换_button_down() -> void:
	if InventoryManager.gold < refresh_gold:
		InventoryManager.notice_msg = "金币不足，无法刷新"
		return
	is_refersh = true
	SignalBus.request_refresh.emit()
	


func _on_返回_button_down() -> void:
	SignalBus.change_scence.emit("tomenu")
