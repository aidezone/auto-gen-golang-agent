package platforms

import (
	"aidezone/auto-gen-golang-agent/defines"
)


type AIPlatform interface {
    Call(ctx []*defines.Message) (*string, error)
}


func NewAIPlatform(platform defines.PlatformName) AIPlatform {
	switch platform {
	case defines.Baidu:
		return NewBaidu()
	case defines.Ollama:
		return NewOllama()
	case defines.Openai:
		return NewOpenai()
	default:
		return nil
	}
}