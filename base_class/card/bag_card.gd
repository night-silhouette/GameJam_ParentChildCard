extends card

# 定义传递给上级的信号
signal hovered(card_data: Dictionary)
signal unhovered()

var stuff_id: int
var price: int

# 💡 新增：用于记录当前卡牌是否被选中 (true 为选中，false 为未选中)
var is_chosen: bool = false 

@onready var texture_rect: TextureRect = $TextureRect
@onready var lab: Label = $"详细"
@onready var ui_name: Label = $name
@onready var border: ReferenceRect = $Border

# 动画参数配置
var HOVER_SCALE := Vector2(1.05, 1.05) 
var NORMAL_SCALE := Vector2(1.0, 1.0)
var TWEEN_DURATION := 0.15 
var scale_tween: Tween

func _ready() -> void:
	mouse_entered.connect(_mouse_entered)
	mouse_exited.connect(_mouse_exited)
	
func setup(data: Dictionary) -> void:
	stuff_id = data.get("stuff_id", 0)
	price = data.get("price", 0)
	
	var res: CardResource = data.get("resource")
	if res == null: return
		
	# 填充父类数据
	id = res.id
	card_name = res.name
	card_texture = res.card_texture
	value = res.value 
	is_combat_card = res.is_combat_card
	is_sub_card = res.is_sub_card
	
	# 把剩下的战斗数据也填上
	card_damage = res.damage
	initial_health = res.initial_health
	max_health = res.max_health
	skill_charge = res.skill_charge
	skill_card_use_num = res.skill_card_use_num
	skill_description = res.skill_description
	notes = res.notes
	sub_card_trigger_effect = res.sub_card_trigger_effect
	
	# 1. 更新卡牌基本纹理
	if texture_rect and card_texture:
		texture_rect.texture = card_texture
		
	# 💡 核心修改 1：让 ui_name 节点直接显示卡牌的名字
	if ui_name:
		ui_name.text = card_name
		
	# 2. 初始化缩放中心点
	pivot_offset = size / 2.0
	
	# 3. 初始化 Label 状态
	if lab:
		lab.visible = false
		
	# 💡 核心修改 2：初始化时确保“选中特效”是关闭的
	is_chosen = false
	if border: border.visible = false

# 快捷获取当前卡牌类型文本
func get_type_string() -> String:
	if is_combat_card: return "母牌"
	if is_sub_card: return "子牌"
	return "未知"


# 3. 监听鼠标悬停与离开事件
func _mouse_entered() -> void:
	_play_scale_tween(HOVER_SCALE)
	
	# 💡 核心修改 3：将名称信息从详细描述中移除
	if lab:
		lab.text = "类型: %s\n价值: %d" % [get_type_string(), value]
		lab.visible = true
	
	hovered.emit({
		"name": card_name,
		"type": get_type_string(),
		"value": value,
		"price": price,
		"description": skill_description
	})


func _mouse_exited() -> void:
	_play_scale_tween(NORMAL_SCALE)
	
	if lab:
		lab.visible = false
		lab.text = ""
	
	unhovered.emit()


# Tween 动画处理函数
func _play_scale_tween(target_scale: Vector2) -> void:
	if scale_tween and scale_tween.is_valid():
		scale_tween.kill()
	
	scale_tween = create_tween()
	scale_tween.tween_property(self, "scale", target_scale, TWEEN_DURATION)\
		.set_trans(Tween.TRANS_QUAD)\
		.set_ease(Tween.EASE_OUT)


# 4. 监听鼠标点击事件
func _gui_input(event: InputEvent) -> void:
	if event is InputEventMouseButton and event.pressed:
		match event.button_index:
			MOUSE_BUTTON_LEFT:
				accept_event() 
				SignalBus.left_clicked.emit(stuff_id)
			
			MOUSE_BUTTON_RIGHT:
				accept_event()
				
				# 💡 核心修改 4：切换选中状态并计算 state
				is_chosen = !is_chosen # 取反：true 变 false，false 变 true
				var state = 1 if is_chosen else 0
				
				# 💡 核心修改 5：触发 _draw() 绘制选中框效果
				if border:
					border.visible = is_chosen

				# 发送带有最新状态的信号
				SignalBus.right_clicked.emit(stuff_id, state)
