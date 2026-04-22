extends TextureRect
@export var animatedSprite2D:AnimatedSprite2D ;
@export var colorRect:ColorRect
func _ready() -> void:
	SignalBus.match_success.connect(_match_success);
	SignalBus.battle_over.connect(_battle_over)
	SignalBus.ws_disconnected.connect(_ws_disconnected)

func _on_笑脸按钮_button_down() -> void:
	SignalBus.to_connect_ws.emit();
	animatedSprite2D.play()
	colorRect.visible = true;
	animatedSprite2D.visible = true;
	

func _match_success():
	SignalBus.change_scence.emit("tobattle");	
	SignalBus.change_ui.emit("tobattle");	
func _battle_over():
	print("退出战斗")
func _ws_disconnected():
	animatedSprite2D.visible = false;
	
