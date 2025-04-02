package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"wen-ai-cli/setup"

	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志级别映射
var levelMap = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
	"fatal": slog.LevelError,
}

var (
	logger *slog.Logger
)

// MultiHandler 实现多处理器
type MultiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler 创建新的多处理器
func NewMultiHandler(handlers ...slog.Handler) slog.Handler {
	return &MultiHandler{handlers: handlers}
}

// Handle 实现 slog.Handler 接口
func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		// 只处理启用了该日志级别的处理器
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

// WithAttrs 实现 slog.Handler 接口
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

// WithGroup 实现 slog.Handler 接口
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}

// Enabled 实现 slog.Handler 接口
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// 检查是否有任何处理器启用了该级别
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// 获取日志级别值
func getLogLevelValue(levelName string, defaultLevel slog.Level) slog.Level {
	if levelName == "" {
		return defaultLevel
	}

	levelLower := strings.ToLower(levelName)
	if val, ok := levelMap[levelLower]; ok {
		return val
	}

	return defaultLevel
}

// 初始化日志配置
func initLogger() {
	// 从配置文件读取日志配置
	logConfig := setup.GetConfig().Logger

	// 设置日志级别
	consoleLevel := getLogLevelValue(logConfig.Console.Level, slog.LevelInfo)
	fileLevel := getLogLevelValue(logConfig.File.Level, slog.LevelInfo)
	// 创建多输出处理器
	var handlers []slog.Handler

	// 控制台日志配置
	if logConfig.Console.Enabled {
		opts := &slog.HandlerOptions{
			Level: consoleLevel,
		}
		if logConfig.Console.Color {
			opts.AddSource = false
		}
		handlers = append(handlers, slog.NewTextHandler(os.Stdout, opts))
	}

	// 文件日志配置
	if logConfig.File.Enabled {
		// 确保日志目录存在
		logDir := filepath.Dir(logConfig.File.Path)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
		} else {
			// 配置日志轮转
			rotator := &lumberjack.Logger{
				Filename:   logConfig.File.Path,
				MaxSize:    logConfig.File.MaxSize,    // 每个日志文件的最大大小（MB）
				MaxBackups: logConfig.File.MaxBackups, // 保留的旧日志文件的最大数量
				MaxAge:     logConfig.File.MaxAge,     // 保留旧日志文件的最大天数
				Compress:   true,                      // 是否压缩旧日志文件
			}

			opts := &slog.HandlerOptions{
				Level: fileLevel,
			}
			handlers = append(handlers, slog.NewJSONHandler(rotator, opts))
		}
	}

	// 创建多处理器
	if len(handlers) > 0 {
		logger = slog.New(NewMultiHandler(handlers...))
	} else {
		// 如果没有配置任何处理器，使用默认的控制台输出
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
}

// 记录日志
func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if logger == nil {
		initLogger()
	}
	if logger != nil {
		logger.Log(ctx, level, msg, args...)
	}
}

// Infof 打印信息日志，格式化方式
func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(context.Background(), slog.LevelInfo, msg)
}

// InfofCtx 打印信息日志，格式化方式，带上下文
func InfofCtx(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(ctx, slog.LevelInfo, msg)
}

// Info 打印信息日志
func Info(msg string) {
	log(context.Background(), slog.LevelInfo, msg)
}

// InfoCtx 打印信息日志，带上下文
func InfoCtx(ctx context.Context, msg string) {
	log(ctx, slog.LevelInfo, msg)
}

// Warnf 打印警告日志，格式化方式
func Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(context.Background(), slog.LevelWarn, msg)
}

// WarnfCtx 打印警告日志，格式化方式，带上下文
func WarnfCtx(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(ctx, slog.LevelWarn, msg)
}

// Warn 打印警告日志
func Warn(msg string) {
	log(context.Background(), slog.LevelWarn, msg)
}

// WarnCtx 打印警告日志，带上下文
func WarnCtx(ctx context.Context, msg string) {
	log(ctx, slog.LevelWarn, msg)
}

// Errorf 打印错误日志，格式化方式
func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(context.Background(), slog.LevelError, msg)
}

// ErrorfCtx 打印错误日志，格式化方式，带上下文
func ErrorfCtx(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(ctx, slog.LevelError, msg)
}

// Error 打印错误日志
func Error(msg string) {
	log(context.Background(), slog.LevelError, msg)
}

// ErrorCtx 打印错误日志，带上下文
func ErrorCtx(ctx context.Context, msg string) {
	log(ctx, slog.LevelError, msg)
}

// Debugf 打印调试日志，格式化方式
func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(context.Background(), slog.LevelDebug, msg)
}

// DebugfCtx 打印调试日志，格式化方式，带上下文
func DebugfCtx(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(ctx, slog.LevelDebug, msg)
}

// Debug 打印调试日志
func Debug(msg string) {
	log(context.Background(), slog.LevelDebug, msg)
}

// DebugCtx 打印调试日志，带上下文
func DebugCtx(ctx context.Context, msg string) {
	log(ctx, slog.LevelDebug, msg)
}

// Fatalf 打印致命错误日志并终止程序，格式化方式
func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(context.Background(), slog.LevelError, msg)
	os.Exit(1)
}

// FatalfCtx 打印致命错误日志并终止程序，格式化方式，带上下文
func FatalfCtx(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	log(ctx, slog.LevelError, msg)
	os.Exit(1)
}

// Fatal 打印致命错误日志并终止程序
func Fatal(msg string) {
	log(context.Background(), slog.LevelError, msg)
	os.Exit(1)
}

// FatalCtx 打印致命错误日志并终止程序，带上下文
func FatalCtx(ctx context.Context, msg string) {
	log(ctx, slog.LevelError, msg)
	os.Exit(1)
}
