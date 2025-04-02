package wenai

import (
	"io"
	"log"
	"regexp"
	"strings"
	"wen-ai-cli/common"
	"wen-ai-cli/model"

	"github.com/cloudwego/eino/schema"

	"github.com/fatih/color"
)

func ReportStream(sr *schema.StreamReader[*schema.Message]) (*schema.Message, *model.HiddenParams, error) {
	defer sr.Close()

	// 创建使用自定义内容颜色的打印器
	printer := common.NewStreamPrinterWithColors(
		"╭──",                        // 顶部边框
		"│",                          // 普通行边框
		"│──",                        // 标题行边框
		"╰──",                        // 底部边框
		"WenAI",                      // 头部文本
		"End",                        // 底部文本
		color.New(color.FgHiCyan),    // 顶部边框颜色
		color.New(color.FgHiCyan),    // 普通行边框颜色
		color.New(color.FgHiCyan),    // 标题行边框颜色
		color.New(color.FgHiCyan),    // 底部边框颜色
		color.New(color.FgHiMagenta), // 头部文本颜色 - 亮紫色
		color.New(color.FgHiMagenta), // 底部文本颜色 - 亮紫色
		color.New(color.FgHiGreen),   // 标题文本颜色 - 亮绿色
		color.New(color.FgHiYellow),  // 有序列表颜色 - 亮紫色
		color.New(color.FgHiYellow),  // 无序列表颜色 - 亮黄色
		color.New(color.FgHiBlue),    // 代码块颜色 - 亮蓝色
	)

	i := 0
	result := &model.HiddenParams{}
	shellCode := ""
	fullContentBuilder := strings.Builder{}
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			// 处理最后一段
			printer.ProcessFragment("\n")
			printer.Flush()

			// 移动正则解析代码到这里
			re := regexp.MustCompile("(?s)```code(.*?)```")
			fullContent := fullContentBuilder.String()
			matches := re.FindAllStringSubmatch(fullContent, -1)
			for _, match := range matches {
				shellCode = strings.TrimSpace(match[1])
			}
			if shellCode != "" {
				result.ShellCode = shellCode
				// 解析shellCode,<下载文件的URL,url>序列化成hideParams
				re := regexp.MustCompile(`<([\p{Han}a-zA-Z0-9]+),(\w+)>`)
				matches := re.FindAllStringSubmatch(shellCode, -1)
				for _, match := range matches {
					paramName := match[1]
					paramType := match[2]
					result.NeedFillParams = append(result.NeedFillParams, model.ParamInfo{
						Param: paramName,
						Type:  paramType,
					})
				}
			}
			fullMessage := &schema.Message{
				Role:    "assistant",
				Content: fullContent,
			}
			return fullMessage, result, nil
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		content := message.Content
		fullContentBuilder.WriteString(content)
		printer.ProcessFragment(content)
		i++
	}
}
