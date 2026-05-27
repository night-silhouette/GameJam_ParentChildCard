extends Node


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.request_get_self_cards_inhand.emit()
	SignalBus.request_get_combat_cards.emit()
	SignalBus.ws_disconnected.connect(_ws_disconnected)

func _ws_disconnected():
	SignalBus.change_scence.emit("tomenu");	
	SignalBus.change_ui.emit("tomenu");	
