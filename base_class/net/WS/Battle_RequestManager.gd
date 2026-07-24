# ActionSender.gd (处理向服务器发送消息)
extends Node
	
func _ready():
	# 显式地将 UI 或逻辑层发出的请求信号，连接到具体的处理函数上
	SignalBus.request_cancel_match.connect(_on_request_cancel_match)
	SignalBus.request_get_self_cards_inhand.connect(_request_get_self_cards_inhand)
	# 假设带参数、
	SignalBus.request_judge.connect(_request_judge)
	SignalBus.request_deploy_magic_card.connect(_on_request_deploy_magic_card)
	SignalBus.request_deploy_parent_card.connect(_request_deploy_parent_card)
	SignalBus.request_deploy_child_card.connect(_request_deploy_child_card)
	SignalBus.request_get_combat_cards.connect(_request_get_combat_cards)
	SignalBus.request_end_animation.connect(_request_end_animation)
	SignalBus.request_combat_movement.connect(_request_combat_movement)
	SignalBus.request_over_battle.connect(_request_over_battle)
	# 新增请求信号连接
	SignalBus.request_get_energy.connect(_request_get_energy)
	SignalBus.request_get_child_card_list.connect(_request_get_child_card_list)
	SignalBus.request_select_weather.connect(_request_select_weather)
	SignalBus.request_active_child_card.connect(_request_active_child_card)
	SignalBus.request_get_discard_list.connect(_request_get_discard_list)
	SignalBus.request_interrupt_select.connect(_request_interrupt_select)
	SignalBus.request_get_weather.connect(_request_get_weather)
	SignalBus.request_get_opponent_cards_inhand.connect(_request_get_opponent_cards_inhand)
	SignalBus.request_reconnect_query.connect(_request_reconnect_query)
	
func _request_deploy_parent_card(card_id,card_temp_id):
	_on_deploy_card(0,card_id,card_temp_id);
	#print("card_id:",card_id,"temp_id:",card_temp_id);
	
func _request_deploy_child_card(card_id,card_temp_id):
	_on_deploy_card(1,card_id,card_temp_id);
	#print("card_id:",card_id,"temp_id:",card_temp_id);
	
func _on_request_deploy_magic_card(card_id,card_temp_id):
	_on_deploy_card(2,card_id,card_temp_id);
	
func _on_request_cancel_match():
	#print("[WS 发送] 取消匹配")
	_send_to_server(NetDef.Action.CANCEL_MATCH, NetDef.Predicate.QUERY, null)

func _request_get_self_cards_inhand():
	#print("[WS 发送] 获取手牌")
	_send_to_server(NetDef.Action.GET_SELF_CARDS, NetDef.Predicate.QUERY, null)
	
func _request_get_combat_cards():
	_send_to_server(NetDef.Action.GET_BT_INFO,NetDef.Predicate.QUERY,null);
	
func _request_end_animation():
	
	_send_to_server(NetDef.Action.ANIMATION_END,NetDef.Predicate.NOTIFY,null);	
	
func _request_get_weather():
	_send_to_server(NetDef.Action.GetWeather,NetDef.Predicate.QUERY,null);
func _request_over_battle():
	_send_to_server(NetDef.Action.OVER_BATTLE,NetDef.Predicate.NOTIFY,null)
func _request_get_opponent_cards_inhand():
	_send_to_server(NetDef.Action.GET_OPPONENT_CARDS, NetDef.Predicate.QUERY, null)
	# 带参数的特殊请求：这就是为什么不能用字典循环绑定的原因

func _on_deploy_card(where,card_id,card_temp_id):
	#print("[WS 发送] 部署卡牌: ", card_id)
	# 自定义中间组装逻辑
	var action_data = {
		"where": where,
		"card_id": card_id,
		"card_temp_id": card_temp_id
	}
	_send_to_server(NetDef.Action.DEPLOY_CARD, NetDef.Predicate.RESULT, action_data)	
	
func _request_judge(judge_data):
	var action_data = {
		"judge_data" = judge_data
	}
	_send_to_server(NetDef.Action.JUDGE, NetDef.Predicate.RESULT, action_data)
	
func _request_reconnect_query():
	_send_to_server(NetDef.Action.SoftReConnect,NetDef.Predicate.QUERY,null)
func _request_combat_movement(combat_list):
	var action_data = combat_list
	_send_to_server(NetDef.Action.COMBAT,NetDef.Predicate.RESULT,action_data)
		
func _request_get_energy():
	#print("[WS 发送] 查看能量值")
	_send_to_server(NetDef.Action.GetEnergy, NetDef.Predicate.QUERY, null)

func _request_get_child_card_list():
	#print("[WS 发送] 查看子牌堆")
	_send_to_server(NetDef.Action.GetChildCardList, NetDef.Predicate.QUERY, null)

func _request_select_weather(weather: int):
	#print("[WS 发送] 选择天气: ", weather)
	var action_data = {
		"weather": weather
	}
	#print(weather)
	_send_to_server(NetDef.Action.SelectWeather, NetDef.Predicate.RESULT, action_data)

func _request_active_child_card(temp_id_list: Array):
	#print("[WS 发送] 激活子卡牌: ", temp_id_list)
	var action_data = {
		"temp_id_list": temp_id_list
	}
	_send_to_server(NetDef.Action.ActiveChildCard, NetDef.Predicate.RESULT, action_data)

func _request_get_discard_list():
	#print("[WS 发送] 查看弃牌堆")
	_send_to_server(NetDef.Action.GetDisCard, NetDef.Predicate.QUERY, null)

func _request_interrupt_select(temp_id_list: Array):
	#print("[WS 发送] 中断选牌: ", temp_id_list)
	var action_data = {
		"temp_id_list": temp_id_list
	}
	_send_to_server(NetDef.Action.Interrupt, NetDef.Predicate.RESULT, action_data)
	
func _send_to_server(action_code: int, predicate: int, action_data: Variant):
	# 调用你实际的 WebSocket 发送函数
	BattleWs.send_action(action_code,action_data,predicate)
	
