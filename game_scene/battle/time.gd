extends Node

# 定义完成时的信号，方便网络层或出牌逻辑监听
signal countdown_finished

# 默认 -1 代表关闭/空闲状态（float，0.1精度递减）
var countdown_time : float = 10.0

@onready var timer = $Timer
@onready var label = $Label

func _ready():	
	# 确保 Timer 初始不要自动运行，完全由我们用代码精确控制
	timer.wait_time = 0.1
	timer.one_shot = false
	timer.timeout.connect(_on_timer_timeout)
	
	# 初始化 UI 状态
	update_ui()

# 【核心核心】供外部（比如网络层收到服务器时间后）调用启动倒计时的函数
func start_countdown(server_duration: int):
	if server_duration <= 0:
		stop_and_hide()
		return
		
	countdown_time = server_duration
	label.visible = true # 显式唤醒 UI
	update_ui()
	
	timer.stop()             # 先停掉旧的
	timer.wait_time = 0.1    # 0.1秒精度
	timer.one_shot = false   # 确保它会循环触发，而不是只走一次！
	timer.start()            # 重新启动

# 【核心核心】将倒计时重置为 -1 并隐藏
func stop_and_hide():
	countdown_time = -1.0
	timer.stop()

# 0.1秒触发一次
func _on_timer_timeout():
	# 拦截：如果刚好在这个瞬间被外部重置为了 -1，直接停掉
	if countdown_time <= -1.0:
		stop_and_hide()
		return

	if countdown_time > 0.0:
		countdown_time -= 0.1
		update_ui()
		
		# 提前结束：剩余时间小于最大延迟阈值，直接进入自由阶段
		if countdown_time < Global.max_delay_time:
			stop_and_hide()
			SignalBus.enter_free.emit()
			return
		
		_play_tick_effect()
	else:
		# 到了 0 秒，停止倒计时，触发结束
		stop_and_hide() 
		_on_countdown_finished()

# 更新 UI 的方法
func update_ui():
	# 状态拦截：如果是 -1，直接隐藏，不进行任何字符换算
	if countdown_time <= -1.0:
		label.visible = false
		return
		
	# 如果不为 -1，确保它是可见的
	if not label.visible:
		label.visible = true
	
	var display_seconds : int = ceil(countdown_time)
	var minutes : int = display_seconds / 60
	var seconds : int = display_seconds % 60
	label.text = "%02d:%02d" % [minutes, seconds]

func _on_countdown_finished():
	label.text = "GO!" # 或者是 "时间到"
	# 发出信号，通知上层逻辑（例如：通知网络层“我本地时间到了，向服务器发强制过牌/空过包”）
	countdown_finished.emit()
	# print("【前端通知】倒计时已结束，信号已发出。")

# 在最后 5 秒文字扣减时调用心跳特效
func _play_tick_effect():
	if countdown_time <= 5 and countdown_time >= 0:
		label.pivot_offset = label.size / 2
		
		var tween = create_tween()
		label.scale = Vector2(1.5, 1.5)
		label.modulate = Color.RED # 顺便给你的最后5秒加个显眼的红光
		
		tween.parallel().tween_property(label, "scale", Vector2.ONE, 0.2).set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
		tween.parallel().tween_property(label, "modulate", Color.WHITE, 0.2)
		
