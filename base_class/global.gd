extends Node
var BASE_URL = "http://120.26.145.68:5300" 
var token_save :bool = false;	
var init_battle_time :int;
var max_delay_time = 0.46;
var cardcalc_animaiton_list:Array;
var game_sence : int

# 注意：Action / Predicate / ACTION_NAME 等网络协议定义已统一迁移到 Net_def.gd
# 请使用 NetDef.Action / NetDef.Predicate 等，避免重复定义

const ZONE_CARD = {
	DECK_ZONE = 1,               # 牌堆
	
	DISCARD_ZONE = 2,       # 弃牌堆
	
	PARENT_BATTLE_ZONE = 3,   # 母牌战斗区
	CHILD_BATTLE_ZONE = 4,    # 子牌战斗区
	SPELL_ZONE = 5,         # 法术牌区
	
	FREE_ZONE = 6,
	
	ENEMY_PARENT_ZONE = 7,    # 敌方母牌战斗区
	ENEMY_CHILD_ZONE = 8,     # 敌方子牌区
	ENEMY_SPELL_ZONE = 9,    # 敌方法术牌区
	
	
	
	#战斗外：
	BAG_ZONE = 10,  #背包区域
	SELL_ZONE  = 11,  #卖出区域
	MATCH_ZONE = 12	, #出战区域
	
	
	# 子牌状态 zone（由 child_state 翻译而来，UI 统一按 zone 读取）
	CHILD_ACTIVE = 13,       # 子牌已激活
	CHILD_NOT_ACTIVE = 14,   # 子牌未激活
	CHILD_DIED = 15,         # 子牌已死亡（不进弃牌堆）
	CHILD_HAS_CATCH = 16,    # 子牌已被捕获
	
	ENEMY_HAND_ZONE = 17,    # 敌方手牌区
}
## 天气枚举
enum Weather {
	NINGJING = 0,
	SHABAO = 1,
	GANMU = 2,
	MINGYANG = 3,
	MIWU = 4,
	DASHU = 5,
	ZHIPAO = 6,
	SHUANGJIANG = 7,
	WUTU = 8,
	FENGDU = 9,
	UNABLE_VIEW = 9999,
}
## 天气中文名字典
const WEATHER_NAME = {
	Weather.NINGJING: "宁静",
	Weather.SHABAO: "沙暴",
	Weather.GANMU: "甘霖",
	Weather.MINGYANG: "明阳",
	Weather.MIWU: "迷雾",
	Weather.DASHU: "大暑",
	Weather.ZHIPAO: "执悖",
	Weather.SHUANGJIANG: "霜降",
	Weather.WUTU: "戊土",
	Weather.FENGDU: "酆都",
	Weather.UNABLE_VIEW: "无法查看",
}
## 天气描述字典
const WEATHER_DESC = {
	Weather.NINGJING: "无",
	Weather.SHABAO: "每回合结束，对所有上场的角色牌，造成一点真伤",
	Weather.GANMU: "每回合结束，对所有上场的角色牌，造成一点回血",
	Weather.MINGYANG: "每次回合判定结束之后，双方都获得1点行动点",
	Weather.MIWU: "双方攻击有15%的概率miss",
	Weather.DASHU: "双方每回合获得1层30%的治疗衰减",
	Weather.ZHIPAO: "回合判定胜者，如果不攻击或者使用技能会受到2点伤害",
	Weather.SHUANGJIANG: "每回合开始双方获得1层40%易伤",
	Weather.WUTU: "每回合开始双方获得1层28%虚弱",
	Weather.FENGDU: "战斗牌死亡时，会变成一只血量为1，伤害为1，没有技能的僵尸，僵尸死亡将不会再次变化",
	Weather.UNABLE_VIEW: "无法查看对方手牌和能量",
}
const WHERE={
	ParentCard = 0,
	ChildCard = 1,
	SkillCard = 2,
	DisCardPool = 4,
	ChildCardPool = 5,
	InHand = 6,
}

