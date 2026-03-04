package main

import (
	"context"
	"encoding/json"
	"fmt"
	"rnctl/pkg/auth"
	"rnctl/pkg/network"
	"rnctl/pkg/storage"
	"rnctl/pkg/utils"
	"rnctl/pkg/zerotier"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// 停止 sudo 管理器
	if sudoMgr := auth.GetGlobalSudoManager(); sudoMgr != nil {
		sudoMgr.Stop()
	}
}

// ========== Sudo 密码管理 ==========

// InitSudoPassword 初始化 sudo 密码（应用启动时调用）
func (a *App) InitSudoPassword(password string) error {
	_, err := auth.InitSudoManager(password)
	return err
}

// CheckSudoInitialized 检查 sudo 是否已初始化
func (a *App) CheckSudoInitialized() bool {
	return auth.GetGlobalSudoManager() != nil
}

// ========== ZeroTier 相关 ==========

// EnsureZeroTierInstalled 确保 ZeroTier 已安装
func (a *App) EnsureZeroTierInstalled() error {
	return zerotier.EnsureInstalled()
}

// GetZeroTierInfo 获取 ZeroTier 状态信息
func (a *App) GetZeroTierInfo() (string, error) {
	return zerotier.GetInfo()
}

// JoinZeroTierNetwork 加入 ZeroTier 网络
func (a *App) JoinZeroTierNetwork(networkID string) (string, error) {
	return zerotier.JoinNetwork(networkID)
}

// LeaveZeroTierNetwork 离开 ZeroTier 网络
func (a *App) LeaveZeroTierNetwork(networkID string) (string, error) {
	return zerotier.LeaveNetwork(networkID)
}

// ListZeroTierNetworks 列出已加入的网络
func (a *App) ListZeroTierNetworks() (string, error) {
	return zerotier.ListNetworks()
}

// GetSystemUptime 获取系统运行时间
func (a *App) GetSystemUptime() (string, error) {
	return utils.GetSystemUptime()
}

// GetHomeDir 获取用户主目录
func (a *App) GetHomeDir() (string, error) {
	return utils.GetHomeDir()
}

// GetRuntimeMountDir 获取运行时挂载目录
func (a *App) GetRuntimeMountDir() (string, error) {
	return utils.GetRuntimeMountDir()
}

// SelectMountDirectory 打开文件夹选择对话框
func (a *App) SelectMountDirectory() (string, error) {
	return utils.SelectDirectory("选择挂载目录")
}

// ========== 设备扫描 ==========

// ScanLocalDevices 扫描局域网内的 Zima 设备
func (a *App) ScanLocalDevices() (string, error) {
	// 1. 获取所有活跃的网段
	subnets, err := network.GetActiveSubnets()
	if err != nil {
		return "[]", err
	}

	// 2. 汇总所有待扫描的 IP
	var allIPs []string
	for _, subnet := range subnets {
		ips, err := zerotier.GenerateIPs(subnet)
		if err == nil {
			allIPs = append(allIPs, ips...)
		}
	}

	// 3. 配置扫描器 - 200个workers并发扫描
	scanner := &zerotier.Scanner{
		TargetPort: 9527,
		Timeout:    500 * time.Millisecond,
		Workers:    200,
	}

	// 4. 执行扫描
	devices := scanner.SmartScanLAN(allIPs)

	// 5. 转换为 JSON
	result, err := json.Marshal(devices)
	if err != nil {
		return "[]", err
	}

	return string(result), nil
}

// ========== 存储管理 ==========

// MountSMB 挂载 SMB 共享
func (a *App) MountSMB(remotePath, localPath, username, password, deviceName, shareName string) error {
	return storage.MountSamba(storage.MountOptions{
		RemotePath: remotePath,
		LocalPath:  localPath,
		Username:   username,
		Password:   password,
		DeviceName: deviceName,
		ShareName:  shareName,
	})
}

// UnmountSMB 卸载 SMB 共享
func (a *App) UnmountSMB(localPath string, force bool) error {
	return storage.UnmountSamba(localPath, force)
}

// UnmountAll 卸载所有 SMB 共享
func (a *App) UnmountAll() error {
	mounts, err := storage.ListMounts()
	if err != nil {
		return err
	}

	var lastErr error
	for _, mount := range mounts {
		if err := storage.UnmountSamba(mount.LocalPath, false); err != nil {
			lastErr = err
			// 继续卸载其他挂载点
		}
	}

	return lastErr
}

// ListMounts 列出已挂载的共享
func (a *App) ListMounts() (string, error) {
	mounts, err := storage.ListMounts()
	if err != nil {
		return "[]", err
	}

	result, err := json.Marshal(mounts)
	if err != nil {
		return "[]", err
	}

	return string(result), nil
}

