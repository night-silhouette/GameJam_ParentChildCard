package CardImpl

import "pcc_card/application/entity/protocol"

type CharacterBaseCard struct {
	BaseCard
}

func (c CharacterBaseCard) Attack(tempId int) {
	c.BtCtx.ProtoColPush(protocol.NewAttack(c.OwnerId, c.TempId, tempId, c.AtkNow))
}

func (c CharacterBaseCard) Hurt(tempId int, HurtHp float64) {

}

func (c CharacterBaseCard) Skill(tempId int) {

}

func (c CharacterBaseCard) Death(tempId int) {

}
