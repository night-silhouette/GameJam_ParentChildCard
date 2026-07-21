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
	
	#print("[WS 接收] -> ", NetDef.get_predicate_name(predicate), "：",NetDef.ACTION_NAME[action_code])
	
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
		NetDef.Action.GET_OPPONENT_CARDS:
			if predicate == NetDef.Predicate.RESULT:
				if action_data is Array:
					SignalBus.oppent_inhand_updated.emit(action_data)
					
		NetDef.Action.GET_BT_INFO:
			if predicate == NetDef.Predicate.RESULT:
				var self_data = action_data.get("self");
				var opp_data = action_data.get("opponent");
				SignalBus.bt_selfinfo_updated.emit(self_data);
				SignalBus.bt_oppinfo_updated.emit(opp_data);
				
				

		NetDef.Action.START_BATTLE:
			if predicate == NetDef.Predicate.NOTIFY:
				SignalBus.battle_started.emit(action_code)
				
		NetDef.Action.OVER_BATTLE:
			if predicate == NetDef.Predicate.RESULT:
				SignalBus.battle_over.emit(str(action_data))#action_data:winner/loser
		NetDef.Action.DEPLOY_CARD:
			if predicate == NetDef.Predicate.NOTIFY:
				if action_data is Dictionary:
					var action = SignalBus.deploy_card_notify.emit.bind(action_data)
					Global.cardcalc_animaiton_list.push_back(action)
					
			if predicate == NetDef.Predicate.QUERY:
				var t = action_data.state_wait_time;
				var where = int(action_data.where);
				# print(where)
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
				#print("init阶段")
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
			if predicate == NetDef.Predicate.QUERY:
				var t = action_data.state_wait_time;
				SignalBus.combat_start_success.emit(t,1);
			if predicate == NetDef.Predicate.NOTIFY:
				var t = action_data.state_wait_time;
				SignalBus.combat_start_success.emit(t,0);
			if predicate == NetDef.Predicate.SUCCEED:
				SignalBus.combat_action_success.emit()
				print("success!!!!!!!!!!!!!")
		NetDef.Action.OpOffline:
			if predicate == NetDef.Predicate.NOTIFY:
				SignalBus.opoffline.emit() 	#对手离线
		NetDef.Action.OpOffline:
			if predicate == NetDef.Predicate.NOTIFY:
				SignalBus.oponline.emit() 	#对手上线
		# 新增：能量值
		NetDef.Action.GetEnergy:
			if predicate == NetDef.Predicate.RESULT:
				SignalBus.energy_updated.emit(action_data)
		
		# 新增：子牌堆
		NetDef.Action.GetChildCardList:
			if predicate == NetDef.Predicate.RESULT:
				if action_data is Array:
					SignalBus.child_card_list_updated.emit(action_data)
					#print(action_data)
					
				else:
					push_error("GET_CHILD_CARD_LIST 返回格式错误，期望 Array")
		
		# 新增：激活子卡牌
		NetDef.Action.ActiveChildCard:
			if predicate == NetDef.Predicate.QUERY:
				var t = action_data.get("state_wait_time", 0)
				var child_list = action_data.get("child_list", [])
				SignalBus.active_child_card_start.emit(t, child_list)
			if predicate == NetDef.Predicate.SUCCEED:
				SignalBus.active_child_card_succeed.emit()
			if predicate == NetDef.Predicate.FINISH:
				var selected_list = action_data.get("selected_temp_id_list", [])
				SignalBus.active_child_card_finish.emit(selected_list)
		
		# 新增：选择天气
		NetDef.Action.SelectWeather:
			if predicate == NetDef.Predicate.NOTIFY:
				var is_more = action_data.get("is_more", false)
				SignalBus.select_weather_notify.emit(is_more)
			if predicate == NetDef.Predicate.QUERY:
				var t = action_data.get("state_wait_time", 0)
				var weather_list = action_data.get("weather_list", [])
				print(action_data)
				SignalBus.select_weather_start.emit(t, weather_list)
			if predicate == NetDef.Predicate.SUCCEED:
				SignalBus.select_weather_succeed.emit(action_data)
		
		# 新增：弃牌堆
		NetDef.Action.GetDisCard:
			if predicate == NetDef.Predicate.RESULT:
				if action_data is Array:
					SignalBus.discard_list_updated.emit(action_data)
				else:
					push_error("GET_DISCARD 返回格式错误，期望 Array")
					
		NetDef.Action.GetWeather:
			if predicate == NetDef.Predicate.RESULT:
				SignalBus.weather_update.emit(int(action_data))
		
		# 新增：卡牌结算
		NetDef.Action.CardCalc:
			if predicate == NetDef.Predicate.FINISH:
				SignalBus.card_calc_finish.emit()
			if predicate == NetDef.Predicate.NOTIFY:
				SignalBus.card_calc_start.emit()
				
		# 新增：中断选牌
		NetDef.Action.Interrupt:
			if predicate == NetDef.Predicate.QUERY:#需要进行操作的终端
				var t = action_data.get("state_wait_time")
				var temp_id_list = action_data.get("temp_id_list")#可以选的牌
				var select_num = action_data.get("select_num")#需要选几张
				var interrupt_type = action_data.get("interrupt_type")
				var call_temp_id = action_data.get("call_temp_id")
				SignalBus.interrupt_start.emit(action_data,1)

			if predicate == NetDef.Predicate.NOTIFY:#不需要操作的终端
				var t = action_data.get("state_wait_time")
				var temp_id_list = action_data.get("temp_id_list")#可以选的牌
				var select_num = action_data.get("select_num")#需要选几张
				var interrupt_type = action_data.get("interrupt_type")
				var call_temp_id = action_data.get("call_temp_id")
				SignalBus.interrupt_start.emit(action_data,0)			
			if predicate == NetDef.Predicate.SUCCEED:
				SignalBus.interrupt_succeed.emit()
				
		NetDef.Action.SkillCardNotify:
			var action = SignalBus.skill_card_notify.emit.bind()
			Global.cardcalc_animaiton_list.push_back(action)
		NetDef.Action.WeatherNotify:
			var action = SignalBus.weather_notify.emit.bind()
			Global.cardcalc_animaiton_list.push_back(action)
		NetDef.Action.BuffCalcNotify:
			var action = SignalBus.buff_notify.emit.bind()
			Global.cardcalc_animaiton_list.push_back(action)
		NetDef.Action.AnimationNotify:
			var action = SignalBus.action_card_notify.emit.bind(action_data)
			Global.cardcalc_animaiton_list.push_back(action)
		NetDef.Action.ChildBelongChange:
			var action = SignalBus.child_belong_change.emit.bind(action_data)
			Global.cardcalc_animaiton_list.push_back(action)
			_add_refresh(action_data)
		NetDef.Action.PositionChange:
			var action = SignalBus.card_pos_change.emit.bind(action_data)
			Global.cardcalc_animaiton_list.push_back(action)
			_add_refresh(action_data)
		NetDef.Action.HpChange:
			var action = SignalBus.hp_change.emit.bind(action_data)
			Global.cardcalc_animaiton_list.push_back(action)
			_add_refresh(action_data)
			
		NetDef.Action.BuffChange:
			_add_refresh(action_data)

			
		
		_:
			# 未处理的 action
			push_warning("New", NetDef.get_action_name(action_code),action_data,predicate)

func _add_refresh(action_data):
	var data_all = action_data.get("data_all")
	var action = SignalBus.refresh_all.emit.bind(data_all)
	Global.cardcalc_animaiton_list.push_back(action)