// ListSMBShares 列出 SMB 共享（带认证）
func (a *App) ListSMBShares(host, username, password string) (string, error) {
	shares, err := storage.ListRemoteShares(host, username, password)
	if err != nil {
		return "[]", err
	}

	result, err := json.Marshal(shares)
	if err != nil {
		return "[]", err
	}

	return string(result), nil
}

// ========== 认证管理 ==========

// SaveCredentials 保存加密的凭证
func (a *App) SaveCredentials(key, value string) error {
	encrypted, err := auth.Encrypt(value)
	if err != nil {
		return err
	}
	return auth.SaveToVault(key, encrypted)
}

// GetCredentials 获取并解密凭证
func (a *App) GetCredentials(key string) (string, error) {
	encrypted, err := auth.LoadFromVault(key)
	if err != nil {
		return "", err
	}
	return auth.Decrypt(encrypted)
}

// ========== 自动挂载配置 ==========

// SaveAutoMountConfig 保存自动挂载配置
func (a *App) SaveAutoMountConfig(config string) error {
	return a.SaveCredentials("auto_mount_config", config)
}

// GetAutoMountConfig 获取自动挂载配置
func (a *App) GetAutoMountConfig() (string, error) {
	return a.GetCredentials("auto_mount_config")
}

// ProcessAutoMount 处理自动挂载（应用启动时调用）
func (a *App) ProcessAutoMount() error {
	configStr, err := a.GetAutoMountConfig()
	if err != nil {
		return err
	}

	var config struct {
		Enabled   bool   `json:"enabled"`
		Configured bool  `json:"configured"`
		DeviceIP  string `json:"deviceIp"`
		Mounts    []struct {
			Name       string `json:"name"`
			RemotePath string `json:"remotePath"`
			LocalPath  string `json:"localPath"`
		} `json:"mounts"`
		QuickAccess *struct {
			Name      string `json:"name"`
			LocalPath string `json:"localPath"`
		} `json:"quickAccess"`
	}

	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return err
	}

	// 如果未启用或未配置，直接返回
	if !config.Enabled || !config.Configured {
		return nil
	}

	// 获取保存的凭证
	credStr, err := a.GetCredentials(config.DeviceIP)
	if err != nil {
		return err
	}

	var cred struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal([]byte(credStr), &cred); err != nil {
		return err
	}

	// 执行挂载
	for _, mount := range config.Mounts {
		// 检查是否已挂载
		mounts, _ := storage.ListMounts()
		alreadyMounted := false
		for _, m := range mounts {
			if m.LocalPath == mount.LocalPath {
				alreadyMounted = true
				break
			}
		}

		if !alreadyMounted {
			// 执行挂载
			err := storage.MountSamba(storage.MountOptions{
				RemotePath: mount.RemotePath,
				LocalPath:  mount.LocalPath,
				Username:   cred.Username,
				Password:   cred.Password,
				DeviceName: config.DeviceIP,
				ShareName:  mount.Name,
			})
			if err != nil {
				// 记录错误但继续挂载其他共享
				continue
			}
		}
	}

	return nil
}

// OpenQuickAccess 打开快速访问文件夹
func (a *App) OpenQuickAccess() error {
	configStr, err := a.GetAutoMountConfig()
	if err != nil {
		return err
	}

	var config struct {
		QuickAccess *struct {
			Name      string `json:"name"`
			LocalPath string `json:"localPath"`
		} `json:"quickAccess"`
	}

	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return err
	}

	if config.QuickAccess == nil {
		return fmt.Errorf("未配置快速访问")
	}

	return utils.OpenFileManager(config.QuickAccess.LocalPath)
}

// ========== 智能连接管理 ==========

// TestDeviceConnection 测试设备连接（智能重试）
func (a *App) TestDeviceConnection(deviceJSON string) (string, error) {
	var device struct {
		DeviceID    string   `json:"device_id"`
		LANIPs      []string `json:"lan_ips"`
		ZeroTierIPs []string `json:"zerotier_ips"`
		PreferredIP string   `json:"preferred_ip"`
	}

	if err := json.Unmarshal([]byte(deviceJSON), &device); err != nil {
		return "", fmt.Errorf("解析设备信息失败: %v", err)
	}

	// 创建连接管理器
	cm := network.NewConnectionManager(
		device.DeviceID,
		device.LANIPs,
		device.ZeroTierIPs,
		device.PreferredIP,
	)

	// 尝试智能连接
	connectedIP, err := cm.SmartConnect(9527)
	if err != nil {
		return "", err
	}

	// 返回连接结果
	result := map[string]interface{}{
		"success":      true,
		"connected_ip": connectedIP,
		"message":      "连接成功",
	}

	resultJSON, _ := json.Marshal(result)
	return string(resultJSON), nil
}
