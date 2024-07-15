package platforms

import (
    "context"

	"github.com/hanyuancheung/gpt-go"
	"aidezone/auto-gen-golang-agent/logger"
	"aidezone/auto-gen-golang-agent/defines"
)


// TODO 改用环境变量或者配置文件
const (
	DEEPSEEK_MODEL_NAME = "deepseek-coder"
	DEEPSEEK_SECRET_KEY = "sk-634609855080435da8429061a7dc012a"
)

type Deepseek struct {
	ctx *context.Context
	client gpt.Client
}

func NewDeepseek() *Deepseek {
	client := gpt.NewClient(DEEPSEEK_SECRET_KEY, gpt.WithBaseURL("https://api.deepseek.com"))
	_ctx := context.Background()
	return &Deepseek{
	    ctx: &_ctx,
	    client: client,
	}
}

func (s *Deepseek) Call(ctx []*defines.Message) (*string, error) {
	req := s.generateRequest(ctx)
	resp, err := s.client.ChatCompletion(*s.ctx, req)
	if err != nil {
		logger.Infof("call deepseek api error: %v", err)
		return nil, err
	}
	if len(resp.Choices) < 1 {
		return nil, nil
	}
	return &resp.Choices[0].Message.Content, nil
}

func (s *Deepseek) generateRequest(ctx []*defines.Message) (*gpt.ChatCompletionRequest) {
	req := &gpt.ChatCompletionRequest{
		Model: DEEPSEEK_MODEL_NAME,
		MaxTokens:   3000,
		Temperature: 1.0,
	}
	for _, msg := range ctx {
		if msg.Actor == defines.USER {
			req.Messages = append(req.Messages, gpt.ChatCompletionRequestMessage{
				Role: "user",
				Content: *msg.Msg,
			})
		}
		if msg.Actor == defines.ROBOT {
			req.Messages = append(req.Messages, gpt.ChatCompletionRequestMessage{
				Role: "assistant",
				Content: *msg.Msg,
			})
		}
	}
	return req
}