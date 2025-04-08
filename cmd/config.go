package cmd

import (
	"wen-ai-cli/action"
	"wen-ai-cli/setup"

	"github.com/gookit/i18n"
	"github.com/urfave/cli/v3"
)

// NewConfigCmd 创建 config 命令
func NewConfigCmd() *cli.Command {
	return &cli.Command{
		Name:    setup.ConfigCmd,
		Aliases: []string{setup.ConfigCmdAlias},
		Usage:   i18n.Dtr("configMode"),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "lang",
				Aliases: []string{"l"},
				Value:   "zh-CN",
				Usage:   i18n.Dtr("configLang"),
			},
			&cli.StringFlag{
				Name:    "apiKey",
				Aliases: []string{"k"},
				Value:   "",
				Usage:   i18n.Dtr("configAk"),
			},
			&cli.StringFlag{
				Name:    "baseURL",
				Aliases: []string{"u"},
				Value:   "",
				Usage:   i18n.Dtr("configBaseURL"),
			},
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Value:   "",
				Usage:   i18n.Dtr("configModel"),
			},
		},
		Action: action.NewConfigAction(),
	}
}
