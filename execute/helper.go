package execute

import (
	"strconv"
	"wen-ai-cli/common"
	"wen-ai-cli/setup"
)

func PrintHelp() *common.StreamPrinter {
	i18n := setup.GetI18n()
	// 打印帮助信息
	var helpPrinter = common.NewStreamPrinterWithAllOptions(false, true, i18n.ChatHelp, setup.CliVersion)
	helpPrinter.Print("1. q/quit->退出\n")
	helpPrinter.Print("2. f/finish->完成对话\n")
	helpPrinter.Print("3. 任意内容->继续对话\n")
	helpPrinter.Flush()
	return helpPrinter
}

// 打印对话轮次
func PrintQuestionTimes(question string, questionTimes int) {
	i18n := setup.GetI18n()
	// 创建使用自定义内容颜色的打印器
	var printer = common.NewStreamPrinterWithAllOptions(false, true, i18n.UserInput, setup.CliVersion)
	printer.Print("## 第 " + strconv.Itoa(questionTimes) + " 次对话\n")
	printer.Print(question)
	printer.Print("\n")
	printer.Flush()
}
