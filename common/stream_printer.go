/*
Stream Printer for Formatted Text Output

This package provides functionality to process streaming text fragments (like those from OpenAI API)
and display them in the terminal with special formatting based on content type:

1. Markdown headings (lines starting with #) - Displayed in RED
2. Ordered list markers (like "1." at start of line) - Displayed in BLUE
3. Unordered list markers (like "*" at start of line) - Displayed in BLUE
4. Code blocks (content between ``` markers) - Displayed in GRAY

Features:
- Handles text fragments of varying sizes (from single chars to larger chunks)
- Maintains proper formatting even when receiving small fragments
- Buffers content until complete lines are available
- Provides typewriter-like effects for streaming text output
- Can be easily integrated with any streaming API
- Correctly handles multibyte characters like Chinese, Japanese, etc.

Usage:
- Use ProcessFragment() to handle individual text chunks as they arrive
- Use ProcessStream() to process a channel of incoming text fragments
- Call Flush() at the end to display any remaining buffered content

Sample integration with OpenAI API:
```go
streamChan := make(chan string)
go ProcessStream(streamChan)

// In your API response handler:
for chunk := range apiResponse {
    streamChan <- chunk
}
close(streamChan)
}
*/

package common

import (
	// 移除未使用的bytes包
	// "bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/fatih/color"
)

// StreamPrinter handles streaming text with special formatting
type StreamPrinter struct {
	// 移除未使用的buffer字段
	// buffer bytes.Buffer

	// State tracking
	lineStart            bool
	inHeading            bool
	inOrderedList        bool
	inUnorderedList      bool
	inCodeBlock          bool // 是否在代码块内
	codeBlockMarkerCount int  // 代码块标记计数

	// Temporary buffers
	headingBuffer string
	listBuffer    string
	codeBuffer    string // 用于识别代码块标记

	// Regular expressions
	headingPattern   *regexp.Regexp
	orderedPattern   *regexp.Regexp
	unorderedPattern *regexp.Regexp

	// Border decoration tracking
	isFirstLine       bool
	hasOutputContent  bool
	isLastLineWritten bool

	// Border characters - 边框字符配置
	headerChar      string // 顶部边框字符
	normalLineChar  string // 普通行边框字符
	headingLineChar string // 标题行边框字符
	footerChar      string // 底部边框字符
	headerText      string // 头部文本
	footerText      string // 底部文本

	// Border colors - 边框颜色配置
	headerColor      *color.Color // 顶部边框颜色
	normalLineColor  *color.Color // 普通行边框颜色
	headingLineColor *color.Color // 标题行边框颜色
	footerColor      *color.Color // 底部边框颜色
	headerTextColor  *color.Color // 头部文本颜色
	footerTextColor  *color.Color // 底部文本颜色

	// Content colors - 内容颜色配置
	headingColor       *color.Color // 标题文本颜色
	orderedListColor   *color.Color // 有序列表颜色
	unorderedListColor *color.Color // 无序列表颜色
	codeBlockColor     *color.Color // 代码块颜色
}

// NewStreamPrinter creates a new printer instance with default border style
func NewStreamPrinter() *StreamPrinter {
	return &StreamPrinter{
		lineStart:            true,
		inHeading:            false,
		inOrderedList:        false,
		inUnorderedList:      false,
		inCodeBlock:          false,
		codeBlockMarkerCount: 0,
		headingBuffer:        "",
		listBuffer:           "",
		codeBuffer:           "",
		headingPattern:       regexp.MustCompile(`^#+\s`),
		orderedPattern:       regexp.MustCompile(`^\d+\.\s`),
		unorderedPattern:     regexp.MustCompile(`^\*\s`),
		isFirstLine:          true,
		hasOutputContent:     false,
		isLastLineWritten:    false,
		// 默认边框字符
		headerChar:      "╭──",
		normalLineChar:  "│",
		headingLineChar: "│──",
		footerChar:      "╰──",
		headerText:      "WenAI CLI",
		footerText:      "END",
		// 默认边框颜色
		headerColor:      color.New(color.FgCyan),
		normalLineColor:  color.New(color.FgWhite),
		headingLineColor: color.New(color.FgYellow),
		footerColor:      color.New(color.FgCyan),
		headerTextColor:  color.New(color.FgHiWhite), // 默认头部文本颜色
		footerTextColor:  color.New(color.FgHiWhite), // 默认底部文本颜色

		// 默认内容颜色
		headingColor:       color.New(color.FgRed),
		orderedListColor:   color.New(color.FgBlue),
		unorderedListColor: color.New(color.FgBlue),
		codeBlockColor:     color.New(color.FgHiBlack),
	}
}

