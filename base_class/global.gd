extends Node
var BASE_URL =  "http://120.26.145.68:10086" #"http://120.26.145.68:5300"
# Global.gd 或 NetworkClient.gd
var token_save :bool = false;
var init_battle_time :int;

const ACTION_NAME = {
	1: "取消匹配",
	2: "获取自己的卡牌信息",
	3: "获取对手的卡牌信息",
	4: "获取场上的战斗信息",
	5: "结束战斗",
	6: "开始战斗",
	7: "部署一张牌",
	8: "战斗回合判断",
	9: "匹配成功",
	10: "动画结束",
	11: "执行战斗行动"
}
enum Predicates {
	EMPTY = 0,
	NOTIFY = 1,
	QUERY = 2,
	RESULT = 3,
	FINISH = 4,
	SUCCEED = 5
}
# 配置表格式：{ action_code: { predicate: [预期类型, 信号] } }
const ACTION = {
	CANCEL_MATCH = 1,
	GET_SELF_CARDS = 2,
	GET_OPPONENT_CARDS = 3,
	GET_BT_INFO = 4,
	OVER_BATTLE = 5,
	START_BATTLE = 6,
	DEPLOY_CARD = 7,
	JUDGE = 8,
	MATCH_SUCCESS = 9,
	ANIMATION_END = 10,
	COMBAT = 11,
}
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
