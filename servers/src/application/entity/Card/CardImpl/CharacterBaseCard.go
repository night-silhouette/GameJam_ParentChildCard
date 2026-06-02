package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
	"pcc_card/global"
	"time"
)

type CharacterBaseCard struct {
	BaseCard
	Card CardAbstract.Card
}

func (c CharacterBaseCard) Attack(TargetId int) {
	c.Notify(BattleData.AnAttack, -1)

	c.EffectAttack(TargetId, c.AtkNow)
}

func (c CharacterBaseCard) Hurt(AttackId int, HurtHp float64) {
	c.Notify(BattleData.AnHurt, -1)
	c.EffectHurt(AttackId, HurtHp)
}

func (c CharacterBaseCard) Skill(TargetId int) {

}

func (c CharacterBaseCard) Death(AttackId int) {
	c.Notify(BattleData.AnDeath, -1)
	SelectCharacterCard := make([]int, 0)
	c.SetCardBt(&SelectCharacterCard)
	c.DisCard(&[]int{c.GetTempId()}) //反向压入
	c.Interrupt(&SelectCharacterCard, global.SelectCharacterTime*time.Second, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 1)
}
func (c CharacterBaseCard) BtCry() {

}
func (c CharacterBaseCard) RoundEnd() {

}

//todo
//---------二次分装---------

func (c CharacterBaseCard) EffectAttack(targetTempId int, AtkHp float64) {
	c.BtCtx.ProtoColPush(protocol.NewAttack(c.OwnerId, c.TempId, targetTempId, AtkHp))
}
func (c CharacterBaseCard) EffectHurt(AttackId int, AtkHp float64) {
	c.BtCtx.ProtoColPush(protocol.NewHurt(c.OwnerId, AttackId, c.TempId, AtkHp))
}
func (c CharacterBaseCard) Notify(Beh BattleData.AnimationBehavior, UserID int) {
	c.BtCtx.Notify(BattleData.NewAnimationDto(c.TempId, Beh, c.BtCtx.GetBtCardInfo(c.OwnerId)), UserID)
}
func (c CharacterBaseCard) Interrupt(res *[]int, time time.Duration, TempIdList []int, SelectNum int) { //res一定要塞到effect函数里处理
	resChan := make(chan []int)
	c.BtCtx.ProtoColPush(&protocol.Interrupt{
		UserId:     c.OwnerId,
		Time:       time,
		TempIdList: TempIdList,
		SelectNum:  SelectNum,
		Res:        resChan,
	})
	go func() {
		val := <-resChan
		*res = val
		c.BtCtx.ProtoColCancelInterrupt()
	}()
}

func (c CharacterBaseCard) DisCard(TempIdList *[]int) {
	c.BtCtx.ProtoColPush(protocol.NewDisCard(c.OwnerId, TempIdList))
}

func (c CharacterBaseCard) SetCardBt(TempIdList *[]int) {
	c.BtCtx.ProtoColPush(protocol.NewSetCardBt(c.OwnerId, TempIdList))
}

func (c CharacterBaseCard) GiveBuff(TempId int, b protocol.Buff) {
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(TempId, b, &c.BuffList))
}
