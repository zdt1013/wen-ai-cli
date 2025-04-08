package cmd

import (
	"wen-ai-cli/action"
	"wen-ai-cli/setup"

	"github.com/gookit/i18n"
	"github.com/urfave/cli/v3"
)

// NewManualCmd 创建 manual 命令
func NewManualCmd() *cli.Command {
	return &cli.Command{
		Name:  setup.ManualCmd,
		Usage: i18n.Dtr("manualCmdUsage"),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "cmd",
				Aliases: []string{"c"},
				Value:   "",
				Usage:   i18n.Dtr("manualCmdFlag"),
			},
		},
		Action: action.NewWenManualAction(),
	}
}
