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

	"github.com/gookit/i18n"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

func NewWenOnceAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
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
			prompt := promptui.Select{
				HideHelp: true,
				Label:    i18n.Dtr("selectOperation"),
				Items:    []string{i18n.Dtr("fillParamsAndRun"), i18n.Dtr("adjustAndRun"), i18n.Dtr("exit")},
			}
			_, result, err := prompt.Run()
			if err != nil {
				logger.Errorf("Prompt failed %v", err)
				return nil
			}
			logger.Debugf(i18n.Dtr("yourChoice"), result)
			if result == i18n.Dtr("fillParamsAndRun") {
				shell_code, shouldExecute := common.HandleParamsCompletion(hidden_params)
				if shouldExecute {
					execute.ExecuteScript(shell_code)
				} else {
					return nil
				}
			} else if result == i18n.Dtr("adjustAndRun") {
				script, shouldExecute := common.HandleScriptAdjustment(hidden_params.ShellCode)
				if shouldExecute {
					execute.ExecuteScript(script)
				} else {
					return nil
				}

			} else {
				logger.Debug(i18n.Dtr("exit"))
			}
		} else {
			if hidden_params.ShellCode == "" {
				// 如果脚本为空，则提示用户，说明无法解析答案
				// 按照微调脚本进行处理
				prompt := promptui.Select{
					Label: i18n.Dtr("selectOperation"),
					Items: []string{i18n.Dtr("exit")},
				}
				_, result, err := prompt.Run()

				if err != nil {
					logger.Errorf("Prompt failed %v", err)
					return nil
				}
				logger.Debugf(i18n.Dtr("yourChoice"), result)
				if result == i18n.Dtr("adjustAndRun") {
					script, shouldExecute := common.HandleScriptAdjustment(hidden_params.ShellCode)
					if shouldExecute {
						execute.ExecuteScript(script)
					}
				} else {
					logger.Debug(i18n.Dtr("exit"))
				}
			} else {
				// 如果脚本不为空，则提示用户，说明可以执行
				logger.Debug(i18n.Dtr("canExecute"))
				prompt := promptui.Select{
					Label: i18n.Dtr("selectOperation"),
					Items: []string{i18n.Dtr("runNow"), i18n.Dtr("adjustAndRun"), i18n.Dtr("exit")},
				}
				_, result, err := prompt.Run()

				if err != nil {
					logger.Errorf("Prompt failed %v", err)
					return nil
				}

				logger.Debugf(i18n.Dtr("yourChoice"), result)
				if result == i18n.Dtr("runNow") {
					execute.ExecuteScript(hidden_params.ShellCode)
				} else if result == i18n.Dtr("adjustAndRun") {
					script, shouldExecute := common.HandleScriptAdjustment(hidden_params.ShellCode)
					if shouldExecute {
						execute.ExecuteScript(script)
					}
				} else {
					logger.Debug(i18n.Dtr("exit"))
				}
			}
		}
		return nil
	}
}
