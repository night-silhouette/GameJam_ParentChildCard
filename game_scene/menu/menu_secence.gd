extends Control

func _ready() -> void:
	SignalBus.ifbattle.connect(_ifbattle)
	SignalBus.request_battle.emit()
	
func _ifbattle(index):
	if index:
		SignalBus.to_reconnect_to.emit()


func _on_万能按钮_button_down() -> void:
	SignalBus.request_shop_get.emit()
	SignalBus.request_refresh_glod_get.emit()
	var data : Array = await SignalBus.shop_get	
	var index = int(data.pop_front().get("goods_id"))
	SignalBus.request_shop_post.emit(index)
	print("goods_id: ",index)

	


func _on_街机按钮_button_down() -> void:
	SignalBus.change_scence.emit("toshop")
