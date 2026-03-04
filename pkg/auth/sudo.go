package auth

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// SudoManager 管理 sudo 密码和凭证刷新
type SudoManager struct {
	password      string
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	refreshTicker *time.Ticker
	usePolkit     bool // 是否使用 polkit
}

var (
	globalSudoManager *SudoManager
	sudoOnce          sync.Once
)

// InitSudoManager 初始化 sudo 管理器
func InitSudoManager(password string) (*SudoManager, error) {
	// 检查是否有 pkexec 可用
	usePolkit := false
	if _, err := exec.LookPath("pkexec"); err == nil {
		usePolkit = true
	}

	// 如果密码为空且有polkit，强制使用polkit
	if password == "" && usePolkit {
		usePolkit = true
	} else if password == "" && !usePolkit {
		return nil, fmt.Errorf("密码不能为空且系统不支持polkit")
	} else {
		// 有密码时，验证密码是否正确
		if err := validateSudoPassword(password); err != nil {
			return nil, fmt.Errorf("密码验证失败: %v", err)
		}
		usePolkit = false // 有密码时不使用polkit
	}

	ctx, cancel := context.WithCancel(context.Background())

	sm := &SudoManager{
		password:      password,
		ctx:           ctx,
		cancel:        cancel,
		refreshTicker: time.NewTicker(5 * time.Minute), // 每 5 分钟刷新一次
		usePolkit:     usePolkit,
	}

	// 只有在不使用 polkit 时才启动后台刷新
	if !usePolkit {
		go sm.refreshLoop()
	}

	globalSudoManager = sm
	return sm, nil
}

// GetGlobalSudoManager 获取全局 sudo 管理器
func GetGlobalSudoManager() *SudoManager {
	return globalSudoManager
}

// validateSudoPassword 验证 sudo 密码是否正确
func validateSudoPassword(password string) error {
	cmd := exec.Command("sudo", "-S", "-v")
	cmd.Stdin = bytes.NewBufferString(password + "\n")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("密码错误: %s", stderr.String())
	}

	return nil
}

// refreshLoop 后台刷新 sudo 凭证
func (sm *SudoManager) refreshLoop() {
	for {
		select {
		case <-sm.ctx.Done():
			sm.refreshTicker.Stop()
			return
		case <-sm.refreshTicker.C:
			sm.refresh()
		}
	}
}

// refresh 刷新 sudo 凭证
func (sm *SudoManager) refresh() {
	sm.mu.RLock()
	password := sm.password
	sm.mu.RUnlock()

	cmd := exec.Command("sudo", "-S", "-v")
	cmd.Stdin = bytes.NewBufferString(password + "\n")
	cmd.Run() // 忽略错误，下次会重试
}

// ExecuteWithSudo 使用缓存的密码执行 sudo 命令
func (sm *SudoManager) ExecuteWithSudo(command string, args ...string) (string, error) {
	// 如果使用 polkit，使用 pkexec
	if sm.usePolkit {
		return sm.executeWithPkexec(command, args...)
	}

	// 否则使用传统 sudo
	sm.mu.RLock()
	password := sm.password
	sm.mu.RUnlock()

	// 构建完整的命令: sudo -S command args...
	fullArgs := append([]string{"-S", command}, args...)
	cmd := exec.Command("sudo", fullArgs...)
	cmd.Stdin = bytes.NewBufferString(password + "\n")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("命令执行失败: %s", stderr.String())
	}

	return stdout.String(), nil
}

// executeWithPkexec 使用 pkexec 执行命令（图形化授权）
func (sm *SudoManager) executeWithPkexec(command string, args ...string) (string, error) {
	// 构建 pkexec 命令
	fullArgs := append([]string{command}, args...)
	cmd := exec.Command("pkexec", fullArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("命令执行失败: %s", stderr.String())
	}

	return stdout.String(), nil
}

// Stop 停止 sudo 管理器
func (sm *SudoManager) Stop() {
	if sm.cancel != nil {
		sm.cancel()
	}

	// 清除密码
	sm.mu.Lock()
	sm.password = ""
	sm.mu.Unlock()
}
