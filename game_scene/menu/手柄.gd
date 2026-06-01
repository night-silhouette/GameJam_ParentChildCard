extends TextureRect
@export var animatedSprite2D:AnimatedSprite2D ;
@export var colorRect:ColorRect
func _ready() -> void:

	SignalBus.battle_over.connect(_battle_over)
	SignalBus.ws_disconnected.connect(_ws_disconnected)

func _on_笑脸按钮_button_down() -> void:
	SignalBus.change_ui.emit("tomatch")
	$"../ColorRect".visible = true


func _battle_over():
	print("退出战斗")
func _ws_disconnected():
	animatedSprite2D.visible = false;
	


func _on_笑脸按钮2_button_down() -> void:
	SignalBus.change_scence.emit("bag")
