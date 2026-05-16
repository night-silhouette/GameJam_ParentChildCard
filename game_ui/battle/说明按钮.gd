extends TextureButton

var flag = false;
func _on_button_down() -> void:
	flag = !flag;
	$"../说明框".visible = flag;
	$"../详细".visible = flag;
