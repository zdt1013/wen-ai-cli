package common

import (
	"fmt"
	"testing"
)

// 修改main函数，添加简化版示例并更新用户指南
func TestStreamPrinter(t *testing.T) {
	// First demonstrate with predictable character-by-character streaming
	SimulateStreamInput()

	// Show an example of how to integrate with a real API
	ExampleWithRealAPI()

	// 演示自定义边框
	SimulateWithCustomBorder()

	// 演示自定义边框颜色
	SimulateWithCustomColors()

	// 演示自定义内容颜色
	SimulateWithCustomContentColors()

	// 演示简化版颜色配置
	SimulateWithSimpleColors()

	// Final instructions for users
	fmt.Println("\n\n=== 如何在您自己的API中使用此库 ===")
	fmt.Println("1. 创建流通道: streamChan := make(chan string)")
	fmt.Println("2. 启动处理器: go ProcessStream(streamChan)")
	fmt.Println("3. 发送API片段: streamChan <- fragment")
	fmt.Println("4. 完成后关闭通道: close(streamChan)")
	fmt.Println("5. 自定义边框: 使用NewStreamPrinterWithBorder创建自定义边框样式")
	fmt.Println("6. 自定义所有颜色: 使用NewStreamPrinterWithColors创建完全自定义样式")
	fmt.Println("7. 简化的颜色配置: 使用NewColoredStreamPrinter快速设置常用颜色")
	fmt.Println("8. 文本颜色可配置: 可单独设置头部文本和底部文本的颜色")
}
