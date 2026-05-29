
extends Node


func _on_debug_action_button_down() -> void:
	##你在这里写就好了，就写成下面这样的函数，按钮直接触发
	pass;

#region 示例和调用本地函数
func _on_request_cancel_match():

	_send_to_server(NetDef.Action.CANCEL_MATCH, NetDef.Predicate.QUERY, null)

func _request_judge(judge_data):
	var action_data = {
		"judge_data" = judge_data
	}
	_send_to_server(NetDef.Action.JUDGE, NetDef.Predicate.RESULT, action_data)
	
func _send_to_server(action_code: int, predicate: int, action_data: Variant):
	# 调用你实际的 WebSocket 发送函数
	BattleWs.send_action(action_code,action_data,predicate)
#endregion