// NewStreamPrinterWithColors creates a printer with custom border style and colors
func NewStreamPrinterWithColors(
	// 边框字符
	headerChar, normalChar, headingChar, footerChar, headerText, footerText string,
	// 边框颜色
	headerColor, normalColor, headingColor, footerColor *color.Color,
	// 头部底部文本颜色
	headerTextColor, footerTextColor *color.Color,
	// 内容颜色
	contentHeadingColor, contentOrderedListColor, contentUnorderedListColor, contentCodeBlockColor *color.Color) *StreamPrinter {

	printer := NewStreamPrinter()

	// 设置自定义边框
	if headerChar != "" {
		printer.headerChar = headerChar
	}
	if normalChar != "" {
		printer.normalLineChar = normalChar
	}
	if headingChar != "" {
		printer.headingLineChar = headingChar
	}
	if footerChar != "" {
		printer.footerChar = footerChar
	}
	if headerText != "" {
		printer.headerText = headerText
	}
	if footerText != "" {
		printer.footerText = footerText
	}

	// 设置自定义边框颜色
	if headerColor != nil {
		printer.headerColor = headerColor
	}
	if normalColor != nil {
		printer.normalLineColor = normalColor
	}
	if headingColor != nil {
		printer.headingLineColor = headingColor
	}
	if footerColor != nil {
		printer.footerColor = footerColor
	}

	// 设置自定义文本颜色
	if headerTextColor != nil {
		printer.headerTextColor = headerTextColor
	}
	if footerTextColor != nil {
		printer.footerTextColor = footerTextColor
	}

	// 设置自定义内容颜色
	if contentHeadingColor != nil {
		printer.headingColor = contentHeadingColor
	}
	if contentOrderedListColor != nil {
		printer.orderedListColor = contentOrderedListColor
	}
	if contentUnorderedListColor != nil {
		printer.unorderedListColor = contentUnorderedListColor
	}
	if contentCodeBlockColor != nil {
		printer.codeBlockColor = contentCodeBlockColor
	}

	return printer
}

// NewStreamPrinterWithBorder creates a printer with custom border style
func NewStreamPrinterWithBorder(header, normal, heading, footer, headerText, footerText string) *StreamPrinter {
	return NewStreamPrinterWithColors(
		header, normal, heading, footer, headerText, footerText,
		nil, nil, nil, nil, // 默认边框颜色
		nil, nil, // 默认文本颜色
		nil, nil, nil, nil, // 默认内容颜色
	)
}

// ProcessFragment handles incoming text fragments
func (p *StreamPrinter) ProcessFragment(fragment string) {
	// Process character by character for other fragments
	for _, char := range fragment {
		p.processChar(char)
	}
}

