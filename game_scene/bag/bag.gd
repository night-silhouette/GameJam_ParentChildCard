extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.request_bag_card.emit();
	SignalBus.request_get_self_gold.emit()


func _on_返回_button_down() -> void:
	SignalBus.change_scence.emit("tomenu")


func _on_街机按钮_button_down() -> void:
	SignalBus.request_sell_card.emit(InventoryManager.get_cards_in_zone(Global.ZONE_CARD.SELL_ZONE))
