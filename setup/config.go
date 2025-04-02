package setup

import (
	"encoding/json"
	"os"
	"path/filepath"

	"wen-ai-cli/model"

	"github.com/gookit/config/v2"
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
)

var cfg *model.Config

// GetConfig 获取配置实例
func GetConfig() *model.Config {
	if cfg == nil {
		cfg = &model.Config{}
		loadConfig()
	}
	return cfg
}

// InitConfig 初始化配置
func InitConfig() {
	// 设置选项支持ENV变量解析：当获取的值为string类型时，会尝试解析其中的ENV变量
	config.WithOptions(config.ParseEnv)
	// 加载配置，可以同时传入多个文件
	configFilePath := GetConfigFilePath()
	if !fsutil.FileExists(configFilePath) {
		// 创建默认配置
		err := createDefaultConfig(configFilePath)
		if err != nil {
			panic("创建默认配置文件失败: " + err.Error())
		}
	}
	// 加载配置文件
	err := config.LoadFiles(configFilePath)
	if err != nil {
		panic("加载配置文件失败: " + err.Error())
	}
}

func createDefaultConfig(configFilePath string) error {
	// Config结构体转换为json，写入文件
	emptyCfg := model.Config{
		DefaultLang: "zh-CN",
		OpenAI: model.OpenAI{
			APIKey:  "",
			BaseURL: "",
			Model:   "",
		},
		Logger: model.Logger{
			Console: model.Console{
				Enabled: true,
				Color:   true,
				Level:   "info",
			},
			File: model.File{
				Enabled:    true,
				Path:       GetLogFilePath(),
				MaxSize:    100,
				MaxBackups: 30,
				MaxAge:     30,
				Level:      "debug",
			},
		}, AnswerConfig: model.AnswerConfig{
			EnableExplain:            true,
			EnableExtendParams:       true,
			EnablePlatformPerception: true,
		},
	}
	jsonData, err := json.Marshal(emptyCfg)
	if err != nil {
		panic("转换配置文件失败: " + err.Error())
	}
	fsutil.WriteFile(configFilePath, jsonData, 0644)
	return nil
}

// loadConfig 加载配置到结构体
func loadConfig() {
	err := config.BindStruct("", cfg)
	if err != nil {
		panic("解析配置文件失败: " + err.Error())
	}
}

// SaveConfig 保存配置
func SaveConfig(cfg *model.Config) {
	// 将配置转换为json
	jsonData, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		panic("转换配置文件失败: " + err.Error())
	}
	// 保存配置
	os.WriteFile(GetConfigFilePath(), jsonData, 0644)
}

// GetAppDir 获取配置目录
func GetAppDir() string {
	homeDir := sysutil.UserHomeDir()
	confDir := filepath.Join(homeDir, ".wenai")
	return confDir
}

// GetLangDir 获取语言目录
func GetLangDir() string {
	appDir := GetAppDir()
	return filepath.Join(appDir, "lang")
}

// GetConfigFilePath 获取配置文件路径
func GetConfigFilePath() string {
	appDir := GetAppDir()
	return filepath.Join(appDir, "conf.json")
}

// GetLogFilePath 获取日志文件路径
func GetLogFilePath() string {
	appDir := GetAppDir()
	return filepath.Join(appDir, "logs/app.log")
}
