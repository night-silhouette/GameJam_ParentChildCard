extends card
   # 所有可进入区域

enum CardState { IDLE, HOVERED, DRAGGING }
var current_state: CardState = CardState.IDLE
const HOVER_ALLOWED_ZONES = [Global.ZONE_CARD.DECK_ZONE,Global.ZONE_CARD.SPELL_ZONE] 
@onready var display: TextureRect = $"卡牌纹理"
var temp_id: int = 0
var hp: int = 0
var damage: int = 0
var buff_id: Array = []
var zone: int = 0

# 拖拽判定距离（防止误触，按住鼠标移动超过这个像素才算拖拽，不需要可以设为 0）
const DRAG_THRESHOLD = 5.0
var mouse_start_pos = Vector2.ZERO


var is_selectable: bool = false
signal card_selected(temp_id: int)



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
	#   "buff_id": Array,    # Buff列表
	#   "zone": int,         # 显示区域
	#   "child_state": int,   # 子牌状态（可选）
	#   "resouce": CardResource  # 本地配置资源
	# }
	
	var res: CardResource = base_res.get("resouce")  # 注意拼写 "resouce"
	if res == null: return
	
	# 1. 先填充卡牌纹理（从 CardResource 获取）
	display.texture = res.card_texture
	
	# 2. 填充运行时数据（从 Dictionary 获取）
	temp_id = base_res.get("temp_id")
	zone = base_res.get("zone")
	buff_id = base_res.get("buff_id", [])
	
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
	buff_id = [];
	
func change_state(new_state: CardState):
	if current_state == new_state:
		return
	
	# 退出旧状态
	match current_state:
		CardState.HOVERED:
			SignalBus.exit_hover.emit()
		CardState.DRAGGING:
			pass

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
	if event is InputEventMouseButton and event.button_index == MOUSE_BUTTON_LEFT:
		if event.pressed:
			if is_selectable:
				card_selected.emit(temp_id)
			else:
				change_state(CardState.DRAGGING)
	
