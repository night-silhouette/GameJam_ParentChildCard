extends Control

@onready var display: TextureRect = $Display
var hp : int;
var atk : int ;

# 每次刷新都会调用这个函数
func update_from_network(base_res: CardResource, net_data: Dictionary) -> void:
	# 1. 基础信息刷新（资源路径、名称等静态内容）
	if base_res and base_res.card_texture:
		display.texture = base_res.card_texture
	
