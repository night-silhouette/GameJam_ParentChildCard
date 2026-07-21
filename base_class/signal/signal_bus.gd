extends Node
#region Game
# 场景 / UI
signal change_scence(path: String)
signal change_ui(name: String)

# 网络状态（给UI用）
signal network_disconnected()

# 按钮 / 游戏行为
signal battle_information
signal online_match
#endregion

#region HTTP
# =========================
# 请求（Outgoing）
# =========================
signal request_login(username: String, password: String)
signal request_validate_token()
signal request_get_user_self()
signal request_get_user_by_id(user_id: int)
signal request_get_user_by_name(name: String)
signal request_register_user(name: String, password: String)
signal request_update_user(id: int, name: String, password: String)
signal request_delete_user_self()
signal request_delete_user_by_id(id: int)
signal request_bag_card
signal request_card_random
signal request_get_self_gold
signal request_debug_addcard
signal request_sell_card(card_list:Array)

signal request_battle
signal request_loot

# 邮件
signal request_send_mail(id: int, txt: String)
signal request_get_mail_numberN()
signal request_get_mail(page: int)
signal request_delete_mail(data: Array[int])
signal request_post_friend_mail()


# 时间
signal request_get_time(time: int)

# =========================
# 响应（Incoming）
# =========================
signal raw_api_responded(api_name: String, method: int, code: int, data: Variant, msg: String)

signal login_success()
signal login_failed(msg: String)

signal token_validated_success()

signal user_info_fetched(id: int, user_name: String, is_admin: bool)

signal user_registered_success()
signal user_updated_success()
signal get_card_random(card_list)
signal get_card_bag(card_list)
signal get_self_gold(gold)
signal sell_card_success()
signal ifbattle(index:bool)

# 邮件
signal send_mail_success()
signal get_mail_numberN_success()
signal get_mail_success()
signal delete_mail_success()

# 时间
signal get_time_success(Tserver: int)
signal get_time_debug(T: int)
#endregion

#region WS
# =========================
# 连接
# =========================
signal to_connect_ws(body)
signal ws_connected
signal ws_disconnected
signal to_reconnect_to(body)
signal soft_reconnect()
# 原始数据（调试用）
signal raw_ws_responded(code, data, msg)

# =========================
# 战斗状态
# =========================
signal battle_started
signal battle_over

# 匹配

signal match_canceled
##进入战斗开始信号
signal match_success(t)

# =========================
# 卡牌数据
# =========================
##全局接口
signal self_inhand_updated(cards)
signal oppent_inhand_updated(cards)
signal bt_selfinfo_updated(cards)
signal bt_oppinfo_updated(cards)
signal energy_updated(energy_list)
signal weather_update(weather_num)
signal get_user_id(self_id:int,oppent:int)
signal get_round_num()

# 子卡牌
signal child_card_list_updated(child_cards)
signal active_child_card_start(t, child_list)
signal active_child_card_succeed
signal active_child_card_finish(selected_temp_id_list)

# 天气
signal select_weather_notify(is_more)
signal select_weather_start(t, weather_list)
signal select_weather_succeed(weather_data)

# 卡牌结算

signal card_calc_start 

signal skill_card_notify()
signal weather_notify
signal buff_notify
signal action_card_notify(action_data)
signal deploy_card_notify()

signal child_belong_change(action_data)
signal card_pos_change()
signal hp_change()
signal buff_change()
signal refresh_all(All_data)
signal card_calc_finish

signal interrupt_start(action_data,is_need)
signal interrupt_succeed
signal combat_action_success
signal opoffline
signal oponline
# 弃牌堆
signal discard_list_updated(cards)

#判定
##判定回合开始信号
signal judge_start(t)
signal judge_finish
signal judge_put

#法术牌
signal  magic_card_start(t)
signal  deploy_magic_success
signal  magic_card_finish
#战斗牌

##判定回合开始信号. 动画开始时间暂停
signal combat_start_success(t,is_win);
signal enter_free()


# =========================
# 请求（发给WS）
# =========================
signal request_cancel_match
signal request_get_self_cards_inhand
signal request_get_opponent_cards_inhand
signal request_over_battle
signal request_deploy_magic_card(card_id,card_temp_id)
signal request_deploy_parent_card(card_id,card_temp_id)
signal request_deploy_child_card(card_id,card_temp_id)
signal request_judge(judge_data)
signal request_get_combat_cards
signal request_end_animation
signal request_combat_finish
signal request_combat_movement(combat_list)
signal request_get_energy
signal request_get_child_card_list
signal request_select_weather(weather)
signal request_active_child_card(temp_id_list)
signal request_get_discard_list()
signal request_interrupt_select(temp_id_list)
signal request_get_weather()
signal request_overbattle()
signal request_reconnect_query()
# 调试
signal request_debug_time
signal request_debug_matchpool

#endregion

#region card
signal enter_freecard(temp_id,zone)
signal exit_freecard(temp_id)
signal detected_area(zone)
signal exit_area(zone)
signal enter_hover(temp_id)
signal exit_hover(temp_id)
signal finishTurn
signal enemy_card_deployed(card_id, pos)
signal ani_end()
signal card_use_dead_enter()
signal card_use_dead_exit()
signal set_change_lock(a:bool)


# 动画状态进入信号（ani_state_machine 进入每个状态时发出，供外部监听）
signal ani_skill_card_notify_enter()
signal ani_weather_notify_enter()
signal ani_buff_notify_enter()
signal ani_action_card_notify_enter(caller, acceptor, behavior)
signal ani_deploy_card_notify_enter(action_data)
signal ani_child_belong_change_enter(origin, object)
signal ani_card_pos_change_enter(object, temp_id)
signal ani_hp_change_enter(temp_id, category, value)
signal ani_buff_change_enter()

signal ani_over_battle()
#endregion

signal left_clicked(stuff_id: int,state: int)
signal right_clicked(stuff_id: int,state: int)	
signal notice_updated(msg:String)
