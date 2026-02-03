@tool
extends TextureRect
var y = 324;
var x = 576;
@export var Duichen : TextureRect ;



func _process(delta: float) -> void:
	

	
	Duichen.position.x=2*x -  position.x;
	Duichen.size.x = size.x;
	
	Duichen.position.y =2*y - position.y;
	Duichen.size.y = size.y;
	
	
