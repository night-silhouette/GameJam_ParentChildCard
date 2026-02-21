extends Node

var peer 
var port 

func _ready() -> void:
	if !OS.has_feature("dedicated_server"):
		queue_free()
		return
	
	peer= WebSocketMultiplayerPeer.new()
	port = 443
	var err=peer.create_server(port)
	if err!=OK:
		print("服务器启动失败,错误码"+str(err))
	else:
		print("服务器启动success")
		
		multiplayer.multiplayer_peer=peer
		multiplayer.peer_connected.connect(_on_peer_connected)
		multiplayer.peer_disconnected.connect(_on_peer_disconnected)		
		
func _process(delta: float) -> void:
	if peer and peer.get_connection_status() !=MultiplayerPeer.CONNECTION_DISCONNECTED:
		peer.poll()

		
		
func _on_peer_connected(id):
	print(str(id)+"连接成功")
func _on_peer_disconnected(id):
	print(str(id)+"断开连接")
