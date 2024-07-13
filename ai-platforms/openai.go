package platforms

import (
    "context"

	"aidezone/auto-gen-golang-agent/defines"
)

type Openai struct {
	ctx *context.Context
}

func NewOpenai() *Openai {
	_ctx := context.TODO()
	return &Openai{
	    ctx: &_ctx,
	}
}

func (s *Openai) Call(ctx []*defines.Message) (*string, error) {
	return nil, nil
}