// processChar handles each character with proper formatting
func (p *StreamPrinter) processChar(char rune) {
	// Write header on first content
	if !p.hasOutputContent {
		p.headerColor.Printf("%s ", p.headerChar)
		p.headerTextColor.Printf("%s \n", p.headerText)
		p.normalLineColor.Printf("%s\n", p.normalLineChar)
		p.hasOutputContent = true
		p.lineStart = true
	}

	// 对于换行符的特殊处理
	if char == '\n' {
		// 打印当前行的结束
		fmt.Println()

		// 关键改变：立即为下一行打印边框字符，即使这是个空行
		if p.inCodeBlock {
			p.normalLineColor.Printf("%s ", p.normalLineChar)
		} else {
			p.normalLineColor.Printf("%s ", p.normalLineChar)
		}

		// 重置状态标志，但边框已经打印
		p.lineStart = false
		p.inHeading = false
		p.inOrderedList = false
		p.inUnorderedList = false
		return
	}

	// 处理行首的边框字符
	if p.lineStart {
		// 行首标志重置，但不再需要在这里打印边框（已在换行符处理时打印）
		p.lineStart = false
	}

	// 处理代码块标记
	if char == '`' {
		p.codeBuffer += string(char)

		// 检查是否已经达到三个反引号（代码块标记）
		if len(p.codeBuffer) == 3 {
			if p.inCodeBlock {
				// 结束代码块
				p.codeBlockColor.Print("```")
				p.inCodeBlock = false
				p.codeBlockMarkerCount = 0
			} else {
				// 开始代码块
				p.codeBlockColor.Print("```")
				p.inCodeBlock = true
				p.codeBlockMarkerCount = 3
			}
			p.codeBuffer = ""
			return
		}
		return
	} else if len(p.codeBuffer) > 0 {
		// 不是完整的代码块标记，打印并重置
		if p.inCodeBlock {
			p.codeBlockColor.Print(p.codeBuffer)
			p.codeBlockColor.Print(string(char))
		} else {
			fmt.Print(p.codeBuffer)
			fmt.Print(string(char))
		}
		p.codeBuffer = ""
		return
	}

	// 如果在代码块内，使用自定义代码颜色显示
	if p.inCodeBlock {
		p.codeBlockColor.Print(string(char))
		return
	}

	// Handle character based on context
	if p.inHeading {
		// In a heading
		p.handleHeading(char)
	} else if p.inOrderedList {
		// In an ordered list
		p.handleOrderedList(char)
	} else if p.inUnorderedList {
		// In an unordered list
		p.handleUnorderedList(char)
	} else if char == '#' {
		// 关键改变：当检测到#字符时，立即替换边框
		if p.headingBuffer == "" {
			// 如果是第一个#，替换边框
			fmt.Print("\r")
			p.headingLineColor.Printf("%s ", p.headingLineChar)
		}
		// Start of potential heading (或继续收集多个#)
		p.headingBuffer += string(char)
	} else if char == '*' && p.listBuffer == "" {
		// Start of potential unordered list
		p.listBuffer += string(char)
	} else if (char >= '0' && char <= '9') && p.listBuffer == "" {
		// Start of potential ordered list
		p.listBuffer += string(char)
	} else if char == '.' && p.listBuffer != "" && regexp.MustCompile(`^\d+$`).MatchString(p.listBuffer) {
		// Part of an ordered list marker
		p.listBuffer += string(char)
	} else if char == ' ' {
		// Check if we have a heading, ordered list, or unordered list
		if p.headingBuffer != "" {
			// Confirm heading
			p.inHeading = true
			// 不打印标题标识符和空格
			p.headingBuffer = ""
		} else if p.listBuffer == "*" {
			// Confirm unordered list
			p.inUnorderedList = true
			p.unorderedListColor.Print(p.listBuffer)
			p.unorderedListColor.Print(string(char))
			p.listBuffer = ""
		} else if regexp.MustCompile(`^\d+\.$`).MatchString(p.listBuffer) {
			// Confirm ordered list
			p.inOrderedList = true
			p.orderedListColor.Print(p.listBuffer)
			p.orderedListColor.Print(string(char))
			p.listBuffer = ""
		} else {
			// Just a regular space
			if p.headingBuffer != "" {
				fmt.Print(p.headingBuffer)
				p.headingBuffer = ""
			}
			if p.listBuffer != "" {
				fmt.Print(p.listBuffer)
				p.listBuffer = ""
			}
			fmt.Print(string(char))
		}
	} else {
		// Not a special format, print any buffered content and the current char
		if p.headingBuffer != "" {
			fmt.Print(p.headingBuffer)
			p.headingBuffer = ""
		}
		if p.listBuffer != "" {
			fmt.Print(p.listBuffer)
			p.listBuffer = ""
		}
		fmt.Print(string(char))
	}
}

// handleHeading processes characters within a heading
func (p *StreamPrinter) handleHeading(char rune) {
	// 在标题内部处理字符，使用自定义标题颜色
	p.headingColor.Print(string(char))
}

