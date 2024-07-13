package platforms

import (
    "context"
	
	"aidezone/auto-gen-golang-agent/defines"
)

type Ollama struct {
	ctx *context.Context
}

func NewOllama() *Ollama {
	_ctx := context.TODO()
	return &Ollama{
	    ctx: &_ctx,
	}
}

func (s *Ollama) Call(ctx []*defines.Message) (*string, error) {
	return nil, nil
}