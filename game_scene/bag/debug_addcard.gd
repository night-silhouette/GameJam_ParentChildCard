extends TextureButton


# Called when the node enters the scene tree for the first time.


func _on_button_down() -> void:
	SignalBus.request_debug_addcard.emit()
	SignalBus.request_bag_card.emit()
	SignalBus.request_get_self_gold.emit()
