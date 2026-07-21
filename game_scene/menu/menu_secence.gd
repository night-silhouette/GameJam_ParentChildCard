extends Control

func _ready() -> void:
	SignalBus.ifbattle.connect(_ifbattle)
	SignalBus.request_battle.emit()
	
func _ifbattle(index):
	if index:
		SignalBus.request_reconnect.emit()
		SignalBus.match_success.emit()
