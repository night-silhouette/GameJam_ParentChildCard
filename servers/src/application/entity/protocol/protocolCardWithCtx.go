package protocol

import (
	"pcc_card/application/entity/BattleData"
	"time"
)

type ProtocolCardWithCtx interface {
	// Notify 传-1，表全部
	Notify(AnimationDto BattleData.AnimationDto, UserId int)
	ProtoColPush(e Effect)

	// ProtoColGetCharacterCard 获取所有手中的角色牌
	ProtoColGetCharacterCard(UserId int) []int
	// ProtoColSetCardBt 上牌(要BT上没牌才可以上)
	ProtoColSetCardBt(UserId int, TempId int)
	ProtoColReduceCardBtHp(SendTempId int, UserId int, TargetTempId int, ReduceHp float64) //死啦，会触发card的death//这些方法都是优先找BT然后找手牌,传递sendid是为了给死亡传递杀死者
	ProtoColHealCardBt(UserId int, TargetTempId int, HealHp float64)                       //设置了不可以恢复到上限
	ProtoColSetDamageCardBt(UserId int, TargetTempId int, NewDamage float64)
	ProtoColCardBtAttack(SendTempId int, UserId int, TargetTempId int, AtkHp float64)
	ProtoColInterrupt(UserId int, InterruptDto *BattleData.InterruptDto, res chan []int, InterruptWaitTime time.Duration) //异步中断，让前端从一定范围内需选牌，结果是tempId的数组，用res管道接受
	// ProtoColMoveDisCardPool 把卡删掉，并且移动到discardpool
	ProtoColMoveDisCardPool(UserId int, TempId int)
	GetBtCardInfo(id int) BattleData.BtCardInfo
	ProtoColCancelInterrupt()
}
