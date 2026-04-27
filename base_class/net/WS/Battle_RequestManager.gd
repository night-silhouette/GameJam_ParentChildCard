# ActionSender.gd (处理向服务器发送消息)
extends Node

func _ready():
	# 显式地将 UI 或逻辑层发出的请求信号，连接到具体的处理函数上
	SignalBus.request_cancel_match.connect(_on_request_cancel_match)
	SignalBus.request_get_self_cards.connect(_on_request_get_self_cards)
	SignalBus.request_deploy_magic_card.connect(_on_request_deploy_magic_card) # 假设带参数、
	SignalBus.request_judge.connect(_request_judge)

# -----------------
# 具体的请求处理函数 (可以在这里封装不同结构的 action_data)
# -----------------
func _on_request_deploy_magic_card(card_id,card_temp_id):
	_on_deploy_card(2,card_id,card_temp_id);
func _on_request_cancel_match():
	print("[WS 发送] 取消匹配")
	_send_to_server(NetDef.Action.CANCEL_MATCH, NetDef.Predicate.QUERY, null)

func _on_request_get_self_cards():
	print("[WS 发送] 获取手牌")
	_send_to_server(NetDef.Action.GET_SELF_CARDS, NetDef.Predicate.QUERY, null)

# 带参数的特殊请求：这就是为什么不能用字典循环绑定的原因
func _on_deploy_card(where,card_id,card_temp_id):
	print("[WS 发送] 部署卡牌: ", card_id)
	
	# 自定义中间组装逻辑
	var action_data = {
		"where": where,
		"card_id": card_id,
		"card_temp_id": card_temp_id
	}
	_send_to_server(NetDef.Action.DEPLOY_CARD, NetDef.Predicate.QUERY, action_data)
func _request_judge(judge_data):
	var action_data = {
		"judge_data" = judge_data
	}
	_send_to_server(NetDef.Action.JUDGE, NetDef.Predicate.RESULT, action_data)

func _send_to_server(action_code: int, predicate: int, action_data: Variant):
	# 调用你实际的 WebSocket 发送函数
	BattleWs.send_action(action_code,action_data,predicate)
