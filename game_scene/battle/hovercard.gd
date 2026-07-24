extends card

# 定义传递给上级的信号
signal unhovered()

var stuff_id: int
var price: int
var zone : int





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
	zone = data.get("zone",0)
	
	
	var res: CardResource = data.get("resouce")
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

		
	# 2. 初始化缩放中心点
	pivot_offset = size / 2.0
	
	# 3. 初始化 Label 状态
	super.populate_detail_page($"卡牌具体页")
	# 💡 核心修改 2：初始化时确保“选中特效”是关闭的
	
# 快捷获取当前卡牌类型文本
func get_type_string() -> String:
	var sub_type = "子牌" if is_sub_card else "母牌"
	
	return sub_type


# 3. 监听鼠标悬停与离开事件
func _mouse_entered() -> void:
	_play_scale_tween(HOVER_SCALE)
	



func _mouse_exited() -> void:
	_play_scale_tween(NORMAL_SCALE)

	unhovered.emit()


# Tween 动画处理函数
func _play_scale_tween(target_scale: Vector2) -> void:
	if scale_tween and scale_tween.is_valid():
		scale_tween.kill()
	scale_tween = create_tween()
	scale_tween.tween_property(self, "scale", target_scale, TWEEN_DURATION)\
		.set_trans(Tween.TRANS_QUAD)\
		.set_ease(Tween.EASE_OUT)



					
func clear_data() -> void:
	# 1. 数据重置
	card_name = "未上牌"
	stuff_id = -1
	price = -1
	zone = -1
