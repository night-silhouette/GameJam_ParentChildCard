shader_type canvas_item;

// 边框加深强度 (0~1)
uniform float border_strength : hint_range(0.0, 1.0) = 0.55;
// 中间去对比度 (0~3, 越大越灰白)
uniform float desaturate : hint_range(0.0, 3.0) = 0.65;
// 中间变淡 (0~3, 越大越淡)
uniform float fade_amount : hint_range(0.0, 3.0) = 0.7;
// 整体亮度 (0.5=暗, 1.0=原图, 1.5=过曝)
uniform float brightness : hint_range(0.5, 1.5) = 1.15;
// 边框宽度 (0~0.5)
uniform float border_width : hint_range(0.0, 0.5) = 0.15;

void fragment() {
	vec4 tex = texture(TEXTURE, UV);

	// 长方形边框距离：离边缘越近值越大 (0=中心, 1=边缘)
	vec2 edge = abs(UV - 0.5) * 2.0;
	float dist = max(edge.x, edge.y);

	// 中间：降对比度 → 偏灰
	float gray = dot(tex.rgb, vec3(0.299, 0.587, 0.114));
	tex.rgb = mix(tex.rgb, vec3(gray), desaturate * (1.0 - dist));

	// 中间：整体变淡变亮
	tex.rgb = mix(tex.rgb, vec3(0.9), fade_amount * (1.0 - dist));

	// 整体亮度
	tex.rgb *= brightness;

	// 长方形边框：加深加黑
	float inner = 1.0 - border_width;
	float border = smoothstep(inner, 1.0, dist) * border_strength;
	tex.rgb *= (1.0 - border);

	COLOR = tex;
}
