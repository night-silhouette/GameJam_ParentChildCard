extends TextureRect


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	position.x = 1152/2 - size.x*scale.x/2
	position.y = 648/2 - size.y*scale.y/2 -40
