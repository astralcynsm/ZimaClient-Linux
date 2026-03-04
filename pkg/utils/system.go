package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// RunCommand 使用 sudo 执行需要权限的命令
// 配合 sudoers 配置，zerotier-cli 可以无需密码运行
func RunCommand(name string, args ...string) (string, error) {
	// 优先使用 sudo
	if _, err := exec.LookPath("sudo"); err == nil {
		fullArgs := append([]string{name}, args...)
		cmd := exec.Command("sudo", fullArgs...)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("执行%s失败: %v, stderr: %s", name, err, stderr.String())
		}
		return strings.TrimSpace(stdout.String()), nil
	}

	// 回退：尝试直接运行
	cmd := exec.Command(name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行%s失败（需要权限）: %v, stderr: %s", name, err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// RunCommandWithoutPrivilege 执行不需要权限的命令
func RunCommandWithoutPrivilege(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("执行%s失败: %v, stderr: %s", name, err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// GetSystemUptime 获取系统运行时间（JSON格式）
func GetSystemUptime() (string, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", fmt.Errorf("读取/proc/uptime失败: %v", err)
	}

	fields := strings.Fields(string(data))
	if len(fields) < 1 {
		return "", fmt.Errorf("解析/proc/uptime失败")
	}

	uptimeSeconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "", fmt.Errorf("解析uptime数值失败: %v", err)
	}

	duration := time.Duration(uptimeSeconds) * time.Second
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	result := map[string]interface{}{
		"uptime_seconds": int64(uptimeSeconds),
		"uptime_text":    fmt.Sprintf("%d天 %d小时 %d分钟", days, hours, minutes),
		"days":           days,
		"hours":          hours,
		"minutes":        minutes,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// GetHomeDir 获取用户主目录
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户主目录失败: %v", err)
	}
	return home, nil
}

// GetRuntimeMountDir 获取运行时挂载目录 /run/user/$(id -u)/zimaclient
func GetRuntimeMountDir() (string, error) {
	// 获取当前用户的 UID
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("获取当前用户失败: %v", err)
	}

	// 构建运行时目录路径
	runtimeDir := filepath.Join("/run/user", currentUser.Uid, "zimaclient")

	// 确保目录存在
	if err := os.MkdirAll(runtimeDir, 0755); err != nil {
		return "", fmt.Errorf("创建运行时目录失败: %v", err)
	}

	return runtimeDir, nil
}

// SelectDirectory 打开文件夹选择对话框
func SelectDirectory(title string) (string, error) {
	// 尝试使用 zenity (GNOME)
	if _, err := exec.LookPath("zenity"); err == nil {
		cmd := exec.Command("zenity", "--file-selection", "--directory", "--title="+title)
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("zenity选择失败: %v", err)
		}
		return strings.TrimSpace(string(output)), nil
	}

	// 尝试使用 kdialog (KDE)
	if _, err := exec.LookPath("kdialog"); err == nil {
		cmd := exec.Command("kdialog", "--getexistingdirectory", ".", "--title", title)
		output, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("kdialog选择失败: %v", err)
		}
		return strings.TrimSpace(string(output)), nil
	}

	return "", fmt.Errorf("未找到文件选择对话框工具（zenity或kdialog）")
}

// OpenFileManager 打开文件管理器并定位到指定路径
func OpenFileManager(path string) error {
	// 检查路径是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("路径不存在: %s", path)
	}

	// 尝试使用 xdg-open (通用)
	if _, err := exec.LookPath("xdg-open"); err == nil {
		cmd := exec.Command("xdg-open", path)
		return cmd.Start()
	}

	// 尝试使用 nautilus (GNOME)
	if _, err := exec.LookPath("nautilus"); err == nil {
		cmd := exec.Command("nautilus", path)
		return cmd.Start()
	}

	// 尝试使用 dolphin (KDE)
	if _, err := exec.LookPath("dolphin"); err == nil {
		cmd := exec.Command("dolphin", path)
		return cmd.Start()
	}

	return fmt.Errorf("未找到文件管理器")
}
