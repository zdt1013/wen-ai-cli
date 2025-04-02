package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/process"
)

func GetSystemInfo() (string, error) {
	info, err := host.Info()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s", info.Platform, info.PlatformVersion), nil
	// 例如 "ubuntu 22.04" 或 "windows 10.0.19045"
}

// golang获取当前运行的shell平台，例如bash、sh、zsh、powershell
func GetShellPlatform() (string, error) {
	ppid := os.Getppid() // 获取父进程ID

	// 创建父进程对象
	p, err := process.NewProcess(int32(ppid))
	if err != nil {
		return "unknown", err
	}

	// 获取父进程的可执行文件路径
	exe, err := p.Exe()
	if err != nil {
		return "unknown", err
	}

	// 提取文件名并处理
	name := filepath.Base(exe)
	name = strings.TrimSuffix(name, ".exe") // 移除.exe扩展名（Windows）
	name = strings.ToLower(name)            // 统一小写

	// 根据文件名判断Shell类型
	switch name {
	case "bash":
		return "bash", nil
	case "zsh":
		return "zsh", nil
	case "sh":
		return "sh", nil
	case "cmd":
		return "cmd", nil
	case "powershell", "pwsh":
		return "powershell", nil
	default:
		return name, nil
	}
}
