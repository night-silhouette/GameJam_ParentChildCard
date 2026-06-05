extends Control

var num :int = 0:
	set(value):
		num = value
		_update_lab()

@onready var lab = $Label

func _ready() -> void:
	_update_lab()

func _update_lab() -> void:
	lab.text = str(num) + "/5"
