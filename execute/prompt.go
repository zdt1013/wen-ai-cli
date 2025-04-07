package execute

import (
	"wen-ai-cli/setup"
	"wen-ai-cli/validate"

	"github.com/manifoldco/promptui"
)

func Prompt(label string, items []string) (string, error) {
	prompt := promptui.Select{
		HideHelp: true,
		Label:    label,
		Items:    items,
	}
	_, result, err := prompt.Run()
	return result, err
}

func InputString(label string) (string, error) {
	i18n := setup.GetI18n()
	// 定义输入验证函数
	validateFn := func(input string) error {
		return validate.ValidateParam(input, "string")
	}

	// 创建用户输入提示
	prompt := promptui.Prompt{
		Label:       i18n.UserInput,
		Validate:    validateFn,
		HideEntered: true,
	}
	// 获取用户输入
	input_quetion, err := prompt.Run()
	return input_quetion, err
}
