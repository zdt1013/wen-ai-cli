package execute

import (
	"os"
	"time"
	"wen-ai-cli/logger"

	"github.com/go-cmd/cmd"
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

// 输出命令执行状态
func printCommandStatus(status cmd.Status) {
	for _, line := range status.Stdout {
		logger.Info(line)
	}
	for _, line := range status.Stderr {
		logger.Error(line)
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

	// 创建cmd实例
	command := cmd.NewCmd(shellName, shellArg, shellCode)

	// 配置超时
	if options.Timeout > 0 {
		go func() {
			<-time.After(options.Timeout)
			command.Stop()
		}()
	}

	// 启动命令
	statusChan := command.Start()

	// 如果需要显示输出，则启动输出刷新协程
	var ticker *time.Ticker
	if options.ShowOutput {
		ticker = time.NewTicker(options.RefreshRate)
		go func() {
			for range ticker.C {
				status := command.Status()
				printCommandStatus(status)
			}
		}()
	}

	// 等待命令执行完成
	finalStatus := <-statusChan

	// 停止输出刷新
	if ticker != nil {
		ticker.Stop()
	}

	// 打印命令执行结果
	if finalStatus.Error != nil {
		logger.Errorf("命令执行出错: %v", finalStatus.Error)
		return finalStatus.Exit, finalStatus.Error
	}

	logger.Debugf("命令执行完成，退出码: %d", finalStatus.Exit)

	// 确保输出所有内容
	if options.ShowOutput {
		printCommandStatus(finalStatus)
	}

	return finalStatus.Exit, nil
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
