package setup

import "github.com/gookit/i18n"

type I18n struct {
	UserInput        string
	UserInputFormat  string
	ChatHelp         string
	SelectOperation  string
	FillParamsAndRun string
	AdjustAndRun     string
	YourChoice       string
	RunNow           string
	CanExecute       string
	Exit             string
	ParamEmptyError  string
	UrlInvalidError  string
}

var i18nInstance *I18n

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

func getDtr(key string) string {
	return i18n.Dtr(key)
}
