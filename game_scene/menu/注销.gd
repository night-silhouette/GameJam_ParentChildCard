extends Button


func _on_button_down() -> void:
	TokenManager.clear_token();
	SignalBus.change_scence.emit("tologin");
	SignalBus.change_ui.emit("tologin");
	
