package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

// StreamPrinter 是一个流式打印对象，用于处理和打印流式文本
type StreamPrinter struct {
	buffer      string // 存储未完成的行
	inCodeBlock bool   // 是否在代码块内
	firstPrint  bool   // 是否是首次打印
	showHashTag bool   // 是否显示标题的#符号
	colorCode   bool   // 是否对代码块内容着色

	// 内容颜色
	titleColor       *color.Color // 标题颜色
	listColor        *color.Color // 列表标记颜色
	codeColor        *color.Color // 代码块边界颜色
	codeContentColor *color.Color // 代码块内容颜色
	normalColor      *color.Color // 普通文本颜色

	// 边框字符
	headerChar      string // 头部边框字符
	normalLineChar  string // 普通行边框字符
	headingLineChar string // 标题行边框字符
	footerChar      string // 尾部边框字符
	headerText      string // 头部文本
	footerText      string // 尾部文本

	// 边框颜色
	headerColor      *color.Color // 头部边框颜色
	normalLineColor  *color.Color // 普通行边框颜色
	headingLineColor *color.Color // 标题行边框颜色
	footerColor      *color.Color // 尾部边框颜色
	headerTextColor  *color.Color // 头部文本颜色
	footerTextColor  *color.Color // 尾部文本颜色
}

// NewStreamPrinter 创建一个新的流式打印器
func NewStreamPrinter() *StreamPrinter {
	return &StreamPrinter{
		buffer:      "",
		inCodeBlock: false,
		firstPrint:  true,
		showHashTag: true, // 默认显示标题的#符号
		colorCode:   true, // 默认对代码块内容着色

		// 内容颜色
		titleColor:       color.New(color.FgBlue),    // 标题使用蓝色
		listColor:        color.New(color.FgMagenta), // 列表标记使用紫色
		codeColor:        color.New(color.FgHiBlack), // 代码块边界使用灰色
		codeContentColor: color.New(color.FgHiCyan),  // 代码块内容使用青色
		normalColor:      color.New(color.Reset),     // 普通文本使用默认颜色

		// 边框字符
		headerChar:      "╭──",
		normalLineChar:  "│",
		headingLineChar: "│──",
		footerChar:      "╰──",
		headerText:      "WenAI CLI",
		footerText:      "END",

		// 边框颜色
		headerColor:      color.New(color.FgCyan),
		normalLineColor:  color.New(color.FgCyan),
		headingLineColor: color.New(color.FgCyan),
		footerColor:      color.New(color.FgCyan),
		headerTextColor:  color.New(color.FgHiWhite),
		footerTextColor:  color.New(color.FgHiWhite),
	}
}

// SetShowHashTag 设置是否显示标题的#符号
func (sp *StreamPrinter) SetShowHashTag(show bool) {
	sp.showHashTag = show
}

// SetColorCode 设置是否对代码块内容着色
func (sp *StreamPrinter) SetColorCode(color bool) {
	sp.colorCode = color
}

// SetHeaderText 设置头部文本
func (sp *StreamPrinter) SetHeaderText(text string) {
	sp.headerText = text
}

// SetFooterText 设置尾部文本
func (sp *StreamPrinter) SetFooterText(text string) {
	sp.footerText = text
}

// NewStreamPrinterWithOptions 创建一个带有显示选项的流式打印器
func NewStreamPrinterWithOptions(showHashTag bool) *StreamPrinter {
	sp := NewStreamPrinter()
	sp.showHashTag = showHashTag
	return sp
}

// NewStreamPrinterWithFullOptions 创建一个带有完整选项的流式打印器
func NewStreamPrinterWithFullOptions(showHashTag bool, colorCode bool) *StreamPrinter {
	sp := NewStreamPrinter()
	sp.showHashTag = showHashTag
	sp.colorCode = colorCode
	return sp
}

// NewStreamPrinterWithTextOptions 创建一个带有文本选项的流式打印器
func NewStreamPrinterWithTextOptions(headerText string, footerText string) *StreamPrinter {
	sp := NewStreamPrinter()
	sp.headerText = headerText
	sp.footerText = footerText
	return sp
}

// NewStreamPrinterWithAllOptions 创建一个带有所有选项的流式打印器
func NewStreamPrinterWithAllOptions(showHashTag bool, colorCode bool, headerText string, footerText string) *StreamPrinter {
	sp := NewStreamPrinter()
	sp.showHashTag = showHashTag
	sp.colorCode = colorCode
	sp.headerText = headerText
	sp.footerText = footerText
	return sp
}

// NewStreamPrinterWithColors 创建一个自定义颜色的流式打印器
func NewStreamPrinterWithColors(titleAttr, listAttr, codeAttr, codeContentAttr, normalAttr color.Attribute) *StreamPrinter {
	sp := NewStreamPrinter()
	sp.titleColor = color.New(titleAttr)
	sp.listColor = color.New(listAttr)
	sp.codeColor = color.New(codeAttr)
	sp.codeContentColor = color.New(codeContentAttr)
	sp.normalColor = color.New(normalAttr)
	return sp
}

