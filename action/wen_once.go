package action

import (
	"context"
	"fmt"
	"strings"
	"wen-ai-cli/common"
	"wen-ai-cli/execute"
	"wen-ai-cli/logger"
	"wen-ai-cli/setup"
	"wen-ai-cli/wenai"

	"github.com/urfave/cli/v3"
)

func NewWenOnceAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		i18n := setup.GetI18n()
		question := strings.Join(cmd.Args().Slice(), " ")
		answerConfig := setup.GetConfig().AnswerConfig
		messages := wenai.CreateOnceMessagesFromTemplate(question, answerConfig.EnableExplain, answerConfig.EnableExtendParams, answerConfig.EnablePlatformPerception)
		cm := wenai.CreateOpenAIChatModel(ctx)
		streamResult := wenai.Stream(ctx, cm, messages)
		_, hidden_params, err := wenai.ReportStream(streamResult)
		if err != nil {
			logger.Errorf("ReportStream failed %v", err)
		}
		fmt.Println("--------------------------------")
		if hidden_params.HasParameters() {
			// 如果存在需要填充的参数，则提示用户，说明可以填充参数
			result, err := execute.Prompt(i18n.SelectOperation, []string{i18n.FillParamsAndRun, i18n.AdjustAndRun, i18n.Exit})
			if err != nil {
				logger.Errorf("Prompt failed %v", err)
				return nil
			}
			logger.Debugf(i18n.YourChoice, result)
			if result == i18n.FillParamsAndRun {
				shell_code, shouldExecute := common.HandleParamsCompletion(hidden_params)
				if shouldExecute {
					execute.ExecuteScript(shell_code)
				} else {
					return nil
				}
			} else if result == i18n.AdjustAndRun {
				script, shouldExecute := common.HandleScriptAdjustment(hidden_params.ShellCode)
				if shouldExecute {
					execute.ExecuteScript(script)
				} else {
					return nil
				}

			} else {
				logger.Debug(i18n.Exit)
			}
		} else {
			if hidden_params.ShellCode == "" {
				// 如果脚本为空，则提示用户，说明无法解析答案
				// 按照微调脚本进行处理
				result, err := execute.Prompt(i18n.SelectOperation, []string{i18n.Exit})
				if err != nil {
					logger.Errorf("Prompt failed %v", err)
					return nil
				}
				logger.Debugf(i18n.YourChoice, result)
				logger.Debug(i18n.Exit)
			} else {
				// 如果脚本不为空，则提示用户，说明可以执行
				logger.Debug(i18n.CanExecute)
				result, err := execute.Prompt(i18n.SelectOperation, []string{i18n.RunNow, i18n.AdjustAndRun, i18n.Exit})
				if err != nil {
					logger.Errorf("Prompt failed %v", err)
					return nil
				}
				logger.Debugf(i18n.YourChoice, result)
				if result == i18n.RunNow {
					execute.ExecuteScript(hidden_params.ShellCode)
				} else if result == i18n.AdjustAndRun {
					script, shouldExecute := common.HandleScriptAdjustment(hidden_params.ShellCode)
					if shouldExecute {
						execute.ExecuteScript(script)
					}
				} else {
					logger.Debug(i18n.Exit)
				}
			}
		}
		return nil
	}
}
