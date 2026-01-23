package execute

import (
	"context"
	"os"
	"os/exec"
	"time"
	"wen-ai-cli/logger"

	"github.com/gookit/i18n"
)

// ExecuteOptions 脚本执行选项
type ExecuteOptions struct {
	ShowOutput  bool          // 是否显示输出
	Timeout     time.Duration // 执行超时时间
	RefreshRate time.Duration // 输出刷新频率
}

// DefaultOptions 默认执行选项
func DefaultOptions() ExecuteOptions {
	return ExecuteOptions{
		ShowOutput:  true,
		Timeout:     0, // 无超时
		RefreshRate: time.Millisecond * 300,
	}
}

// 获取系统对应的Shell
func getSystemShell() (string, string) {
	shellName := "bash"
	shellArg := "-c"
	if os.PathSeparator == '\\' { // Windows
		shellName = "powershell"
		shellArg = "-Command"
	}
	return shellName, shellArg
}

// ExecuteScriptWithOptions 使用指定选项执行shell脚本
func ExecuteScriptWithOptions(shellCode string, options ExecuteOptions) (int, error) {
	// 根据操作系统选择合适的shell
	shellName, shellArg := getSystemShell()

	// 创建命令实例
	cmdObj := exec.Command(shellName, shellArg, shellCode)

	// 如果需要显示输出，直接连接到终端（保持原文格式）
	if options.ShowOutput {
		cmdObj.Stdout = os.Stdout
		cmdObj.Stderr = os.Stderr
	}

	// 配置超时
	var ctx context.Context
	var cancel context.CancelFunc
	if options.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), options.Timeout)
		defer cancel()
		cmdObj = exec.CommandContext(ctx, shellName, shellArg, shellCode)

		// 重新设置输出（CommandContext 会创建新实例）
		if options.ShowOutput {
			cmdObj.Stdout = os.Stdout
			cmdObj.Stderr = os.Stderr
		}
	}

	// 执行命令（输出会自动流式打印到终端）
	err := cmdObj.Run()

	// 获取退出码
	exitCode := 0
	if err != nil {
		// 尝试从 exit error 获取退出码
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			// 其他错误（如命令未找到）
			logger.Errorf("命令执行出错: %v", err)
			return -1, err
		}
	}

	logger.Debugf("命令执行完成，退出码: %d", exitCode)

	return exitCode, nil
}

// ExecuteScript 使用默认选项执行shell脚本
func ExecuteScript(shellCode string) {
	logger.Debugf(i18n.Dtr("executingScript"), shellCode)
	exitCode, _ := ExecuteScriptWithOptions(shellCode, DefaultOptions())

	// 如果退出码不为0，可以记录日志等操作
	if exitCode != 0 {
		logger.Warnf("警告：脚本执行异常，退出码: %d", exitCode)
	}
}
