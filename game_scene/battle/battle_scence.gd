extends Node


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.request_get_self_cards_inhand.emit()
	SignalBus.request_get_combat_cards.emit()
	
