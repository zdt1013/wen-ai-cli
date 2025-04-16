package wenai

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// Generate 使用大语言模型生成回复
// ctx: 上下文
// llm: 大语言模型实例
// in: 输入的消息列表
// 返回: 生成的回复消息
func Generate(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

// Stream 使用大语言模型进行流式生成回复
// ctx: 上下文
// llm: 大语言模型实例
// in: 输入的消息列表
// 返回: 生成的流式回复读取器
func Stream(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {
	result, err := llm.Stream(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}
