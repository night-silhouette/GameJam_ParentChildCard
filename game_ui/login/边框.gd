extends TextureRect


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.login_success.connect(_to_menu)

	SignalBus.login_failed.connect(_fail)
	
	
func _to_menu():
	SignalBus.change_scence.emit("tomenu");
	SignalBus.change_ui.emit("tomenu");
func _fail(msg: String):
	# print("登陆失败：",msg);
	pass
