extends Node

@onready var button: TextureButton = $"转换"
@onready var active_card: HBoxContainer = $"卡牌显示/active_child_card"
@onready var op_card: Control = $"卡牌显示/敌方牌库"
@onready var weather_name_label: Label = $"weather_name"
@onready var card_manager: Node = $"数据层/card_manager"

var _showing_active: bool = false

func _ready() -> void:
	Global.fake_death(active_card)
	Global.fake_death(op_card)

	# 临时：测试阶段全部允许输入
	$"block/全局block".allow_input()
	$"block/战斗牌block".allow_input()
	$"block/法术牌block".allow_input()
	
	if card_manager:
		card_manager.UI_date_update.connect(_update_active_child_display)
		card_manager.weather_num_changed.connect(_on_weather_num_changed)


func _ws_disconnected() -> void:
	SignalBus.change_scence.emit("tomenu")
	SignalBus.change_ui.emit("tomenu")


func _on_转换_button_down() -> void:
	if _showing_active:
		Global.fake_death(active_card)
		Global.revive(op_card)
	else:
		Global.fake_death(op_card)
		Global.revive(active_card)
	_showing_active = not _showing_active


func _update_active_child_display() -> void:
	if card_manager == null:
		return
	var active_cards = card_manager.get_cards_by_zone(Global.ZONE_CARD.CHILD_ACTIVE)
	var children = active_card.get_children()
	var index = 0
	for child in children:
		if not child.has_method("setup"):
			continue
		if index < active_cards.size():
			child.setup(active_cards[index])
			child.visible = true
		else:
			child.visible = false
		index += 1


func _on_weather_num_changed(weather_num: int) -> void:
	if weather_name_label == null:
		return
	weather_name_label.text = Global.WEATHER_NAME.get(weather_num, "未知")
	weather_name_label.visible = true
