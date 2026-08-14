extends Control
@export var card_node : Control
var lab_string:String
# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	card_node.buff_change.connect(_buff_change)
	card_node.dead.connect(_dead)


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _buff_change(delta: float) -> void:
	visible = true
	lab_string = ""
	if card_node.buff_list.is_empty():
		return
	for i  in card_node.buff_list:
		var string = Global.BUFF_NAME[i.get("buff_id")]+"(" + String(i.get("buff_value")) + ")"+"*" +String(i.get("buff_stacks"))
		string = string + "\n"
		lab_string = lab_string + string
	$Label.text = lab_string
	
func _dead():
	visible = false
