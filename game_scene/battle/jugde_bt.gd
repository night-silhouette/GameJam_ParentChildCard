extends Node

# 在检查器(Inspector)中直接拖入你 0, 1, 2 对应的粗线条手绘纹理
@export var texture_state_0 : Texture2D  
@export var texture_state_1 : Texture2D  
@export var texture_state_2 : Texture2D  

# 引用你的两个子节点 TextureRect 和胜负 Label
@onready var texture_rect_1 : TextureRect = $TextureRect
@onready var texture_rect_2 : TextureRect = $TextureRect2
@onready var label : Label = $Label

# 用于管理节点的独立 Tween，防止频繁连续赋值时动画打架冲突
var tween_1 : Tween
var tween_2 : Tween
var label_tween : Tween # 专门控制 Label 的动画器

## 定义带 setter 的变量。当外部执行 `judge_data = [...]` 时会自动触发下面的代码
var judge_data : Array = []:
	set(new_value):
		judge_data = new_value
		_refresh_ui_by_data()

## 核心需求：is_win 带 setter 控制
var is_win : int = -2:
	set(value):
		is_win = value
		_update_label_state(value) # 每次胜负数据改变，自动触发 Label 动效

func _ready() -> void:
	# 初始化时先刷新一次状态
	_refresh_ui_by_data()
	_update_label_state(is_win)


## 内部核心刷新函数
func _refresh_ui_by_data() -> void:
	if not is_inside_tree():
		await ready

	# 安全检查：如果数组为空或长度不够，默认全部隐藏处理
	if judge_data.size() < 2:
		_hide_node_instantly(texture_rect_1)
		_hide_node_instantly(texture_rect_2)
		return
		
	# 提取前两个元素
	var val_1 : int = judge_data[0]
	var val_2 : int = judge_data[1]
	
	# 分别更新两个子节点的纹理、显示状态与浮现动画
	_update_node_state(texture_rect_1, val_1, 1)
	_update_node_state(texture_rect_2, val_2, 2)

## 专门用来安全修改单个数据的方法
func update_single_judge_data(index: int, value: int) -> void:
	# 确保数组长度至少为 2
	while judge_data.size() < 2:
		judge_data.append(-1)
		
	# 安全赋值
	judge_data[index] = value
	
	# 💡 手动触发刷新（因为直接改索引不会走 set）
	_refresh_ui_by_data()
## 核心控制函数：处理 Label 的延迟、文字映射和华丽出现效果
func _update_label_state(value: int) -> void:
	if not is_instance_valid(label):
		return
		
	# 1. 如果为 -2，直接隐藏并杀掉残留动画，退出
	if value == -2:
		if label_tween: label_tween.kill()
		label.visible = false
		label.modulate.a = 0.0
		return
		
	# 2. 杀掉正在运行的旧 Label 动画，防止疯狂连击时画面闪烁
	if label_tween: 
		label_tween.kill()
		
	# 3. 映射对应的英文文本
	match value:
		0:
			label.text = "DRAW"
			label.modulate = Color.DARK_GRAY # 平局给个中性的灰色（可自行修改或删掉颜色修改）
		1:
			label.text = "WIN"
			label.modulate = Color.GREEN_YELLOW # 赢了给个亮眼的金绿/黄色
		-1:
			label.text = "LOSE"
			label.modulate = Color.CRIMSON # 输了给个醒目的猩红色
		_:
			label.visible = false
			return

	# 4. 【华丽动效配置】
	# 设置 Label 动画轴心为中心点（极其重要！确保你的 Label 在场景里的 Size 是正常的，否则回弹会歪）
	label.pivot_offset = label.size / 2
	
	# 设定初始状态：完全透明，且缩小到 0.5 倍（从深处弹出的感觉）
	label.modulate.a = 0.0
	label.scale = Vector2(0.5, 0.5)
	label.visible = true # 唤醒显示
	
	# 创建全新的独立 Tween
	label_tween = create_tween()
	
	# ─── 延时处理 ───
	# 延迟 0.4 秒。等上面的 texture_rect 差不多弹完，它再接上，视觉节奏感最好
	label_tween.tween_interval(0.4)
	
	# ─── 出现动画（并行：透明度+果冻回弹放大） ───
	label_tween.set_parallel(true)
	
	# A. 透明度淡入到 1.0
	label_tween.tween_property(label, "modulate:a", 1.0, 0.3)\
		.set_trans(Tween.TRANS_CUBIC).set_ease(Tween.EASE_OUT)
		
	# B. 大气回弹放大（从 0.5 弹到 1.0）
	# 使用 TRANS_BACK，它会故意放大到像 1.15 倍，然后像果冻一样缩回 1.0，效果非常动感
	label_tween.tween_property(label, "scale", Vector2.ONE, 0.45)\
		.set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)


## 核心控制函数：处理纹理、显隐切换以及“慢慢浮现”动画
func _update_node_state(target_rect: TextureRect, value: int, index: int) -> void:
	if not is_instance_valid(target_rect):
		return
		
	# 1. 如果输入是 -1，直接让它消失
	if value == -1:
		_hide_node_instantly(target_rect)
		return
		
	# 2. 杀掉该节点旧的未完成动画，防止连续赋值时画面疯狂闪烁
	if index == 1 and tween_1: tween_1.kill()
	if index == 2 and tween_2: tween_2.kill()
	
	# 3. 映射对应的手绘纹理
	match value:
		0: target_rect.texture = texture_state_0
		1: target_rect.texture = texture_state_1
		2: target_rect.texture = texture_state_2
		_:
			_hide_node_instantly(target_rect)
			return
	if index == 2:
		target_rect.rotation = 180;
	# 4. 【核心动效】：如果节点原本是隐藏的，或者刚刚被唤醒，触发慢慢浮现效果
	if not target_rect.visible or target_rect.modulate.a < 0.1:
		# 设置动画轴心为中心点（防止缩放时往左上角歪）
		target_rect.pivot_offset = target_rect.size / 2
		
		# 准备初始状态：完全透明，且缩小到 0.8 倍
		target_rect.modulate.a = 0.0
		target_rect.scale = Vector2(0.8, 0.8)
		target_rect.visible = true # 显式唤醒
		
		# 创建全新的 Tween 动画器
		var new_tween = create_tween().set_parallel(true) # 允许透明度和缩放同时进行
		
		# A. 透明度从 0 慢慢淡入到 1（耗时 0.3 秒）
		new_tween.tween_property(target_rect, "modulate:a", 1.0, 0.3)\
			.set_trans(Tween.TRANS_SINE).set_ease(Tween.EASE_OUT)
			
		# B. 尺寸从 0.8 平滑张开到 1.0
		new_tween.tween_property(target_rect, "scale", Vector2.ONE, 0.35)\
			.set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)
		
		# 保存引用以便后续清理
		if index == 1: tween_1 = new_tween
		if index == 2: tween_2 = new_tween
	else:
		# 如果节点本来就是可见的，只是换了个皮肤，那就直接显示，不重复放动画
		target_rect.visible = true
		target_rect.modulate.a = 1.0
		target_rect.scale = Vector2.ONE
	
	SignalBus.request_end_animation.emit();

# 内部辅助：瞬间隐藏节点并归零动画状态
func _hide_node_instantly(target_rect: TextureRect) -> void:
	if is_instance_valid(target_rect):
		target_rect.visible = false
		target_rect.modulate.a = 0.0
		target_rect.scale = Vector2(0.8, 0.8)
