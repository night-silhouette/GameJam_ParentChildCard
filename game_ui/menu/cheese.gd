extends TextureRect


func spawn_piece():

	var piece = preload("res://game_ui/menu/cheese_piece.tscn").instantiate()

	piece.global_position = global_position

	get_tree().current_scene.add_child(piece);
	
	
func animationplayer():
	$cheese_animationpalyertion_Player.play("抖动");
