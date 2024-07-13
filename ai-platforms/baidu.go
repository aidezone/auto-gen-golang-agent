package platforms

import (
    "os"
    "context"

    // "aidezone/auto-gen-golang-agent/logger"
    "aidezone/auto-gen-golang-agent/defines"

    "github.com/baidubce/bce-qianfan-sdk/go/qianfan"
)

// TODO 改用环境变量或者配置文件
const (
	BAIDU_MODEL_NAME = "ERNIE-Speed-128K"
	BAIDU_ACCESS_KEY = "d8d0391ad22045f2906e0d7a32bbf5b0"
	BAIDU_SECRET_KEY = ""
)

type Baidu struct {
	chat *qianfan.ChatCompletion
	ctx *context.Context
}

func NewBaidu() *Baidu {
	os.Setenv("QIANFAN_ACCESS_KEY", BAIDU_ACCESS_KEY)
	os.Setenv("QIANFAN_SECRET_KEY", BAIDU_SECRET_KEY)
	_ctx := context.TODO()
	return &Baidu{
		// 调用对话Chat，可以通过 WithModel 指定模型，例如指定ERNIE-3.5-8K，参数对应ERNIE-Bot
		chat: qianfan.NewChatCompletion(
	        qianfan.WithModel(BAIDU_MODEL_NAME),
	    ),
	    ctx: &_ctx,
	}
}

func (s *Baidu) Call(ctx []*defines.Message) (*string, error) {
	req := s.generateRequest(ctx)

    // 发起对话
    resp, err := s.chat.Do(*s.ctx, req)
    if err != nil {
        return nil, err
    }
    return &resp.Result, nil
}

func (s *Baidu) generateRequest(ctx []*defines.Message) (*qianfan.ChatCompletionRequest) {
	req := &qianfan.ChatCompletionRequest{
		ResponseFormat: "json_object",
        // Messages: []qianfan.ChatCompletionMessage{
        //     qianfan.ChatCompletionUserMessage("帮我设计一个功能丰富的ktv点歌系统"),
        //     qianfan.ChatCompletionAssistantMessage("你的回答"),
        //     qianfan.ChatCompletionUserMessage("这个点歌系统里可以有哪些增值服务？"),
        // },
    }
	for _, msg := range ctx {
		if msg.Actor == defines.USER {
			req.Messages = append(req.Messages, qianfan.ChatCompletionUserMessage(*msg.Msg))
		}
		if msg.Actor == defines.ROBOT {
			req.Messages = append(req.Messages, qianfan.ChatCompletionAssistantMessage(*msg.Msg))
		}
	}
	return req
}