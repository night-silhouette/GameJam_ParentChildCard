extends card

#region 导出变量
## [核心] 选择模式指示码，由子类或场景中设置
@export var match_code: int = 0

## 数据管理器引用
@export var card_manager: Node = null
#endregion

#region 运行时变量
var temp_id: int = 0
var buff_id: Array = []
var zone: int = 0

## 选中状态（true=选中，false=未选中）
var is_chosen: bool = false:
	set(value):
		is_chosen = value
		_update_chosen_visual()
#endregion

#region 节点引用
@onready var texture_rect: TextureRect = $TextureRect
@onready var lab: Label = $"详细"
@onready var ui_name: Label = $name
@onready var check_mark: Control = $"勾"   # 选中勾号
#endregion

#region 动画参数
var HOVER_SCALE := Vector2(1.05, 1.05)
var NORMAL_SCALE := Vector2(1.0, 1.0)
var TWEEN_DURATION := 0.15
var scale_tween: Tween
#endregion


func _ready() -> void:
	mouse_entered.connect(_on_mouse_entered)
	mouse_exited.connect(_on_mouse_exited)
	
	# 如果 card_manager 未导出设置，尝试自动查找
	if card_manager == null:
		_auto_find_card_manager()
	
	# 监听数据层的选中状态变更
	if card_manager and card_manager.has_signal("selection_changed"):
		card_manager.selection_changed.connect(_on_selection_changed)


## 自动查找 card_manager（BattleScene 根节点下的数据管理器）
func _auto_find_card_manager() -> void:
	var root = get_tree().current_scene
	if root:
		card_manager = root.find_child("data_manager_BT", true, false)
		if card_manager == null:
			card_manager = root.find_child("card_manager", true, false)


## 初始化卡牌数据
## data 格式与 data_manager_BT.card_list 中条目一致：
## { "id": int, "temp_id": int, "zone": int, "resouce": CardResource, ... }
func setup(data: Dictionary) -> void:
	# 1. 提取运行时标识
	temp_id = data.get("temp_id", 0)
	zone = data.get("zone", 0)
	buff_id = data.get("buff_id", [])

	# 2. 获取本地配置资源
	var res: CardResource = data.get("resouce")
	if res == null:
		# 没有 resouce 时，尝试从 data 直接读取 card_name
		card_name = data.get("card_name", "")
		if ui_name:
			ui_name.text = card_name
		return

	# 3. 填充父类静态配置
	id = res.id
	card_name = res.name
	card_texture = res.card_texture
	value = res.value
	is_combat_card = res.is_combat_card
	is_sub_card = res.is_sub_card

	# 4. 填充战斗配置
	card_damage = res.damage
	initial_health = res.initial_health
	max_health = res.max_health
	skill_charge = res.skill_charge
	skill_card_use_num = res.skill_card_use_num
	skill_description = res.skill_description
	notes = res.notes
	sub_card_trigger_effect = res.sub_card_trigger_effect

	# 5. 更新纹理
	if texture_rect and card_texture:
		texture_rect.texture = card_texture

	# 6. 更新名称
	if ui_name:
		ui_name.text = card_name

	# 7. 初始化缩放中心
	pivot_offset = size / 2.0

	# 8. 初始化 Label
	if lab:
		lab.visible = false

	# 9. 初始化选中视觉效果
	_update_chosen_visual()


## 更新选中视觉效果（子类可重写）
func _update_chosen_visual() -> void:
	if is_instance_valid(check_mark):
		check_mark.visible = is_chosen


#region 鼠标交互

func _on_mouse_entered() -> void:
	_play_scale_tween(HOVER_SCALE)
	
	if lab:
		lab.text = " %s/ %d" % [_get_type_string(), value]
		lab.visible = true


func _on_mouse_exited() -> void:
	_play_scale_tween(NORMAL_SCALE)
	
	if lab:
		lab.visible = false
		lab.text = ""


## 左键点击：切换选中状态，通知数据层
func _gui_input(event: InputEvent) -> void:
	if event is InputEventMouseButton and event.pressed:
		match event.button_index:
			MOUSE_BUTTON_LEFT:
				accept_event()
				_on_left_click()

#endregion


## 左键点击处理：调用数据层 toggle_selection
func _on_left_click() -> void:
	if card_manager == null:
		return
	
	if not card_manager.has_method("toggle_selection"):
		push_warning("choose_card: card_manager 没有 toggle_selection 方法")
		return
	
	# 调用数据层的选中管理函数
	# 返回 true=选中, false=取消选中
	card_manager.toggle_selection(match_code, temp_id)


## 监听数据层的 selection_changed 信号，同步视觉效果
## 只响应属于自己 match_code 和 temp_id 的变更
func _on_selection_changed(changed_match_code: int, changed_temp_id: int, selected: bool) -> void:
	if changed_match_code != match_code:
		return
	if changed_temp_id != temp_id:
		return
	
	is_chosen = selected


#region 工具函数

func _get_type_string() -> String:
	return "子牌" if is_sub_card else "母牌"


func _play_scale_tween(target_scale: Vector2) -> void:
	if scale_tween and scale_tween.is_valid():
		scale_tween.kill()
	
	scale_tween = create_tween()
	scale_tween.tween_property(self, "scale", target_scale, TWEEN_DURATION)\
		.set_trans(Tween.TRANS_QUAD)\
		.set_ease(Tween.EASE_OUT)

#endregion
