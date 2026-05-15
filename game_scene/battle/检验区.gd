extends Area2D
@export var zone : int

func _on_area_entered(area: Area2D) -> void:
	SignalBus.detected_area.emit(zone)


func _on_area_exited(area: Area2D) -> void:
	SignalBus.exit_area.emit()
