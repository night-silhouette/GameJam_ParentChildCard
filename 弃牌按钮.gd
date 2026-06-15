extends Control
var index:int = 0
@export var dis_card :Control 

func _on_button_button_down() -> void:
	index += 1;
	if index%2 == 1:
		Global.revive(dis_card);
		dis_card.visible = false
	if index%2 == 0:
		Global.fake_death(dis_card)
		dis_card.visible = true
		index = 0
