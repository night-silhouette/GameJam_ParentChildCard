extends card

# 定义传递给上级的信号


# 保留信号（如果上层其他地方不需要，也可以删掉，现在主控直接走本地 lab）
signal hovered(card_data: Dictionary)
signal unhovered()

var stuff_id: int
var price: int

@onready var texture_rect: TextureRect = $TextureRect
@onready var lab: Label = $"详细"

# 动画参数配置
var HOVER_SCALE := Vector2(1.05, 1.05) # 微微放大 1.05 倍
var NORMAL_SCALE := Vector2(1.0, 1.0)
var TWEEN_DURATION := 0.15 # 动画过渡时间（秒）
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
	value = res.value # 卡牌价值
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
		
	# 2. 初始化缩放中心点（核心：必须设置在中心，否则会往右下角放大）
	pivot_offset = size / 2.0
	
	# 3. 初始化 Label 状态：默认隐藏，防止一生成就挂着文字
	if lab:
		lab.visible = false


# 快捷获取当前卡牌类型文本
func get_type_string() -> String:
	if is_combat_card: return "母牌"
	if is_sub_card: return "子牌"
	return "未知"


# 3. 监听鼠标悬停与离开事件
func _mouse_entered() -> void:
	print("jianting")
	# 协调1：先放大卡牌
	_play_scale_tween(HOVER_SCALE)
	
	# 协调2：更新并显示本地的 Label
	if lab:
		lab.text = "名称: %s\n类型: %s\n价值: %d" % [card_name, get_type_string(), value]
		lab.visible = true
	
	# 向上层发出的信号（保留，供其他 UI 监听，不用可无视）
	hovered.emit({
		"name": card_name,
		"type": get_type_string(),
		"value": value,
		"price": price,
		"description": skill_description
	})


func _mouse_exited() -> void:
	# 协调1：卡牌缩回原样
	_play_scale_tween(NORMAL_SCALE)
	
	# 协调2：隐藏本地 Label 并清空文本
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
				SignalBus.right_clicked.emit(stuff_id)
