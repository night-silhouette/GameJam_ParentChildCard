package protocol

type ProtocolCardWithCtx interface {
	ProtoColPush(e Effect)
	ProtoColReduceCardBtHp(SendTempId int, UserId int, TargetTempId int, ReduceHp float64) //死啦，会触发card的death
	ProtoColHealCardBt(UserId int, TargetTempId int, HealHp float64)                       //设置了不可以恢复到上限
	ProtoColSetDamageCardBt(UserId int, TargetTempId int, NewDamage float64)
}
