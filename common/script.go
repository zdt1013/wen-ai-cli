package common

import (
	"errors"
	"log/slog"
	"strings"
	"wen-ai-cli/logger"
	"wen-ai-cli/model"
	"wen-ai-cli/validate"

	"github.com/gookit/i18n"
	"github.com/manifoldco/promptui"
)

// HandleScriptAdjustment 处理脚本微调运行逻辑
func HandleScriptAdjustment(defaultScript string) (string, bool) {
	validateFn := func(input string) error {
		if len(input) < 1 {
			return errors.New(i18n.Dtr("paramEmptyError"))
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    ">",
		Validate: validateFn,
		Default:  defaultScript,
	}

	result, err := prompt.Run()

	if err != nil {
		logger.Errorf("Prompt failed %v", err)
		return "", false
	}

	logger.Debugf(i18n.Dtr("adjustedScript"), slog.String("script", result))

	// 要求用户确认是否运行
	confirm := promptui.Select{
		HideHelp: true,
		Label:    i18n.Dtr("confirmRunScript"),
		Items:    []string{i18n.Dtr("yes"), i18n.Dtr("no")},
	}
	_, confirmResult, err := confirm.Run()
	if err != nil {
		logger.Errorf("Prompt failed %v", err)
		return "", false
	}

	return result, confirmResult == i18n.Dtr("yes")
}

// ConfirmExecution 确认是否执行脚本
func ConfirmExecution() (bool, error) {
	confirm := promptui.Select{
		HideHelp: true,
		Label:    i18n.Dtr("confirmRunScript"),
		Items:    []string{i18n.Dtr("yes"), i18n.Dtr("no")},
	}
	_, confirmResult, err := confirm.Run()
	if err != nil {
		logger.Errorf("Prompt failed %v", err)
		return false, err
	}

	return confirmResult == i18n.Dtr("yes"), nil
}

// HandleParamsCompletion 处理参数补全运行逻辑
func HandleParamsCompletion(hiddenParams *model.HiddenParams) (string, bool) {
	// 遍历参数获取用户输入
	for i := range hiddenParams.NeedFillParams {
		// 使用指针引用，确保修改能保存到原始数据
		param := &hiddenParams.NeedFillParams[i]
		// 获取参数类型
		param_type := param.Type
		// 获取参数名称
		param_name := param.Param
		// 判断参数类型，做不同校验规则
		validateFn := func(input string) error {
			return validate.ValidateParam(input, param_type)
		}

		prompt := promptui.Prompt{
			Label:       param_name,
			Validate:    validateFn,
			HideEntered: true,
		}
		result, err := prompt.Run()

		if err != nil {
			logger.Errorf("Prompt failed %v", err)
			return "", false
		}
		// 回填参数值
		param.Value = result
	}

	// 替换参数
	shell_code := hiddenParams.ShellCode
	for _, param := range hiddenParams.NeedFillParams {
		// 构建可能的参数形式：原始参数名和<>包裹的参数名
		paramName := param.Param
		paramType := param.Type
		bracketParamName := "<" + paramName + "," + paramType + ">"
		// 替换参数
		shell_code = strings.Replace(shell_code, bracketParamName, param.Value, -1)
	}

	// 打印最终脚本
	logger.Debugf(i18n.Dtr("scriptToExecute"), shell_code)

	shouldExecute, err := ConfirmExecution()
	if err != nil {
		return "", false
	}

	return shell_code, shouldExecute
}
