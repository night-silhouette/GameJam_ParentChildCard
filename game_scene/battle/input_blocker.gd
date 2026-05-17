extends Control
class_name InputBlocker

func _ready() -> void:
	# 初始状态：默认允许鼠标通过（不拦截）
	allow_input()


## 接口 1：不允许鼠标通过（开启拦截锁定）
func block_input() -> void:
	# MOUSE_FILTER_STOP 会消耗掉所有鼠标事件，使其无法传递给后面的节点
	mouse_filter = Control.MOUSE_FILTER_STOP
	
	# 如果你想在拦截时让游戏画面微微变暗（可选），可以取消注释下面这行：
	#modulate = Color(0, 0, 0, 0.3) 
	



## 接口 2：允许鼠标通过（解除拦截解锁）
func allow_input() -> void:
	# MOUSE_FILTER_IGNORE 让该节点完全对鼠标透明，事件会直接穿透过去
	mouse_filter = Control.MOUSE_FILTER_IGNORE
	
	# 恢复完全透明
	# modulate = Color(1, 1, 1, 0)
	
