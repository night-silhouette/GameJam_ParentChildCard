extends Control

func _ready() -> void:
	SignalBus.ifbattle.connect(_ifbattle)
	SignalBus.request_battle.emit()
	
func _ifbattle(index):
	if index:
		SignalBus.to_reconnect_to.emit()
