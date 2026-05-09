package CardAbstract

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/protocol"
)

type Card interface {
	GetID() int
	SetInfo(info map[string]any)
	GetInfo() map[string]any
	Clone() Card
	GetStateCodeChan() chan protocol.Effect
	SetStateCodeChan(chan protocol.Effect)
	GetTempId() int
	SetTempId(id int)
	GetOwnerId() int
	SetOwnerId(id int)
	GetHpNow() float64
	SetHpNow(hpNow float64)
	GetAtkNow() float64
	SetAtkNow(atkNow float64)
	SetBtCtx(btCtx protocol.ProtocolCardWithCtx)
	ReInitialize()
}

func GetCardDto(c Card) BattleData.CardDto {
	res := BattleData.CardDto{}
	res.Id = c.GetID()
	res.TempId = c.GetTempId()
	res.Hp = c.GetHpNow()
	res.Damage = c.GetAtkNow()
	return res
}
