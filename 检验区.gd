extends Area2D
@export var zone : int
@export var card_manage : Node

func _on_body_entered(body: Node2D) -> void:
	card_manage.free_card_enter.emit(zone);
