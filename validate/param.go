package validate

import (
	"errors"
	"net/url"
	"strings"
	"wen-ai-cli/setup"
)

// ValidateParam 根据参数类型验证输入值
func ValidateParam(input string, paramType string) error {
	i18n := setup.GetI18n()
	if len(input) < 1 {
		return errors.New(i18n.ParamEmptyError)
	}

	// 根据参数类型执行不同的验证
	switch strings.ToLower(paramType) {
	case "url":
		// URL类型验证
		_, err := url.Parse(input)
		if err != nil {
			return errors.New(i18n.UrlInvalidError)
		}

		// 检查URL是否包含协议
		if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			return errors.New(i18n.UrlInvalidError)
		}
	case "string":
		// 字符串类型验证 - 基本检查
		if strings.TrimSpace(input) == "" {
			return errors.New(i18n.ParamEmptyError)
		}
	}

	return nil
}
