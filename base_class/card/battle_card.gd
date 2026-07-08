extends card
   # 所有可进入区域

enum CardState { IDLE, HOVERED, DRAGGING ,NEED_OPERATE}
var current_state: CardState = CardState.IDLE
const HOVER_ALLOWED_ZONES = [Global.ZONE_CARD.DECK_ZONE,Global.ZONE_CARD.SPELL_ZONE] 
@onready var display: TextureRect = $"卡牌纹理"
@onready var name_diapaly: Label = $name
@onready var buff_label: Label = $"buff"

var temp_id: int = 0
var hp: int = 0
var damage: int = 0
var buff_list: Array = []
var zone: int = 0

# 拖拽判定距离（防止误触，按住鼠标移动超过这个像素才算拖拽，不需要可以设为 0）
const DRAG_THRESHOLD = 5.0
var mouse_start_pos = Vector2.ZERO
var _is_pressing: bool = false


var is_selectable: bool = false
signal card_selected(temp_id: int)

## 中断选牌：是否被选中
var is_chosen: bool = false:
	set(value):
		is_chosen = value
		_update_chosen_visual()

@export var card_manager: Node


func _ready():
	mouse_entered.connect(_on_mouse_entered)
	mouse_exited.connect(_on_mouse_exited)
	

func update_card_data(base_res: Dictionary) -> void:
	# base_res 是 data_manager_BT 中的 card 字典，格式如下：
	# {
	#   "id": int,           # 卡牌配置ID
	#   "temp_id": int,      # 运行时唯一标识
	#   "hp": int,           # 当前生命值
	#   "damage": int,       # 当前伤害值
	#   "buff_list": Array[Dictionary],    # BuffDto数组 [{buff_id, buff_stacks, buff_value}]
	#   "zone": int,         # 显示区域
	#   "child_state": int,   # 子牌状态（可选）
	#   "resouce": CardResource  # 本地配置资源
	# }
	
	var res: CardResource = base_res.get("resouce")  # 注意拼写 "resouce"
	if res == null: return
	
	# 1. 先填充卡牌纹理（从 CardResource 获取）
	display.texture = res.card_texture
	
	# 显示卡牌名字
	name_diapaly.text = res.name
	
	# 2. 填充运行时数据（从 Dictionary 获取）
	temp_id = base_res.get("temp_id")
	zone = base_res.get("zone")
	buff_list = base_res.get("buff_list", [])
	_update_buff_display()
	
	# 战斗数据：如果这张卡在战场上，hp/damage 才有意义
	if res.is_combat_card:
		hp = base_res.get("hp", 0)
		damage = base_res.get("damage", 0)
	
	# 3. 填充父类静态配置数据（从 CardResource 获取）
	id = res.id
	card_name = res.name
	card_texture = res.card_texture
	value = res.value
	is_combat_card = res.is_combat_card
	is_sub_card = res.is_sub_card
	
	# 4. 填充战斗相关配置
	card_damage = res.damage
	initial_health = res.initial_health
	max_health = res.max_health
	skill_charge = res.skill_charge
	skill_card_use_num = res.skill_card_use_num
	skill_description = res.skill_description
	notes = res.notes
	sub_card_trigger_effect = res.sub_card_trigger_effect
	
func free_card():
	display.texture = null;
	hp = 0;
	damage = 0;
	temp_id = 0;
	buff_list = [];
	

func _update_buff_display() -> void:
	if not buff_label:
		return
	var lines: Array[String] = []
	for b in buff_list:
		if b is Dictionary:
			var bid = b.get("buff_id", 0)
			var name = Global.BUFF_NAME.get(bid, "???")
			var stacks = b.get("buff_stacks", 1)
			var value = b.get("buff_value")
			lines.append("%s(%d) x%d" % [name,value, stacks])
	buff_label.text = "\n".join(lines)
	buff_label.visible = not lines.is_empty()


func change_state(new_state: CardState):
	if current_state == new_state:
		return
	
	# 退出旧状态
	match current_state:
		CardState.HOVERED:
			SignalBus.exit_hover.emit()
		CardState.DRAGGING:
			pass
		CardState.NEED_OPERATE:
			is_chosen = false

	# 进入新状态
	current_state = new_state
	match current_state:
		CardState.IDLE:
			pass
			
			#var tween = create_tween()
			#tween.tween_property(self, "scale", Vector2.ONE, 0.1)
		
		CardState.HOVERED:
			SignalBus.enter_hover.emit(temp_id)
			#var tween = create_tween()
			#tween.tween_property(self, "scale", Vector2(1.05, 1.05), 0.1)
		
		CardState.DRAGGING:
			# 进入拖拽 → 发信号给自由节点
			SignalBus.exit_hover.emit()
			SignalBus.enter_freecard.emit(temp_id, zone)
			# ✅ 发完信号立刻回到 IDLE
			change_state(CardState.IDLE)
		CardState.NEED_OPERATE:
			pass

# --- 鼠标悬停事件（✅ 只有 zone 符合才允许 hover）---
func _on_mouse_entered():
	# ✅ 关键：只有 zone == HOVER_ENABLE_ZONE 时才允许进入悬停
	if current_state == CardState.IDLE and  zone in HOVER_ALLOWED_ZONES:
		change_state(CardState.HOVERED)

func _on_mouse_exited():
	if current_state == CardState.HOVERED:
		change_state(CardState.IDLE)

# --- 拖拽判定 ---
func _gui_input(event):
	# NEED_OPERATE 状态下：点击切换选中
	if current_state == CardState.NEED_OPERATE:
		if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
			accept_event()
			_on_need_operate_click()
		return

	if event is InputEventMouseButton and event.button_index == MOUSE_BUTTON_LEFT:
		if event.pressed:
			mouse_start_pos = event.global_position
			_is_pressing = true
		else:
			_is_pressing = false

	if _is_pressing and event is InputEventMouseMotion:
		var distance = event.global_position.distance_to(mouse_start_pos)
		if distance >= DRAG_THRESHOLD:
			_is_pressing = false
			change_state(CardState.DRAGGING)


## NEED_OPERATE 点击：切换选中状态，由 card_manager 校验上限
func _on_need_operate_click() -> void:
	if card_manager and card_manager.has_method("toggle_selection"):
		# card_manager 会校验上限，通过 selection_changed 信号回调结果
		card_manager.toggle_selection(99, temp_id)


## 外部信号：进入中断选牌模式
func enter_need_operate() -> void:
	change_state(CardState.NEED_OPERATE)


## 外部信号：退出中断选牌模式
func exit_need_operate() -> void:
	change_state(CardState.IDLE)


## 选中视觉反馈
func _update_chosen_visual() -> void:
	modulate = Color(1.0, 0.8, 0.5, 1.0) if is_chosen else Color.WHITE
