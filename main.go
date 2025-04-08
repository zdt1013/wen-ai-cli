package main

import (
	"context"
	"os"
	"wen-ai-cli/action"
	"wen-ai-cli/cmd"
	"wen-ai-cli/logger"
	"wen-ai-cli/setup"

	"github.com/gookit/i18n"
	"github.com/urfave/cli/v3"
)

func main() {
	// 初始化配置
	setup.InitConfig()
	// 初始化多语言
	setup.InitLang()
	// 初始化命令
	app := &cli.Command{
		Name:   "wen",
		Usage:  i18n.Dtr("usage"),
		Action: action.NewWenOnceAction(),
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			// 获取当前要运行的command
			command := cmd.Args().First()
			// 如果command是config，则不检查必要配置
			if command == setup.ConfigCmd || command == setup.ConfigCmdAlias {
				return ctx, nil
			}
			// 检查必要配置
			config := setup.GetConfig()
			if config.OpenAI.APIKey == "" || config.OpenAI.BaseURL == "" || config.OpenAI.Model == "" {
				return nil, cli.Exit(i18n.Dtr("configError"), 400)
			}
			return ctx, nil
		},
		Commands: []*cli.Command{
			cmd.NewChatCmd(),
			cmd.NewConfigCmd(),
			cmd.NewManualCmd(),
		},
	}
	// 运行命令
	if err := app.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err.Error())
	}
}
