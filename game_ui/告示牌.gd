extends TextureRect

@onready var lab = $Label
@export var show_duration: float = 3.0
@export var fade_duration: float = 0.3

var hide_timer: Timer

func _ready() -> void:
	visible = false
	modulate.a = 0.0
	
	hide_timer = Timer.new()
	hide_timer.wait_time = show_duration
	hide_timer.one_shot = true
	hide_timer.timeout.connect(_hide_notice)
	add_child(hide_timer)
	
	SignalBus.notice_updated.connect(_show_notice)

func _show_notice(msg: String) -> void:
	hide_timer.stop()
	
	lab.text = msg
	visible = true
	modulate.a = 0.0
	
	var tween = create_tween()
	tween.set_trans(Tween.TRANS_QUAD)
	tween.set_ease(Tween.EASE_OUT)
	tween.tween_property(self, "modulate:a", 1.0, fade_duration)
	
	hide_timer.start()

func _hide_notice() -> void:
	var tween = create_tween()
	tween.set_trans(Tween.TRANS_QUAD)
	tween.set_ease(Tween.EASE_IN)
	tween.tween_property(self, "modulate:a", 0.0, fade_duration)
	tween.finished.connect(_on_hide_finished)
	

func _on_hide_finished() -> void:
	visible = false
	SignalBus.notic_end.emit()
