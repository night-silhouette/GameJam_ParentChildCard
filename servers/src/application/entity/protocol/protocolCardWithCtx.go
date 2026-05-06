package protocol

import "pcc_card/application/entity/BattleData"

type ProtocolCardWithCtx interface {
	Notify(AnimationDto BattleData.AnimationDto)

	ProtoColReduceCardBtHp(SendTempId int, UserId int, TargetTempId int, ReduceHp float64) //死啦，会触发card的death//这些方法都是优先找BT然后找手牌
	ProtoColHealCardBt(UserId int, TargetTempId int, HealHp float64)                       //设置了不可以恢复到上限
	ProtoColSetDamageCardBt(UserId int, TargetTempId int, NewDamage float64)
	ProtoColCardBtAttack(SendTempId int, UserId int, TargetTempId int, AtkHp float64)
	ProtoColInterrupt(UserId int, InterruptDto BattleData.InterruptDto, res chan []int)
}
