package wenai

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/gookit/config/v2"
)

func CreateOpenAIChatModel(ctx context.Context) model.BaseChatModel {
	apiKey := config.String("openai.apiKey")
	baseURL := config.String("openai.baseURL")
	modelName := config.String("openai.model")
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  apiKey,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
