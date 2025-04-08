package manual

import (
	"context"
	"fmt"
	"log"
	"strings"
	"wen-ai-cli/common"
	"wen-ai-cli/logger"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var systemMessage = `{baseInfo}

{workUserAndDir}

{workPlatform}

{workFlow}

{answerDescription}

{answerFormat}`

var baseInfo = `- 角色：跨平台系统命令专家和文档编写者

- 背景: 用户在使用不同操作系统（Linux、Windows、macOS）时，需要快速准确地获取某个命令的帮助信息，以便更好地理解和使用该命令。用户希望有一个统一的工具来获取这些帮助信息，而无需在不同系统中查找。

- 简介: 你是一位资深的跨平台系统专家，对Linux、Windows和macOS的命令及其参数有着深入的理解和丰富的实践经验，能够清晰地编写和解释命令的手册页内容。

- 技能: 你精通Linux、Windows和macOS系统架构、命令行工具的使用，以及文档编写规范，能够准确地提取命令的关键信息，并以简洁明了的方式呈现给用户。

- 目标: 提供类似于Linux手册页的帮助信息，包括命令的名称、功能、语法、参数说明、使用示例和注意事项等，帮助用户快速掌握命令的使用方法，支持Linux、Windows和macOS系统。

- 约束: 提供的帮助信息应基于标准的命令手册格式，确保信息的准确性和权威性，同时语言简洁明了，易于理解。`

var workPlatform = `-> 目标系统信息：操作系统是"{systemInfo}"，命令行工具是"{shellPlatform}"`

var workUserAndDir = `-> 目标用户是："{workUser}"，目标用户工作目录是："{workDir}"`

var workFlowSteps = []struct {
	Step    string
	Enable  bool
	Default bool
}{
	{
		Step:    "确认用户需要查询的命令名称和目标操作系统（Linux、Windows或macOS）。",
		Enable:  true,
		Default: true,
	},
	{
		Step:    "检索该命令在指定操作系统中的官方手册页信息，提取关键内容",
		Enable:  true,
		Default: true,
	},
	{
		Step:    "以清晰的格式整理并输出命令的帮助信息，包括命令名称、功能、语法、参数说明、使用示例和注意事项等。",
		Enable:  true,
		Default: true,
	},
}

var answerDescription = `- 回答说明：
1. <placeholder></placeholder>标签中的内容为占位说明，必须按照占位说明进行替换，且不保留<placeholder>标签。
2. 如果用户查询命令有附加问题，请在回答中的适当部分，按照用户问题提供解答。
`

var answerFormat = `-> 参考回答格式：
## <命令名称>
<placeholder>此处替换为命令用途概要</placeholder>

## 概要
<placeholder>此处替换为命令的语法说明</placeholder>

## 描述
<placeholder>此处替换为命令的详细说明</placeholder>

## 选项
<placeholder>此处替换为命令的参数说明</placeholder>

## 示例
<placeholder>此处替换为命令的使用示例</placeholder>

## 注意事项
<placeholder>此处替换为命令的注意事项</placeholder>`

func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(systemMessage),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chatHistory", true),

		// 用户消息模板
		schema.UserMessage("查询命令名称：{cmdName}\n附加问题：{question}"),
	)
}

func getWorkPlatform(enablePlatformPerception bool) string {
	// 获取配置信息，是否启用平台感知
	if !enablePlatformPerception {
		return ""
	}
	systemInfo, err := common.GetSystemInfo()
	if err != nil {
		logger.Errorf("get system info failed: %v\n", err)
	}
	shellPlatform, err := common.GetShellPlatform()
	if err != nil {
		logger.Errorf("get shell platform failed: %v\n", err)
	}
	workPlatform = strings.Replace(workPlatform, "{systemInfo}", systemInfo, -1)
	workPlatform = strings.Replace(workPlatform, "{shellPlatform}", shellPlatform, -1)
	return workPlatform
}

func getWorkUserAndDir(enableWorkUserAndDir bool) string {
	if !enableWorkUserAndDir {
		return ""
	}
	user, err := common.GetUser()
	if err != nil {
		logger.Errorf("get user failed: %v\n", err)
	}
	pwd, err := common.GetPwd()
	if err != nil {
		logger.Errorf("get pwd failed: %v\n", err)
	}
	workUserAndDir = strings.Replace(workUserAndDir, "{workUser}", user, -1)
	workUserAndDir = strings.Replace(workUserAndDir, "{workDir}", pwd, -1)
	return workUserAndDir
}

func getWorkFlow(enablePlatformPerception bool, enableWorkUserAndDir bool) string {
	// 重置所有步骤为默认状态
	for i := range workFlowSteps {
		workFlowSteps[i].Enable = workFlowSteps[i].Default
	}

	// 根据enable参数设置工作流程步骤的启用状态
	if !enableWorkUserAndDir {
		workFlowSteps[0].Enable = false
	}
	if !enablePlatformPerception {
		workFlowSteps[1].Enable = false
	}
	var steps []string
	for i, step := range workFlowSteps {
		if step.Enable {
			steps = append(steps, fmt.Sprintf("  %d. %s", i+1, step.Step))
		}
	}
	return "- 工作流程:\n" + strings.Join(steps, "\n")
}

func CreateOnceMessagesFromTemplate(cmdName, question string, enableExplain bool, enableExtendParams bool, enablePlatformPerception bool, enableWorkUserAndDir bool) []*schema.Message {
	template := createTemplate()
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"baseInfo":          baseInfo,
		"workFlow":          getWorkFlow(enablePlatformPerception, enableWorkUserAndDir),
		"workPlatform":      getWorkPlatform(enablePlatformPerception),
		"workUserAndDir":    getWorkUserAndDir(enableWorkUserAndDir),
		"answerDescription": answerDescription,
		"answerFormat":      answerFormat,
		"cmdName":           cmdName,
		"question":          question,
		// 对话历史
		"chatHistory": []*schema.Message{},
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	logger.Debugf("messages: %v\n", messages)
	return messages
}

func CreateMoreMessagesFromTemplate(cmdName, question string, chatHistory []*schema.Message, enableExplain bool, enableExtendParams bool, enablePlatformPerception bool, enableWorkUserAndDir bool) []*schema.Message {
	template := createTemplate()
	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"baseInfo":          baseInfo,
		"workFlow":          getWorkFlow(enablePlatformPerception, enableWorkUserAndDir),
		"workPlatform":      getWorkPlatform(enablePlatformPerception),
		"workUserAndDir":    getWorkUserAndDir(enableWorkUserAndDir),
		"answerDescription": answerDescription,
		"answerFormat":      answerFormat,
		"cmdName":           cmdName,
		"question":          question,
		// 对话历史
		"chatHistory": chatHistory,
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	logger.Debugf("messages: %v\n", messages)
	return messages
}