// NewStreamPrinterWithBorder 创建一个自定义边框的流式打印器
func NewStreamPrinterWithBorder(
	headerChar, normalLineChar, headingLineChar, footerChar,
	headerText, footerText string,
	headerColor, normalLineColor, headingLineColor, footerColor,
	headerTextColor, footerTextColor color.Attribute) *StreamPrinter {

	sp := NewStreamPrinter()

	sp.headerChar = headerChar
	sp.normalLineChar = normalLineChar
	sp.headingLineChar = headingLineChar
	sp.footerChar = footerChar
	sp.headerText = headerText
	sp.footerText = footerText

	sp.headerColor = color.New(headerColor)
	sp.normalLineColor = color.New(normalLineColor)
	sp.headingLineColor = color.New(headingLineColor)
	sp.footerColor = color.New(footerColor)
	sp.headerTextColor = color.New(headerTextColor)
	sp.footerTextColor = color.New(footerTextColor)

	return sp
}

// 打印头部
func (sp *StreamPrinter) printHeader() {
	sp.headerColor.Print(sp.headerChar)
	sp.headerTextColor.Print(" " + sp.headerText)
	fmt.Println()
	sp.firstPrint = false
}

// 打印尾部
func (sp *StreamPrinter) printFooter() {
	sp.footerColor.Print(sp.footerChar)
	sp.footerTextColor.Print(" " + sp.footerText)
	fmt.Println()
}

// Print 接收一段文本并处理打印
func (sp *StreamPrinter) Print(text string) {
	// 如果是首次打印，先打印头部
	if sp.firstPrint {
		sp.printHeader()
	}

	// 将新的文本添加到缓冲区
	sp.buffer += text

	// 查找所有完整的行（以换行符结尾）
	for {
		// 查找第一个换行符的位置
		newlineIndex := strings.Index(sp.buffer, "\n")
		if newlineIndex == -1 {
			// 没有完整的行，保留所有内容在缓冲区中
			break
		}

		// 提取完整的行（不包括换行符）
		completeLine := sp.buffer[:newlineIndex]

		// 打印这一行（包括边框）
		sp.printLineWithBorder(completeLine)
		fmt.Println() // 打印换行符

		// 更新缓冲区，移除已处理的行（包括换行符）
		sp.buffer = sp.buffer[newlineIndex+1:]
	}
}

// Flush 刷新缓冲区中的所有内容，并打印尾部
func (sp *StreamPrinter) Flush() {
	// 如果是首次打印并且有内容，先打印头部
	if sp.firstPrint && sp.buffer != "" {
		sp.printHeader()
	}

	if sp.buffer != "" {
		sp.printLineWithBorder(sp.buffer)
		fmt.Println()
		sp.buffer = ""
	}

	// 打印尾部
	sp.printFooter()
}

// 处理并打印带边框的单行文本
func (sp *StreamPrinter) printLineWithBorder(line string) {
	// 确定是否是标题行
	isHeading := false
	if !sp.inCodeBlock && regexp.MustCompile(`^#{1,6}\s`).MatchString(line) {
		isHeading = true
	}

	// 打印行边框
	if isHeading {
		sp.headingLineColor.Print(sp.headingLineChar)
	} else {
		sp.normalLineColor.Print(sp.normalLineChar)
	}

	// 打印空格，美观
	fmt.Print(" ")

	// 打印实际内容
	sp.printFormattedContent(line)
}

// 处理并打印格式化内容（不带边框）
func (sp *StreamPrinter) printFormattedContent(line string) {
	// 检查是否是代码块边界
	if strings.HasPrefix(line, "```") {
		sp.codeColor.Print(line)
		sp.inCodeBlock = !sp.inCodeBlock
		return
	}

	// 如果在代码块内，使用代码内容颜色或普通颜色打印
	if sp.inCodeBlock {
		if sp.colorCode {
			sp.codeContentColor.Print(line)
		} else {
			sp.normalColor.Print(line)
		}
		return
	}

	// 处理标题行（# 开头）
	if match, _ := regexp.MatchString(`^#{1,6}\s`, line); match {
		if sp.showHashTag {
			// 显示全部内容，包括#符号
			sp.titleColor.Print(line)
		} else {
			// 不显示#符号，只显示标题内容
			re := regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
			parts := re.FindStringSubmatch(line)
			if len(parts) == 3 {
				// 直接打印标题内容，不显示#符号
				sp.titleColor.Print(parts[2])
			} else {
				// 如果解析失败，打印整行
				sp.titleColor.Print(line)
			}
		}
		return
	}

	// 处理有序列表（数字加点）
	if match, _ := regexp.MatchString(`^\d+\.\s`, line); match {
		// 找到数字和点的部分
		re := regexp.MustCompile(`^(\d+\.\s)(.*)$`)
		parts := re.FindStringSubmatch(line)
		if len(parts) == 3 {
			sp.listColor.Print(parts[1])
			sp.normalColor.Print(parts[2])
		} else {
			sp.normalColor.Print(line)
		}
		return
	}

	// 处理无序列表（*号）
	if strings.HasPrefix(line, "* ") {
		sp.listColor.Print("* ")
		sp.normalColor.Print(line[2:])
		return
	}

	// 其他普通文本行
	sp.normalColor.Print(line)
}
