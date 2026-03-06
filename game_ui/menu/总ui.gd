extends Control




func _ready():
	for b in $"按钮界面".get_children():
		print("发现子节点:", b.name);
		b.pressed.connect(_on_button_pressed.bind(b));

func _on_button_pressed(button):
	print("是这个按钮:", button.name);
