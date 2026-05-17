# ActionReceiver.gd (处理服务器发来的消息)
extends Node

func _ready():
	SignalBus.raw_ws_responded.connect(_handle_ws_data)
	
func _handle_ws_data(code: int, data: Variant, msg: String):
	if code != 0 or data == null: return
	if data is String: data = JSON.parse_string(data)

	var action_code = int(data.get("action_code", -1))
	var action_data = data.get("action_data", null)
	var predicate = int(data.get("predicates", 0)) 
	
	print("[WS 接收] -> ", NetDef.get_predicate_name(predicate), "：",action_code)
	
	_dispatch(action_code, action_data, predicate)

# 核心的分发器：明明白白写清楚每一个 Action 是怎么处理的
func _dispatch(action_code: int, action_data: Variant, predicate: int):
	match action_code:
		
		NetDef.Action.GET_SELF_CARDS:
			# 针对获取手牌，只有结果返回时才处理
			if predicate == NetDef.Predicate.RESULT:
				# 在这里你可以自由地做中间处理，比如数据转换、校验
				if action_data is Array:
					SignalBus.self_inhand_updated.emit(action_data)
					
				else:
					push_error("GET_SELF_CARDS 返回格式错误，期望 Array")
					
		NetDef.Action.GET_BT_INFO:
			if predicate == NetDef.Predicate.RESULT:
				var self_data = action_data.get("self");
				var opp_data = action_data.get("opponent");
				SignalBus.bt_selfinfo_updated.emit(self_data);
				SignalBus.bt_oppinfo_updated.emit(opp_data);
		
		NetDef.Action.START_BATTLE:
			if predicate == NetDef.Predicate.NOTIFY:
				SignalBus.battle_started.emit(action_code)
				
		NetDef.Action.DEPLOY_CARD:
			if predicate == NetDef.Predicate.NOTIFY:
				if action_data is Dictionary:
					var card_id = action_data.get("card_id", "")
					var pos = action_data.get("position", Vector2.ZERO)
					SignalBus.enemy_card_deployed.emit(card_id, pos)
			if predicate == NetDef.Predicate.QUERY:
				var t = action_data.state_wait_time;
				var where = action_data.where;
				match where:
					2:
						SignalBus.magic_card_start.emit(t)
			if predicate == NetDef.Predicate.SUCCEED:
				SignalBus.deploy_magic_success.emit();
			if predicate	 == NetDef.Predicate.FINISH:
				SignalBus.magic_card_finish.emit();
						
		NetDef.Action.CANCEL_MATCH:
			if predicate == NetDef.Predicate.RESULT:
				SignalBus.match_canceled.emit()
		NetDef.Action.MATCH_SUCCESS:
			if predicate == NetDef.Predicate.NOTIFY:
				var t = action_data.state_wait_time;
				SignalBus.match_success.emit(t);
		NetDef.Action.JUDGE:
			if predicate == NetDef.Predicate.QUERY:
				var t = action_data.state_wait_time;
				SignalBus.judge_start.emit(t);
			if predicate == NetDef.Predicate.FINISH:
				SignalBus.judge_finish.emit(action_data);
			if predicate == NetDef.Predicate.SUCCEED:
				SignalBus.judge_put.emit()
		NetDef.Action.COMBAT:
			var t = action_data.state_wait_time;
			if predicate == NetDef.Predicate.QUERY:
				SignalBus.combat_start_success.emit(t,1);
			if predicate == NetDef.Predicate.NOTIFY:
				SignalBus.combat_start_success.emit(t,0);
		_:
			# 未处理的 action
			push_warning("未处理的下发动作 -> ", NetDef.get_action_name(action_code))