// handleOrderedList processes characters in an ordered list
func (p *StreamPrinter) handleOrderedList(char rune) {
	// 在有序列表内部处理字符
	fmt.Print(string(char))
}

// handleUnorderedList processes characters in an unordered list
func (p *StreamPrinter) handleUnorderedList(char rune) {
	// 在无序列表内部处理字符
	fmt.Print(string(char))
}

// Flush outputs any remaining content in buffers
func (p *StreamPrinter) Flush() {
	// Print any remaining content in buffers
	if p.headingBuffer != "" {
		fmt.Print(p.headingBuffer)
	}
	if p.listBuffer != "" {
		fmt.Print(p.listBuffer)
	}
	if p.codeBuffer != "" {
		fmt.Print(p.codeBuffer)
	}

	// 无论如何，确保结束前有换行
	fmt.Print("\r")
	p.footerColor.Printf("%s\n", p.normalLineChar)
	p.footerColor.Printf("%s ", p.footerChar)
	p.footerTextColor.Printf("%s\n", p.footerText)
}

// SimulateStreamInput demonstrates the stream printer with example content
func SimulateStreamInput() {
	printer := NewStreamPrinter()

	// 打印说明文本
	fmt.Println("\n=== 逐字符流式输出演示 ===")

	// 测试片段，包含中文字符，模拟字符流式输出
	fragments := []string{
		"#", " ", "一级标题示例", "\n",
		"这是普通文本，包含一个链接。", "\n",
		"\n", // 添加一个空行测试
		"\n", // 连续的空行测试
		"0. 没有\n",
		"1", ".", " ", "这", "是", "有", "序", "列", "表", "的", "第", "一", "项", "\n",
		"2", ".", " ", "这", "是", "有", "序", "列", "表", "的", "第", "二", "项", "\n",
		"*", " ", "这", "是", "无", "序", "列", "表", "的", "第", "一", "项", "\n",
		"*", " ", "这", "是", "无", "序", "列", "表", "的", "第", "二", "项", "\n",
		"## ", "二级标题", "\n",
		"以下是一个代码示例:", "\n",
		"`", "`", "`", "\n",
		"f", "u", "n", "c", " ", "g", "e", "t", "N", "a", "m", "e", "(", ")", " ", "{", "\n",
		"\n", // 代码块内空行
		"\n", // 代码块内连续空行
		" ", " ", " ", " ", "r", "e", "t", "u", "r", "n", " ", "\"", "逐", "字", "符", "演", "示", "\"", "\n",
		"}", "\n",
		"`", "`", "`", "\n",
		"### ", "三级标题", "\n",
		"继续显示普通文本。", "\n",
		"\n", // 再加一个空行
		"\n", // 最后两个连续空行
	}

	// Process each fragment with a delay
	for _, fragment := range fragments {
		printer.ProcessFragment(fragment)
		time.Sleep(50 * time.Millisecond) // Typing speed simulation
	}

	// Flush any remaining content
	printer.Flush()
}

// Process a stream from any source that provides a text channel
func ProcessStream(stream <-chan string) {
	printer := NewStreamPrinter()

	// Process each fragment as it arrives
	for fragment := range stream {
		// 特殊处理可能的空行
		for _, char := range fragment {
			printer.processChar(char)
		}
	}

	// Flush any remaining content
	printer.Flush()
}

