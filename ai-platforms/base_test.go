package platforms

import (
	"os"
	"testing"
	"aidezone/auto-gen-golang-agent/logger"
)

// 测试前的初始化
func TestMain(m *testing.M) {
	// 初始化日志
	logger.InitLogger("", true, false)

	// 调用 testing.M 的 Run 方法执行所有测试
	exitCode := m.Run()

	// 退出测试程序，返回测试结果
	os.Exit(exitCode)
}