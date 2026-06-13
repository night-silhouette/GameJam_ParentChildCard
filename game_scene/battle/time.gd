extends Node

@onready var timer: Timer = $Timer
@onready var label: Label = $Label


func _ready():
	timer.autostart = false
	timer.one_shot = true
	timer.timeout.connect(_on_timeout)
	timer.stop()
	label.visible = false


func start_countdown(server_duration: int):
	if server_duration <= 0:
		timer.stop()
		label.visible = false
		return

	timer.stop()
	timer.wait_time = float(server_duration - Global.max_delay_time)
	timer.start()
	label.visible = true


func _process(_delta: float):
	if not timer.is_stopped():
		var remaining = timer.time_left
		_update_label(remaining)
		_play_tick_effect(remaining)
	else:
		label.visible = false


func _on_timeout():
	label.visible = false
	SignalBus.enter_free.emit()


func _update_label(remaining: float):
	var display_seconds: int = ceil(remaining)
	var minutes: int = display_seconds / 60
	var seconds: int = display_seconds % 60
	label.text = "%02d:%02d" % [minutes, seconds]


func _play_tick_effect(remaining: float):
	if remaining <= 5.0 and remaining >= 0.0:
		label.pivot_offset = label.size / 2

		var tween = create_tween()
		label.scale = Vector2(1.5, 1.5)
		label.modulate = Color.RED

		tween.parallel().tween_property(label, "scale", Vector2.ONE, 0.2).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(label, "modulate", Color.WHITE, 0.2)