// Example showing usage with a real API call (though simulated here)
func ExampleWithRealAPI() {
	fmt.Println("\n\n=== 与OpenAI风格API集成的示例 ===")

	// Create a channel to receive fragments
	streamChan := make(chan string)

	// Start a goroutine to process the stream
	go func() {
		ProcessStream(streamChan)
	}()

	// Simulate API sending fragments (in Chinese)
	sentences := []string{
		"# API响应标题\n\n\n", // 三个连续换行符
		"以下是您请求的编程语言相关数据：\n\n",
		"1. Python是一种高级、通用型编程语言。\n",
		"2. JavaScript通常用于Web开发。\n",
		"3. Go（或Golang）以其性能和简洁性著称。\n\n",
		"代码示例:\n```\npackage main\n\n\nimport \"fmt\"\n\nfunc main() {\n\n    fmt.Println(\"Hello, 世界\")\n\n}\n```\n\n", // 代码块内有空行
		"获取更多信息，请访问相关链接。\n\n",                                                                                           // 确保最后有连续的空行
	}

	// Send fragments with delays to simulate API streaming
	for _, sentence := range sentences {
		// 转换为[]rune确保正确处理中文字符
		runeText := []rune(sentence)
		remaining := runeText

		for len(remaining) > 0 {
			// Random fragment size between 1-4 runes
			chunkSize := 1 + (time.Now().Nanosecond() % 4)
			if chunkSize > len(remaining) {
				chunkSize = len(remaining)
			}

			// 确保每个片段都是有效的UTF-8字符序列
			fragment := string(remaining[:chunkSize])

			// Send fragment to channel
			streamChan <- fragment
			remaining = remaining[chunkSize:]

			// Simulate network delay
			time.Sleep(time.Duration(20+(time.Now().Nanosecond()%60)) * time.Millisecond)
		}
	}

	// Close the channel when done
	close(streamChan)

	// Wait for processing to complete
	time.Sleep(100 * time.Millisecond)
}

// 添加一个使用自定义边框的示例函数
func SimulateWithCustomBorder() {
	// 创建使用自定义边框的打印器
	printer := NewStreamPrinterWithBorder(
		"╭┈┈",  // 顶部边框
		"|",    // 普通行边框
		"├┈┈",  // 标题行边框
		"╰┈┈",  // 底部边框
		"用户问题", // 头部文本
		"完成",   // 底部文本
	)

	// 打印说明文本
	fmt.Println("\n=== 自定义边框示例 ===")

	// 使用相同的测试片段
	fragments := []string{
		"#", " ", "一级标题示例", "\n",
		"这是普通文本，使用自定义边框。", "\n",
		"## ", "二级标题", "\n",
		"这是使用自定义边框的内容。", "\n",
	}

	// Process each fragment with a delay
	for _, fragment := range fragments {
		printer.ProcessFragment(fragment)
		time.Sleep(50 * time.Millisecond) // Typing speed simulation
	}

	// Flush any remaining content
	printer.Flush()
}

// 添加一个使用自定义边框颜色的示例函数
func SimulateWithCustomColors() {
	// 创建使用自定义边框和颜色的打印器
	printer := NewStreamPrinterWithColors(
		"★━━",                        // 顶部边框
		"│",                          // 普通行边框
		"☆━━",                        // 标题行边框
		"★━━",                        // 底部边框
		"彩色边框示例",                     // 头部文本
		"结束",                         // 底部文本
		color.New(color.FgHiMagenta), // 顶部边框颜色 - 亮紫色
		color.New(color.FgGreen),     // 普通行边框颜色 - 绿色
		color.New(color.FgHiYellow),  // 标题行边框颜色 - 亮黄色
		color.New(color.FgHiMagenta), // 底部边框颜色 - 亮紫色
		color.New(color.FgHiWhite),   // 头部文本颜色 - 亮白色
		color.New(color.FgHiWhite),   // 底部文本颜色 - 亮白色
		color.New(color.FgRed),       // 标题文本颜色 - 红色
		color.New(color.FgBlue),      // 有序列表颜色 - 蓝色
		color.New(color.FgBlue),      // 无序列表颜色 - 蓝色
		color.New(color.FgHiBlack),   // 代码块颜色 - 灰色
	)

	// 打印说明文本
	fmt.Println("\n=== 自定义边框颜色示例 ===")

	// 使用测试片段
	fragments := []string{
		"#", " ", "彩色边框标题", "\n",
		"这是普通文本，使用自定义颜色边框。", "\n",
		"## ", "二级标题", "\n",
		"* 颜色可以通过参数单独设置", "\n",
		"* 可以设置不同部分使用不同颜色", "\n",
	}

	// Process each fragment with a delay
	for _, fragment := range fragments {
		printer.ProcessFragment(fragment)
		time.Sleep(50 * time.Millisecond) // Typing speed simulation
	}

	// Flush any remaining content
	printer.Flush()
}

