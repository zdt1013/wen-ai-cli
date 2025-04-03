package action

import (
	"context"
	"strconv"
	"strings"
	"wen-ai-cli/common"
	"wen-ai-cli/execute"
	"wen-ai-cli/logger"
	"wen-ai-cli/setup"
	"wen-ai-cli/validate"

	"wen-ai-cli/wenai"

	"github.com/cloudwego/eino/schema"
	"github.com/gookit/i18n"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v3"
)

func NewWenChatAction() cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		// 获取配置信息
		answerConfig := setup.GetConfig().AnswerConfig
		// 初始化聊天历史记录
		chat_history := []*schema.Message{}
		// 将命令行参数拼接为问题
		question := strings.Join(cmd.Args().Slice(), " ")
		question_times := 0
		if question == "" {
			// 定义输入验证函数
			validateFn := func(input string) error {
				return validate.ValidateParam(input, "string")
			}

			// 创建用户输入提示
			prompt := promptui.Prompt{
				Label:       i18n.Dtr("userInput"),
				Validate:    validateFn,
				HideEntered: true,
			}
			// 获取用户输入
			first_question, err := prompt.Run()

			if err != nil {
				logger.Errorf("Prompt failed %v", err)
				return nil
			}

			// 记录用户输入
			logger.Debugf(i18n.Dtr("userInput"), first_question)
			question = first_question
		}

		// 进入主循环，持续与用户交互
		for {
			question_times++
			// 创建使用自定义内容颜色的打印器
			var printer = common.NewStreamPrinterWithAllOptions(false, true, i18n.Dtr("userInput"), setup.CliVersion)
			printer.Print("## 第 " + strconv.Itoa(question_times) + " 次对话\n")
			printer.Print(question)
			printer.Print("\n")
			printer.Flush()
			// 创建聊天消息模板

			messages := wenai.CreateMoreMessagesFromTemplate(question, chat_history, answerConfig.EnableExplain, answerConfig.EnableExtendParams, answerConfig.EnablePlatformPerception)
			// 创建OpenAI聊天模型
			cm := wenai.CreateOpenAIChatModel(ctx)
			// 获取流式处理结果
			streamResult := wenai.Stream(ctx, cm, messages)
			// 解析流式结果，获取完整消息和隐藏参数
			fullMessage, hidden_params, err := wenai.ReportStream(streamResult)
			if err != nil {
				logger.Errorf("ReportStream failed %v", err)
			}

			// 打印帮助信息
			var helpPrinter = common.NewStreamPrinterWithAllOptions(false, true, i18n.Dtr("chatHelp"), setup.CliVersion)
			helpPrinter.Print("1. q/quit->退出\n")
			helpPrinter.Print("2. f/finish->完成对话\n")
			helpPrinter.Print("3. 任意内容->继续对话\n")
			helpPrinter.Flush()

			// 定义输入验证函数
			validateFn := func(input string) error {
				return validate.ValidateParam(input, "string")
			}

			// 创建用户输入提示
			prompt := promptui.Prompt{
				Label:       i18n.Dtr("userInput"),
				Validate:    validateFn,
				HideEntered: true,
			}
			// 获取用户输入
			input_quetion, err := prompt.Run()
			helpPrinter.Clear0()

			if err != nil {
				logger.Errorf("Prompt failed %v", err)
				return nil
			}

			// 记录用户输入
			logger.Debugf(i18n.Dtr("userInputFormat"), input_quetion)

			// 处理退出命令
			if input_quetion == "q" || input_quetion == "Q" {
				logger.Debug(i18n.Dtr("exit"))
				return nil
			}

			// 处理功能命令
			if input_quetion == "f" || input_quetion == "F" {
				// 如果有需要填充的参数
				if hidden_params.HasParameters() {
					// 创建操作选择提示
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

					// 记录用户选择
					logger.Debugf(i18n.Dtr("yourChoice"), result)

					// 根据选择执行相应操作
					if result == i18n.Dtr("fillParamsAndRun") {
						shell_code, shouldExecute := common.HandleParamsCompletion(hidden_params)
						if shouldExecute {
							execute.ExecuteScript(shell_code)
						}
					} else if result == i18n.Dtr("adjustAndRun") {
						script, shouldExecute := common.HandleScriptAdjustment(hidden_params.ShellCode)
						if shouldExecute {
							execute.ExecuteScript(script)
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
						// 创建操作选择提示
						prompt := promptui.Select{
							Label: i18n.Dtr("selectOperation"),
							Items: []string{i18n.Dtr("runNow"), i18n.Dtr("adjustAndRun"), i18n.Dtr("exit")},
						}
						_, result, err := prompt.Run()

						if err != nil {
							logger.Errorf("Prompt failed %v", err)
							return nil
						}

						// 记录用户选择
						logger.Debugf(i18n.Dtr("yourChoice"), result)

						// 根据选择执行相应操作
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

			// 其他情况，继续对话，并更新聊天历史记录
			// 保留最近10条消息
			chat_history = messages[max(1, len(messages)-10):]
			// 添加最新消息到历史记录
			chat_history = append(chat_history, fullMessage)
			// 更新问题为最新输入
			question = input_quetion
		}
	}
}
