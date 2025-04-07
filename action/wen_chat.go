package action

import (
	"context"
	"strings"
	"wen-ai-cli/common"
	"wen-ai-cli/execute"
	"wen-ai-cli/logger"
	"wen-ai-cli/setup"

	"wen-ai-cli/wenai"

	"github.com/cloudwego/eino/schema"
	"github.com/urfave/cli/v3"
)

func NewWenChatAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// 获取配置信息
		answerConfig := setup.GetConfig().AnswerConfig
		// 获取语言包
		i18n := setup.GetI18n()
		// 初始化聊天历史记录
		chatHistory := []*schema.Message{}
		// 将命令行参数拼接为问题
		question := strings.Join(cmd.Args().Slice(), " ")
		questionTimes := 0
		if question == "" {
			firstQuestion, err := execute.InputString(i18n.UserInput)
			if err != nil {
				logger.Errorf("Prompt failed %v", err)
				return nil
			}

			// 记录用户输入
			logger.Debugf(i18n.UserInput, firstQuestion)
			question = firstQuestion
		}

		// 进入主循环，持续与用户交互
		for {
			questionTimes++
			// 打印对话轮次
			execute.PrintQuestionTimes(question, questionTimes)
			// 创建聊天消息模板

			messages := wenai.CreateMoreMessagesFromTemplate(question, chatHistory, answerConfig.EnableExplain, answerConfig.EnableExtendParams, answerConfig.EnablePlatformPerception)
			// 创建OpenAI聊天模型
			cm := wenai.CreateOpenAIChatModel(ctx)
			// 获取流式处理结果
			streamResult := wenai.Stream(ctx, cm, messages)
			// 解析流式结果，获取完整消息和隐藏参数
			fullMessage, hiddenParams, err := wenai.ReportStream(streamResult)
			if err != nil {
				logger.Errorf("ReportStream failed %v", err)
			}

			// 打印帮助信息
			var helpPrinter = execute.PrintHelp()
			inputQuetion, err := execute.InputString(i18n.UserInput)
			helpPrinter.Clear0()
			if err != nil {
				logger.Errorf("Prompt failed %v", err)
				return nil
			}

			// 记录用户输入
			logger.Debugf(i18n.UserInputFormat, inputQuetion)

			// 处理退出命令
			if inputQuetion == "q" || inputQuetion == "Q" {
				logger.Debug(i18n.Exit)
				return nil
			}

			// 处理功能命令
			if inputQuetion == "f" || inputQuetion == "F" {
				// 如果有需要填充的参数
				if hiddenParams.HasParameters() {
					// 创建操作选择提示
					result, err := execute.Prompt(i18n.SelectOperation, []string{i18n.FillParamsAndRun, i18n.AdjustAndRun, i18n.Exit})
					if err != nil {
						logger.Errorf("Prompt failed %v", err)
						return nil
					}

					// 记录用户选择
					logger.Debugf(i18n.YourChoice, result)

					// 根据选择执行相应操作
					if result == i18n.FillParamsAndRun {
						shellCode, shouldExecute := common.HandleParamsCompletion(hiddenParams)
						if shouldExecute {
							execute.ExecuteScript(shellCode)
						}
					} else if result == i18n.AdjustAndRun {
						script, shouldExecute := common.HandleScriptAdjustment(hiddenParams.ShellCode)
						if shouldExecute {
							execute.ExecuteScript(script)
						}
					} else {
						logger.Debug(i18n.Exit)
					}
				} else {
					if hiddenParams.ShellCode == "" {
						// 如果脚本为空，则提示用户，说明无法解析答案
						// 按照微调脚本进行处理
						result, err := execute.Prompt(i18n.SelectOperation, []string{i18n.Exit})
						if err != nil {
							logger.Errorf("Prompt failed %v", err)
							return nil
						}
						logger.Debugf(i18n.YourChoice, result)
						if result == i18n.AdjustAndRun {
							script, shouldExecute := common.HandleScriptAdjustment(hiddenParams.ShellCode)
							if shouldExecute {
								execute.ExecuteScript(script)
							}
						} else {
							logger.Debug(i18n.Exit)
						}
					} else {
						// 如果脚本不为空，则提示用户，说明可以执行
						logger.Debug(i18n.CanExecute)
						// 创建操作选择提示
						result, err := execute.Prompt(i18n.SelectOperation, []string{i18n.RunNow, i18n.AdjustAndRun, i18n.Exit})
						if err != nil {
							logger.Errorf("Prompt failed %v", err)
							return nil
						}

						// 记录用户选择
						logger.Debugf(i18n.YourChoice, result)

						// 根据选择执行相应操作
						if result == i18n.RunNow {
							execute.ExecuteScript(hiddenParams.ShellCode)
						} else if result == i18n.AdjustAndRun {
							script, shouldExecute := common.HandleScriptAdjustment(hiddenParams.ShellCode)
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

			// 其他情况，继续对话，并更新聊天历史记录
			// 保留最近10条消息
			chatHistory = messages[max(1, len(messages)-10):]
			// 添加最新消息到历史记录
			chatHistory = append(chatHistory, fullMessage)
			// 更新问题为最新输入
			question = inputQuetion
		}
	}
}
