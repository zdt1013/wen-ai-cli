package validate

import (
	"errors"
	"net/url"
	"strings"
)

// ValidateParam 根据参数类型验证输入值
func ValidateParam(input string, paramType string) error {
	if len(input) < 1 {
		return errors.New("参数不能为空")
	}

	// 根据参数类型执行不同的验证
	switch strings.ToLower(paramType) {
	case "url":
		// URL类型验证
		_, err := url.Parse(input)
		if err != nil {
			return errors.New("请输入有效的URL地址")
		}

		// 检查URL是否包含协议
		if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
			return errors.New("URL必须包含http://或https://协议")
		}
	case "string":
		// 字符串类型验证 - 基本检查
		if strings.TrimSpace(input) == "" {
			return errors.New("字符串不能只包含空白字符")
		}
		// 可以添加其他类型的验证逻辑
	}

	return nil
}