// 添加一个使用自定义内容颜色的示例函数
func SimulateWithCustomContentColors() {
	// 创建使用自定义内容颜色的打印器
	printer := NewStreamPrinterWithColors(
		"┏━━",                        // 顶部边框
		"┃",                          // 普通行边框
		"┣━━",                        // 标题行边框
		"┗━━",                        // 底部边框
		"自定义内容颜色示例",                  // 头部文本
		"结束",                         // 底部文本
		color.New(color.FgHiCyan),    // 顶部边框颜色
		color.New(color.FgHiCyan),    // 普通行边框颜色
		color.New(color.FgHiCyan),    // 标题行边框颜色
		color.New(color.FgHiCyan),    // 底部边框颜色
		color.New(color.FgHiMagenta), // 头部文本颜色 - 亮紫色
		color.New(color.FgHiMagenta), // 底部文本颜色 - 亮紫色
		color.New(color.FgHiGreen),   // 标题文本颜色 - 亮绿色
		color.New(color.FgHiMagenta), // 有序列表颜色 - 亮紫色
		color.New(color.FgHiYellow),  // 无序列表颜色 - 亮黄色
		color.New(color.FgHiBlue),    // 代码块颜色 - 亮蓝色
	)

	// 打印说明文本
	fmt.Println("\n=== 自定义内容颜色示例 ===")

	// 使用测试片段
	fragments := []string{
		"#", " ", "自定义绿色标题", "\n",
		"这是普通文本内容。", "\n",
		"1", ".", " ", "紫色有序列表项目", "\n",
		"2", ".", " ", "也是紫色的有序列表", "\n",
		"*", " ", "黄色无序列表项目", "\n",
		"*", " ", "另一个黄色无序列表项目", "\n",
		"以下是代码示例：", "\n",
		"`", "`", "`", "\n",
		"package main\n\nfunc main() {\n    // 这是蓝色代码块\n    fmt.Println(\"Hello World\")\n}\n",
		"`", "`", "`", "\n",
	}

	// Process each fragment with a delay
	for _, fragment := range fragments {
		printer.ProcessFragment(fragment)
		time.Sleep(50 * time.Millisecond) // Typing speed simulation
	}

	// Flush any remaining content
	printer.Flush()
}

// 添加一个简化版的自定义颜色构造函数
func NewColoredStreamPrinter(
	borderColor, textColor, headingColor, listColor, codeColor *color.Color,
	headerText, footerText string) *StreamPrinter {

	// 使用相同的边框字符，但允许自定义颜色
	return NewStreamPrinterWithColors(
		"┏━━", "┃", "┣━━", "┗━━", // 使用默认边框字符
		headerText, footerText, // 自定义文本
		borderColor, borderColor, borderColor, borderColor, // 所有边框使用相同颜色
		textColor, textColor, // 头部和底部文本颜色
		headingColor, listColor, listColor, codeColor, // 内容颜色
	)
}

// 添加一个使用简化版构造函数的示例
func SimulateWithSimpleColors() {
	// 使用简化版构造函数创建打印器
	printer := NewColoredStreamPrinter(
		color.New(color.FgHiBlue),  // 边框颜色 - 亮蓝色
		color.New(color.FgHiCyan),  // 文本颜色 - 亮青色
		color.New(color.FgHiRed),   // 标题颜色 - 亮红色
		color.New(color.FgHiGreen), // 列表颜色 - 亮绿色
		color.New(color.FgYellow),  // 代码颜色 - 黄色
		"简化版颜色配置",                  // 头部文本
		"结束",                       // 底部文本
	)

	// 打印说明文本
	fmt.Println("\n=== 简化版颜色配置示例 ===")
	fmt.Println("(使用亮青色的头部和底部文本)")

	// 使用测试片段
	fragments := []string{
		"#", " ", "红色标题", "\n",
		"普通文本内容。", "\n",
		"1", ".", " ", "绿色有序列表", "\n",
		"*", " ", "绿色无序列表", "\n",
		"```\n// 黄色代码块\nfunc example() {\n    return\n}\n```\n",
	}

	// Process each fragment with a delay
	for _, fragment := range fragments {
		printer.ProcessFragment(fragment)
		time.Sleep(50 * time.Millisecond) // Typing speed simulation
	}

	// Flush any remaining content
	printer.Flush()
}
