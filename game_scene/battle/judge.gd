extends Container # 如果你的父节点是 Control，这里就写 extends Control


const COLOR_NORMAL = Color.WHITE

const COLOR_PRESSED = Color(0.3, 0.3, 0.3)

var index_judge = -1;

func _ready() -> void:
	# 动态遍历父节点下的所有子节点
	for child in get_children():
		if child is BaseButton:
			# 1. 初始化所有按钮的变色状态
			_update_button_color(child)
			child.toggled.connect(_on_button_toggled.bind(child))
	init_all_unpressed()

func init_all_unpressed() -> void:
	# 1. 寻找任意一个子按钮，拿到它们共用的 ButtonGroup
	var group: ButtonGroup = null
	for child in get_children():
		if child is BaseButton and child.button_group != null:
			group = child.button_group
			break
	
	# 2. 如果找到了按钮组，用代码【强制允许全弹起】
	if group:
		group.allow_unpress = true
	
	# 3. 遍历所有按钮，强行解除按下状态，并刷回原色
	for child in get_children():
		if child is BaseButton:
			# 临时断开信号连接，防止重置状态时疯狂触发 _on_button_toggled 里的业务逻辑
			if child.toggled.is_connected(_on_button_toggled):
				child.toggled.disconnect(_on_button_toggled.bind(child))
			
			# 核心：状态归零
			child.button_pressed = false
			child.self_modulate = COLOR_NORMAL
			
			# 重新接回信号，不影响玩家后续的点击
			child.toggled.connect(_on_button_toggled.bind(child))
			

func _on_button_toggled(is_pressed: bool, button_node: BaseButton) -> void:
	if is_pressed and button_node.button_group:
		button_node.button_group.allow_unpress = false
	# 1. 刷新当前这个按钮的颜色
	_update_button_color(button_node)
	
	# 2. 如果这个按钮是被【按下】了，执行对应的游戏逻辑
	if is_pressed:
		_switch_game_logic(button_node.name)


# 统一控制按钮变色的内部函数
func _update_button_color(button_node: BaseButton) -> void:
	if button_node.button_pressed:
		# 选中的按钮变暗
		button_node.self_modulate = COLOR_PRESSED
	else:
		# 没选中的按钮恢复原色
		button_node.self_modulate = COLOR_NORMAL


# 供你扩展的卡牌游戏切换逻辑
func _switch_game_logic(button_name: String) -> void:
	match button_name:
		"剪刀":
			index_judge = 0
		"包袱":
			index_judge = 1
		"石头":
			index_judge = 2
