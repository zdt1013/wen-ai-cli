package action

import (
	"context"
	"strings"
	"wen-ai-cli/logger"
	"wen-ai-cli/setup"
	"wen-ai-cli/wenai"
	"wen-ai-cli/wenai/manual"

	"github.com/urfave/cli/v3"
)

// NewWenManualAction 创建 manual action执行
func NewWenManualAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		i18n := setup.GetI18n()
		cmdName := cmd.String("cmd")
		question := strings.Join(cmd.Args().Slice(), " ")
		answerConfig := setup.GetConfig().AnswerConfig
		messages := manual.CreateOnceMessagesFromTemplate(cmdName, question, answerConfig.EnableExplain, answerConfig.EnableExtendParams, answerConfig.EnablePlatformPerception, answerConfig.EnableWorkUserAndDir)
		cm := wenai.CreateOpenAIChatModel(ctx)
		streamResult := wenai.Stream(ctx, cm, messages)
		_, _, err := wenai.ReportStream(streamResult)
		if err != nil {
			logger.Errorf("ReportStream failed %v", err)
		}
		logger.Debug(i18n.Exit)
		return nil
	}
}
