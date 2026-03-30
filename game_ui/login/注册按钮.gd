extends TextureButton
@export var username := LineEdit;
@export var password := LineEdit;


func _on_button_down() -> void:
	var _username:String = username.text;
	var _password:String = password.text;
	if validate_name(_username) && validate_password(_password) :
		SignalBus.request_register_user.emit(_username,_password);	
	
	
	
func validate_name(name: String) -> bool:
	var  index = 0;
	if name.length() < 3:
		print("用户名长度不得小于3")
		index = 1;
	if name.length() > 16:
		index = 1;
		print("用户名长度不得大于16")
	if index == 0 :
		return true;
	else :
		return false;
		
func validate_password(name: String) -> bool:
	var  index = 0;
	if name.length() < 6:
		print("密码长度不得小于3")
		index = 1;
	if name.length() > 25:
		index = 1;
		print("密码名长度不得大于16")
	if index == 0 :
		return true;
	else :
		return false;
