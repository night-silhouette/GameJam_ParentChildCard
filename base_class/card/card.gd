extends Control
   # 所有可进入区域

enum CardState { IDLE, HOVERED, DRAGGING }
var current_state: CardState = CardState.IDLE
const HOVER_ALLOWED_ZONES = [Global.ZONE_CARD.DECK_ZONE,Global.ZONE_CARD.SPELL_ZONE] 
@onready var display: TextureRect = $"卡牌纹理"
var temp_id ;
var hp ;
var damage  ;
var buff_id ;
var zone ;
var spell_des ;

# 拖拽判定距离（防止误触，按住鼠标移动超过这个像素才算拖拽，不需要可以设为 0）
const DRAG_THRESHOLD = 5.0
var mouse_start_pos = Vector2.ZERO


var dragging = false
var click_offset = Vector2.ZERO
var original_position = Vector2.ZERO



func _ready():
	mouse_entered.connect(_on_mouse_entered)
	mouse_exited.connect(_on_mouse_exited)
	

func update_card_data(base_res) -> void:
	display.texture = base_res.texture
	if base_res.is_combat_card :
		hp = base_res.hp;
		damage = base_res.damage;
	temp_id = base_res.temp_id;
	buff_id = base_res.buff_id;
	zone = base_res.zone;
	spell_des = base_res.spell_des;
func update():
	pass;

func free_card():
	display.texture = null;
	hp = null;
	damage = null;
	temp_id = null;
	buff_id = null;
	
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
	# 鼠标按下
	if event is InputEventMouseButton and event.button_index == MOUSE_BUTTON_LEFT:
		if event.pressed:
			change_state(CardState.DRAGGING)
	
