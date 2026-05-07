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
signal to_connect_ws
signal ws_connected
signal ws_disconnected

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
signal self_cards_updated(cards)
signal opponent_cards_updated(cards)

#判定
##判定回合开始信号
signal judge_start(t)
signal judge_finish

#法术牌
signal  magic_card_start(t)
#战斗牌

##判定回合开始信号. 动画开始时间暂停
signal combat_start_success(t);
signal combat_start_fail(t);

# =========================
# 请求（发给WS）
# =========================
signal request_cancel_match
signal request_get_self_cards
signal request_get_opponent_cards
signal request_over_battle
signal request_deploy_magic_card(card_id,card_temp_id)
signal request_deploy_parent_card(card_id,card_temp_id)
signal request_deploy_child_card(card_id,card_temp_id)
signal request_judge(judge_data)

# 调试
signal request_debug_time
signal request_debug_matchpool
#endregion
