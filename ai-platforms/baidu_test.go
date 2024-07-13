package platforms

import (
	"testing"

	"aidezone/auto-gen-golang-agent/defines"
	"aidezone/auto-gen-golang-agent/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewBaidu(t *testing.T) {
	// 测试 NewBaidu 是否正确初始化 Baidu 实例
	baiduInstance := NewBaidu()
	assert.NotNil(t, baiduInstance, "Baidu instance should not be nil")
}

func TestCall(t *testing.T) {

	// 创建 Baidu 实例
	baiduInstance := NewBaidu()

	// 模拟生成请求
	chat := defines.NewChat()
	req := chat.AppendAsk("你好")

	resp, err := baiduInstance.Call(req)

	logger.Infof("TestNewBaidu.TestCall resp: %v, err: %v", *resp, err)
}
