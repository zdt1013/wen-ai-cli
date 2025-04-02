package model

type OpenAI struct {
	APIKey  string `mapstructure:"apiKey" json:"apiKey"`
	BaseURL string `mapstructure:"baseURL" json:"baseURL"`
	Model   string `mapstructure:"model" json:"model"`
}

type Console struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled"`
	Color   bool   `mapstructure:"color" json:"color"`
	Level   string `mapstructure:"level" json:"level"`
}

type File struct {
	Enabled    bool   `mapstructure:"enabled" json:"enabled"`
	Path       string `mapstructure:"path" json:"path"`
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize"`
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups"`
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge"`
	Level      string `mapstructure:"level" json:"level"`
}

type Logger struct {
	Console Console `mapstructure:"console" json:"console"`
	File    File    `mapstructure:"file" json:"file"`
}

type AnswerConfig struct {
	EnableExplain            bool `mapstructure:"enableExplain" json:"enableExplain"`
	EnableExtendParams       bool `mapstructure:"enableExtendParams" json:"enableExtendParams"`
	EnablePlatformPerception bool `mapstructure:"enablePlatformPerception" json:"enablePlatformPerception"`
}

type Config struct {
	DefaultLang  string       `mapstructure:"defaultLang" json:"defaultLang"`
	OpenAI       OpenAI       `mapstructure:"openai" json:"openai"`
	Logger       Logger       `mapstructure:"logger" json:"logger"`
	AnswerConfig AnswerConfig `mapstructure:"answerConfig" json:"answerConfig"`
}
