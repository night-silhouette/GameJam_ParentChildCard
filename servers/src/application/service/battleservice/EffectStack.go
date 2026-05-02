package battleservice

import (
	"pcc_card/application/entity/protocol"
)

type EffectStack []protocol.Effect

func NewEffectStack() *EffectStack {
	return &EffectStack{}
}

func (s *EffectStack) Push(e protocol.Effect) {
	*s = append(*s, e)
}

func (s *EffectStack) Pop() protocol.Effect {
	if len(*s) == 0 {
		return nil
	}
	stack := *s
	top := stack[len(stack)-1]
	*s = stack[:len(stack)-1]
	return top
}

func (s *EffectStack) IsEmpty() bool {
	return len(*s) == 0
}