enum  BUFF  {
	BonusDamage = 1,
	Powerful = 2,
	Weakness = 3,
	DamageImmunity = 4,
	Vulnerability = 5,
	Block = 6,
	HealingBoost = 8,
	HealingDecay = 9,
	Wither = 10,
	Binding = 11,
	Retaliate = 12,
	Confine = 13,
	Giant = 14,
	Disarm = 15,
	XuFeng = 16,
}
const BUFF_NAME = {
	BUFF.BonusDamage : "额外伤害",
	BUFF.Powerful : "强盛",
	BUFF.Weakness : "虚弱",
	BUFF.DamageImmunity : "免伤",
	BUFF.Vulnerability : "易伤",
	BUFF.Block : "格挡",
	BUFF.HealingBoost : "治疗增强",
	BUFF.HealingDecay : "治疗衰减",
	BUFF.Wither : "凋零",
	BUFF.Binding : "束缚",
	BUFF.Retaliate : "反击",
	BUFF.Confine : "禁锢",
	BUFF.Giant : "巨人化",
	BUFF.Disarm : "缴械",
	BUFF.XuFeng : "续风",
}

enum ANI_BEHAVIOR {
	AnAttack = 0,
	AnHurt  = 1,
	AnDeath = 2,
	AnSkill = 3,
}
enum HP_CATEGORY  {
	Damage  = 0,
	Heal  = 1,
	TrueDamage  = 2,
}
enum GAME_SENCE{
	START = 0,
	LOGIN = 1,
	MENU  = 2,
	BATTLE = 3,
	BAG = 4,
}
## 接口 1：让节点进入【假死】状态
func fake_death(target_node: Node) -> void:
	if not is_instance_valid(target_node):
		return
		
	# 1. 掐断逻辑心跳：全面禁用该节点及其所有子节点的 _process, _physics_process 和 Timer
	target_node.process_mode = Node.PROCESS_MODE_DISABLED
	
	# 2. 掐断画面渲染：如果是 UI 节点或 2D 节点，直接隐藏
	if target_node is CanvasItem:
		target_node.visible = false
		
	# 3. 拦截鼠标与输入输入：
	# 如果是 UI 控件，让它彻底对鼠标透明（无法被点击）
	if target_node is Control:
		target_node.mouse_filter = Control.MOUSE_FILTER_IGNORE
	
	# 4. 强行关闭碰撞体（防止假死物体还能挡住网络弹道或卡牌拖拽判定）
	_set_all_collisions(target_node, false)

## 接口 2：把节点从假死中【救活】
func revive(target_node: Node) -> void:
	if not is_instance_valid(target_node):
		return
		
	# 1. 恢复画面显示
	if target_node is CanvasItem:
		target_node.visible = true
		
	# 2. 恢复鼠标接收
	if target_node is Control:
		target_node.mouse_filter = Control.MOUSE_FILTER_STOP # 或者是 MOUSE_FILTER_PASS，取决于你原厂设置
		
	# 3. 恢复物理碰撞
	_set_all_collisions(target_node, true)
	
	# 4. 最后一步：接回逻辑心跳。强制恢复和父节点一样的处理模式（通常是 INHERIT 恢复正常）
	target_node.process_mode = Node.PROCESS_MODE_INHERIT
	
## 辅助函数：递归开启/关闭节点下所有的物理碰撞体（如果有的话）
func _set_all_collisions(root_node: Node, enabled: bool) -> void:
	# 遍历目标节点下的所有子节点，把碰撞体全部禁用/启用
	for child in root_node.get_children():
		if child is CollisionShape2D or child is CollisionPolygon2D:
			child.set_deferred("disabled", !enabled)
		# 递归向下查找（比如卡牌内部嵌套的复合节点）
		if child.get_child_count() > 0:
			_set_all_collisions(child, enabled)
			
