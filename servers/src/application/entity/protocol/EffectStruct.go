package protocol

type Attack struct {
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
}

func NewAttack(UserId int, SendTempId int, TargetTempId int, AtkValue float64) *Attack {
	res := Attack{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	return &res
}

func (A *Attack) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColCardBtAttack(A.SendTempId, A.UserId, A.TargetTempId, A.AtkValue)
}

type Hurt struct {
	UserId       int
	SendTempId   int
	TargetTempId int
	AtkValue     float64
}

func (A *Hurt) Execute(pc ProtocolCardWithCtx) {
	pc.ProtoColReduceCardBtHp(A.SendTempId, A.UserId, A.TargetTempId, A.AtkValue)
}

func NewHurt(UserId int, SendTempId int, TargetTempId int, AtkValue float64) *Hurt {
	res := Hurt{}
	res.UserId = UserId
	res.SendTempId = SendTempId
	res.TargetTempId = TargetTempId
	res.AtkValue = AtkValue
	return &res
}
