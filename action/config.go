package action

import (
	"context"
	"fmt"
	"wen-ai-cli/setup"

	"github.com/urfave/cli/v3"
)

// NewConfigAction 创建 config action执行
func NewConfigAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		cfg := setup.GetConfig()
		// 获取命令行参数，如果存在则设置默认语言
		lang := cmd.String("lang")
		if lang != "" {
			cfg.DefaultLang = lang
		}
		// 设置openaiApiKey
		apiKey := cmd.String("apiKey")
		if apiKey != "" {
			cfg.OpenAI.APIKey = apiKey
		}
		// 设置openaiBaseURL
		baseURL := cmd.String("baseURL")
		if baseURL != "" {
			cfg.OpenAI.BaseURL = baseURL
		}
		// 设置openaiModel
		model := cmd.String("model")
		if model != "" {
			cfg.OpenAI.Model = model
		}
		// 保存配置
		setup.SaveConfig(cfg)
		fmt.Println("配置保存成功")
		return nil
	}
}
