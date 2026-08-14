extends Area2D
@export var zone : int

func _ready() -> void:
	collision_layer = 1
	collision_mask = 1
	monitoring = true
	monitorable = true

func _on_area_entered(area: Area2D) -> void:

	SignalBus.detected_area.emit(zone)


func _on_area_exited(area: Area2D) -> void:

	SignalBus.exit_area.emit(zone)
