package wenai

import (
	"io"
	"log"
	"regexp"
	"strings"
	"wen-ai-cli/common"
	"wen-ai-cli/model"
	"wen-ai-cli/setup"

	"github.com/cloudwego/eino/schema"
)

func ReportStream(sr *schema.StreamReader[*schema.Message]) (*schema.Message, *model.HiddenParams, error) {
	defer sr.Close()

	// 创建使用自定义内容颜色的打印器
	printer := common.NewStreamPrinterWithAllOptions(false, true, setup.CliName, setup.CliVersion)

	i := 0
	result := &model.HiddenParams{}
	shellCode := ""
	fullContentBuilder := strings.Builder{}
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			// 处理最后一段
			printer.Print("\n")
			printer.Flush()

			// 移动正则解析代码到这里
			re := regexp.MustCompile("(?s)```code(.*?)```")
			fullContent := fullContentBuilder.String()
			// 查找如果有多个代码块匹配，则认为最后一个代码块是shellCode
			matches := re.FindAllStringSubmatch(fullContent, -1)
			if len(matches) > 0 {
				// 获取最后一个匹配的代码块
				lastMatch := matches[len(matches)-1]
				if len(lastMatch) > 1 {
					shellCode = strings.TrimSpace(lastMatch[1])
				}
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
		printer.Print(content)
		i++
	}
}
