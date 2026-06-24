package CardAbstract

import (
	"context"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
)

type Card interface {
	GetID() int
	SetInfo(info map[string]any)
	GetInfo() map[string]any
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
	GetBuffList() *[]*protocol.Buff
	AppendBuff(b *protocol.Buff)
	GetDec() *CardMeta.Decorator
	AddBuff(buff *protocol.Buff, pc protocol.ProtocolCardWithCtx)
	BuffRoundEnd(pc protocol.ProtocolCardWithCtx)
	GetForm() BattleData.Form
	SetForm(BattleData.Form)
	ShareInit(goctx context.Context, ctx protocol.ProtocolCardWithCtx)
	PutBroadInfo(v *CardMeta.BroadInfo) //把广播信息丢到管道里
	BroadCallBack(v *CardMeta.BroadInfo)
	NextRound()
}

func GetCardDto(c Card) BattleData.CardDto {
	res := BattleData.CardDto{}
	res.Id = c.GetID()
	res.TempId = c.GetTempId()
	res.Hp = c.GetHpNow()
	res.Damage = c.GetAtkNow()
	BuffDtoList := make([]BattleData.BuffDto, 0, 8)
	for _, buff := range *c.GetBuffList() {
		BuffDtoList = append(BuffDtoList, buff.GetBuffDto())
	}
	res.BuffDtoList = BuffDtoList
	res.Form = c.GetForm()
	return res
}
