package protocol

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/presentation/handler/battlehandler/BattleDto"
	"time"
)

type ProtocolCardWithCtx interface {
	// Notify 传-1，表全部
	Notify(AnimationDto BattleData.AnimationDto, UserId int)
	ProtoColPush(e Effect)
	GiveBuff(TempId int, buff *Buff)
	// ProtoColGetCharacterCard 获取所有手中的角色牌
	ProtoColGetCharacterCard(UserId int) []int
	// ProtoColSetCardBt 上牌(要BT上没牌才可以上)
	ProtoColSetCardBt(UserId int, TempId int)
	ProtoColReduceCardBtHp(SendTempId int, TargetTempId int, ReduceHp float64) //死啦，会触发card的death//这些方法都是优先找BT然后找手牌,传递sendid是为了给死亡传递杀死者
	ProtoColHealCardBt(TargetTempId int, HealHp float64)                       //设置了不可以恢复到上限
	ProtoColSetDamageCardBt(UserId int, TargetTempId int, NewDamage float64)
	ProtoColCardBtAttack(SendTempId int, UserId int, TargetTempId int, AtkHp float64, Category BattleData.ValueChange)
	ProtoColInterrupt(UserId int, InterruptDto *BattleData.InterruptDto, res chan []int, InterruptWaitTime time.Duration) //异步中断，让前端从一定范围内需选牌，结果是tempId的数组，用res管道接受
	ProtoColMoveDisCardPool(UserId int, TempId int)                                                                       // ProtoColMoveDisCardPool 把卡删掉，并且移动到discardpool
	GetBtCardInfo(id int) BattleData.BtCardInfo
	ProtoColCancelInterrupt()
	ProtoNotifyValue(Category BattleData.ValueChange, Value float64, TempId int, IsMiss bool) //数值变化的时候的通知
	ProtoNotifyCardMove(Object BattleData.Where, TempId int)
	ProtoColUpdateEnergy(UserId int, offset int)
	ProtoColCanUpdateEnergy(UserId int, offset int) bool
	CheckCard(id int) bool                                                           //检查是否还有卡出战,看pbt和cbt这两个位置主要是
	CreateTempId() int                                                               //产生tempId,计数用的
	ProtoColAttackNoHurt(CardTempId int, Value int, Category BattleData.ValueChange) //无主攻击
	ProtoColSetMaxHp(TargetTempId int, MaxHp float64)                                //设置最大生命
	GetIds() []int                                                                   //获取用户id数组
	GetWinnerIsAction() bool                                                         //给那个执悖天气用的
	GetWinnerId() int
	ProtoSendAction(UserId int, action BattleDto.Action)
	GetWeather() Weather
	GetDataAll(UseId int) *BattleData.DataAll
	Broad(v *CardMeta.BroadInfo)
}
