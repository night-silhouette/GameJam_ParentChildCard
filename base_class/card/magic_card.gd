extends "res://base_class/card/battle_card.gd"

var _dissolve_material: ShaderMaterial
var _burn_glow: ColorRect

func _ready() -> void:
	super._ready()
	_setup_burn_glow()
	_setup_dissolve_material()
	SignalBus.ani_skill_card_notify_enter.connect(play_burn_dissolve)


func _setup_burn_glow() -> void:
	_burn_glow = ColorRect.new()
	_burn_glow.name = "BurnGlow"
	_burn_glow.color = Color(1.0, 0.4, 0.05, 0.0)  # 烈焰橙，初始全透明
	_burn_glow.size = size * 1.4
	_burn_glow.position = -Vector2(size.x * 0.2, size.y * 0.2)
	_burn_glow.mouse_filter = Control.MOUSE_FILTER_IGNORE
	_burn_glow.visible = false
	add_child(_burn_glow)
	move_child(_burn_glow, 0)  # 放到最底层，在卡牌纹理下面


func _setup_dissolve_material() -> void:
	_dissolve_material = ShaderMaterial.new()
	_dissolve_material.shader = Shader.new()
	_dissolve_material.shader.code = """
shader_type canvas_item;

uniform float dissolve_amount : hint_range(0.0, 1.0) = 0.0;
uniform float edge_width : hint_range(0.0, 0.5) = 0.12;
uniform vec4 burn_color : source_color = vec4(1.0, 0.2, 0.0, 1.0);
uniform sampler2D noise_tex : repeat_enable;

void fragment() {
	vec4 tex = texture(TEXTURE, UV);
	float noise = texture(noise_tex, UV).r;
	
	float edge = dissolve_amount + edge_width;
	
	if (noise < dissolve_amount) {
		// 已消融区域：丢弃
		discard;
	} else if (noise < edge) {
		// 燃烧边缘：渐变混合
		float t = (noise - dissolve_amount) / edge_width;
		// 边缘发光：越靠近消融前线越亮
		float glow = 1.0 - t;
		vec4 burn = mix(burn_color * (1.0 + glow * 0.8), tex, t);
		COLOR = burn;
	} else {
		// 未消融区域：正常显示
		COLOR = tex;
	}
}
"""
	_dissolve_material.set_shader_parameter("dissolve_amount", 0.0)
	_dissolve_material.set_shader_parameter("edge_width", 0.12)
	_dissolve_material.set_shader_parameter("burn_color", Color(1.0, 0.2, 0.0, 1.0))

	# 生成噪声纹理（消融图案）
	var noise_tex = NoiseTexture2D.new()
	var noise = FastNoiseLite.new()
	noise.noise_type = FastNoiseLite.TYPE_SIMPLEX
	noise.fractal_octaves = 3
	noise.frequency = 0.015
	noise_tex.noise = noise
	noise_tex.width = 256
	noise_tex.height = 256
	_dissolve_material.set_shader_parameter("noise_tex", noise_tex)


## 触发：燃烧消融动画
## duration: 动画总时长（秒），默认 1.2
func play_burn_dissolve(duration: float = 0.7) -> void:
	# 锁定状态，动画期间禁止交互
	set_change_lock(true)

	# 挂 shader 材质
	display.material = _dissolve_material
	_dissolve_material.set_shader_parameter("dissolve_amount", 0.0)

	# 显示背光
	_burn_glow.visible = true
	_burn_glow.modulate.a = 0.0

	var tween = create_tween()
	tween.set_parallel(true)

	# 背光：快速亮起 → 慢慢消退
	tween.tween_property(_burn_glow, "modulate:a", 0.95, duration * 0.2)
	tween.tween_property(_burn_glow, "modulate:a", 0.0, duration * 0.8).set_delay(duration * 0.2)

	# 消融：dissolve_amount 从 0 到 1
	tween.tween_method(_set_dissolve, 0.0, 1.0, duration)

	# 完成后清理
	tween.chain().tween_callback(_on_dissolve_finished)


func _set_dissolve(value: float) -> void:
	_dissolve_material.set_shader_parameter("dissolve_amount", value)


func _on_dissolve_finished() -> void:
	_burn_glow.visible = false
	display.material = null       # 移除 shader
	set_change_lock(false)        # 解锁
	visible = false               # 隐藏卡牌
	SignalBus.ani_end.emit()      # 通知 ani_state_machine 动画结束
