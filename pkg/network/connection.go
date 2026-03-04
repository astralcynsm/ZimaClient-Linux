package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	deviceID    string
	lanIPs      []string
	ztIPs       []string
	preferredIP string
	currentIP   string
	mu          sync.RWMutex

	// 回调函数
	onSwitching func(from, to, reason string) // 切换时的回调
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(deviceID string, lanIPs, ztIPs []string, preferredIP string) *ConnectionManager {
	return &ConnectionManager{
		deviceID:    deviceID,
		lanIPs:      lanIPs,
		ztIPs:       ztIPs,
		preferredIP: preferredIP,
		currentIP:   preferredIP,
	}
}

// SetSwitchCallback 设置切换回调
func (cm *ConnectionManager) SetSwitchCallback(callback func(from, to, reason string)) {
	cm.onSwitching = callback
}

// GetCurrentIP 获取当前使用的IP
func (cm *ConnectionManager) GetCurrentIP() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.currentIP
}

// TestConnection 测试连接是否可用
func (cm *ConnectionManager) TestConnection(ip string, port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// SmartConnect 智能连接（带重试和自动切换）
func (cm *ConnectionManager) SmartConnect(port int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. 优先尝试当前IP
	currentIP := cm.GetCurrentIP()
	if cm.TestConnection(currentIP, port, 2*time.Second) {
		return currentIP, nil
	}

	// 2. 当前IP失败，尝试LAN IP（优先级最高）
	for _, ip := range cm.lanIPs {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("连接超时")
		default:
		}

		if cm.TestConnection(ip, port, 2*time.Second) {
			cm.switchTo(ip, "LAN连接恢复")
			return ip, nil
		}
	}

	// 3. LAN全部失败，尝试ZeroTier IP
	for _, ip := range cm.ztIPs {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("连接超时")
		default:
		}

		if cm.TestConnection(ip, port, 3*time.Second) {
			cm.switchTo(ip, "通过ZeroTier连接")
			return ip, nil
		}
	}

	// 4. 所有IP都失败，进行重试（使用pkg/zerotier的重试逻辑）
	return cm.retryWithBackoff(port, ctx)
}

// retryWithBackoff 带退避的重试
func (cm *ConnectionManager) retryWithBackoff(port int, ctx context.Context) (string, error) {
	retries := []time.Duration{
		1 * time.Second,
		2 * time.Second,
		3 * time.Second,
	}

	allIPs := append(cm.lanIPs, cm.ztIPs...)

	for i, delay := range retries {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("连接超时")
		case <-time.After(delay):
		}

		// 每次重试都测试所有IP
		for _, ip := range allIPs {
			if cm.TestConnection(ip, port, 2*time.Second) {
				reason := fmt.Sprintf("第%d次重试成功", i+1)
				cm.switchTo(ip, reason)
				return ip, nil
			}
		}
	}

	return "", fmt.Errorf("所有连接尝试均失败")
}

// switchTo 切换到新IP
func (cm *ConnectionManager) switchTo(newIP, reason string) {
	cm.mu.Lock()
	oldIP := cm.currentIP
	cm.currentIP = newIP
	cm.mu.Unlock()

	// 触发回调
	if cm.onSwitching != nil && oldIP != newIP {
		cm.onSwitching(oldIP, newIP, reason)
	}
}

// HealthCheck 健康检查（定期轮询）
func (cm *ConnectionManager) HealthCheck(port int, interval time.Duration, ctx context.Context) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentIP := cm.GetCurrentIP()

			// 测试当前IP
			if !cm.TestConnection(currentIP, port, 2*time.Second) {
				// 当前IP失败，尝试切换
				_, _ = cm.SmartConnect(port)
			}
		}
	}
}
