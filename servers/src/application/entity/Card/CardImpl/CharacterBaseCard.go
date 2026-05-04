package CardImpl

import (
	"fmt"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/protocol"
)

type CharacterBaseCard struct {
	BaseCard
}

func (c CharacterBaseCard) Attack(TargetId int) {
	c.Notify(BattleData.AnAttack)
	c.EffectAttack(TargetId, c.AtkNow)
}

func (c CharacterBaseCard) Hurt(AttackId int, HurtHp float64) {
	c.Notify(BattleData.AnHurt)
	c.EffectHurt(AttackId, HurtHp)
}

func (c CharacterBaseCard) Skill(TargetId int) {

}

func (c CharacterBaseCard) Death(AttackId int) {
	fmt.Println("我死了")
}

func (c CharacterBaseCard) EffectAttack(targetTempId int, AtkHp float64) {
	c.StateCodeChan <- protocol.NewAttack(c.OwnerId, c.TempId, targetTempId, AtkHp)
}
func (c CharacterBaseCard) EffectHurt(AttackId int, AtkHp float64) {
	c.StateCodeChan <- protocol.NewHurt(c.OwnerId, AttackId, c.TempId, AtkHp)
}
func (c CharacterBaseCard) Notify(Beh BattleData.AnimationBehavior) {
	c.BtCtx.Notify(BattleData.MewAnimationDto(c.ID, c.TempId, Beh))
}
