package setup

import (
	"path/filepath"

	"wen-ai-cli/assets"

	"github.com/gookit/i18n"
)

// InitLang 初始化多语言配置
func InitLang() {
	config := GetConfig()
	targetLangDir := GetLangDir()
	langFiles, err := filepath.Glob(filepath.Join(targetLangDir, "*.ini"))
	if err != nil {
		panic("获取语言文件失败: " + err.Error())
	}
	if len(langFiles) == 0 {
		if err := assets.CopyLangFiles(targetLangDir); err != nil {
			panic("复制语言文件失败: " + err.Error())
		}
	}

	defaultLang := config.DefaultLang
	languages := map[string]string{
		"zh-CN": "简体中文",
		"en":    "English",
	}

	// 这里直接初始化的默认实例
	i18n.Init(targetLangDir, defaultLang, languages)
}
