package ssabuilder

import "golang.org/x/tools/go/ssa"

type Send struct {
	*send
}

type send struct {
	*ssa.Send
	calls []ssa.CallInstruction
}

func (c *Send) String() string {
	return c.Send.String()
}

func (c *Send) SetSend(s *ssa.Send) {
	c.Send = s
}

func (c *Send) GetSend() *ssa.Send {
	return c.Send
}

func (c *Send) AddCall(ci ssa.CallInstruction) {
	if c.calls == nil {
		c.calls = make([]ssa.CallInstruction, 1)
	}
	c.calls = append(c.calls, ci)
}

func (c *Send) GetCalls() []ssa.CallInstruction {
	return c.calls
}
