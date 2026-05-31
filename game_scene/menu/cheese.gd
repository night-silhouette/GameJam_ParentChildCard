extends TextureRect
var press_case :int = 0;
@export var _animationplayer : AnimationPlayer



func spawn_piece():
	var piece = preload("res://game_scene/menu/cheese_piece.tscn").instantiate()
	
	# 获取当前节点的全局中心位置
	var central_position = global_position + (size * scale / 2) 

	piece.global_position = central_position
	
	# 稍微给生成点一点偏移，让它们不要完全重合在一个点
	piece.global_position += Vector2(randf_range(-20, 20), randf_range(-20, 20))
	piece.scale = scale

	# 建议添加到当前场景，或者专门的特效层
	get_tree().current_scene.add_child(piece)
	
	
func animationplayer():
	$cheese_animationpalyer.play("抖动");
	_animationplayer.play("手柄进入")
	await _animationplayer.animation_finished
	$"../游戏开始".visible = true	
	$"../牌库".visible = true	

		
		
		
		
