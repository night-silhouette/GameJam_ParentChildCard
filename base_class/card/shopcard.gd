extends card

# 定义传递给上级的信号

signal unhovered()

var goods_id: int = 0
var card_id: int = 0
var price: int = 0
var zone: int = 0
var is_buy : bool

# 💡 新增：用于记录当前卡牌是否被选中 (true 为选中，false 为未选中)
var is_chosen: bool = false 


@onready var card: Control = $"卡牌具体页"



func _ready() -> void:
	mouse_entered.connect(_mouse_entered)
	mouse_exited.connect(_mouse_exited)

	
	
func setup(data: Dictionary) -> void:
	goods_id = int(data.get("goods_id", 0))
	card_id = int(data.get("card_id", 0))
	price = int(data.get("price", 0))
	zone = int(data.get("zone", 0))

	
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
		
	# 💡 核心修改 1：让 ui_name 节点直接显示卡牌的名字
		
	# 2. 初始化缩放中心点
	pivot_offset = size / 2.0
	# 3. 初始化 Label 状态
	# 💡 核心修改 2：初始化时确保"选中特效"是关闭的
	super.populate_detail_page(card)

# 3. hover抖动效果（位置不变）
const SHAKE_ANGLE := 3.0
const SHAKE_DURATION := 0.06

var _hover_tween: Tween


func _mouse_entered():
	if _hover_tween and _hover_tween.is_valid():
		_hover_tween.kill()
	_hover_tween = create_tween()
	_hover_tween.set_loops()
	_hover_tween.tween_property(self, "rotation", deg_to_rad(SHAKE_ANGLE), SHAKE_DURATION)
	_hover_tween.tween_property(self, "rotation", deg_to_rad(-SHAKE_ANGLE), SHAKE_DURATION)


func _mouse_exited():
	if _hover_tween and _hover_tween.is_valid():
		_hover_tween.kill()
	_hover_tween = create_tween()
	_hover_tween.tween_property(self, "rotation", 0.0, SHAKE_DURATION)


# 4. 监听鼠标点击事件
func _gui_input(event: InputEvent) -> void:
	if event is InputEventMouseButton and event.pressed:
		match event.button_index:
			MOUSE_BUTTON_LEFT:
				accept_event()
				SignalBus.buy_card.emit(goods_id)
