package setup

import "github.com/gookit/i18n"

// I18n 定义了所有需要国际化的文本内容
type I18n struct {
	UserInput        string // 用户输入提示
	UserInputFormat  string // 用户输入格式提示
	ChatHelp         string // 聊天帮助信息
	SelectOperation  string // 选择操作提示
	FillParamsAndRun string // 填充参数并运行提示
	AdjustAndRun     string // 调整并运行提示
	YourChoice       string // 你的选择提示
	RunNow           string // 立即运行提示
	CanExecute       string // 可以执行提示
	Exit             string // 退出提示
	ParamEmptyError  string // 参数为空错误
	UrlInvalidError  string // URL无效错误
}

// i18nInstance 是 I18n 的单例实例
var i18nInstance *I18n

// GetI18n 获取 I18n 的单例实例
// 如果实例不存在则初始化一个新实例
func GetI18n() *I18n {
	if i18nInstance == nil {
		i18nInstance = &I18n{
			UserInput:        getDtr("userInput"),
			UserInputFormat:  getDtr("userInputFormat"),
			ChatHelp:         getDtr("chatHelp"),
			SelectOperation:  getDtr("selectOperation"),
			FillParamsAndRun: getDtr("fillParamsAndRun"),
			AdjustAndRun:     getDtr("adjustAndRun"),
			YourChoice:       getDtr("yourChoice"),
			RunNow:           getDtr("runNow"),
			CanExecute:       getDtr("canExecute"),
			Exit:             getDtr("exit"),
			ParamEmptyError:  getDtr("paramEmptyError"),
			UrlInvalidError:  getDtr("urlInvalidError"),
		}
	}
	return i18nInstance
}

// getDtr 获取指定key的国际化文本
func getDtr(key string) string {
	return i18n.Dtr(key)
}
