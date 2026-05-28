extends card

# 定义传递给上级的信号
signal left_clicked(stuff_id: int)
signal right_clicked(stuff_id: int)

var stuff_id: int
var price: int

@onready var texture_rect: TextureRect = $TextureRect


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
	
	# 把剩下的战斗数据也填上，方便后续详细页读取
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
		
	# 2. 设置悬停提示 (Hover 核心)
	_setup_hover_tooltip()


# 格式化悬停显示的信息
func _setup_hover_tooltip() -> void:
	var type_str = "未知"
	if is_combat_card:
		type_str = "母牌 "
	elif is_sub_card:
		type_str = "子牌 "
		
	# 赋值给 Control 节点自带的 tooltip_text，鼠标悬停时会自动弹出
	tooltip_text = "名称: %s\n类型: %s\n价值: %d" % [card_name, type_str, value]


# 3. 监听鼠标点击事件 (点击核心)
func _gui_input(event: InputEvent) -> void:
	# 确保是鼠标点击事件，且是按下状态（防止弹起时重复触发）
	if event is InputEventMouseButton and event.pressed:
		match event.button_index:
			MOUSE_BUTTON_LEFT:
				# 必须消耗掉事件，防止点穿到下层 UI
				accept_event() 
				left_clicked.emit(stuff_id)
			MOUSE_BUTTON_RIGHT:
				accept_event()
				right_clicked.emit(stuff_id)
