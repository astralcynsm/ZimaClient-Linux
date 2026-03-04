package storage

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type ShareInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

func ListRemoteShares(host, user, pswd string) ([]ShareInfo, error) {
	authString := fmt.Sprintf("%s%%%s", user, pswd)

	// 创建带超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 移除 -N 参数，因为我们要使用密码认证
	// -g: 机器可读格式输出
	cmd := exec.CommandContext(ctx, "smbclient", "-L", host, "-U", authString, "-g")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("连接超时，请检查网络或主机地址")
		}
		// 返回更详细的错误信息
		return nil, fmt.Errorf("获取Samba列表失败: %v, 输出: %s", err, string(output))
	}

	var shares []ShareInfo
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// 输出格式应该是Disk|Name|Desc之类的格式
		parts := strings.Split(line, "|")
		if len(parts) >= 2 && parts[0] == "Disk" {
			share := ShareInfo{
				Name: parts[1],
				Type: "Disk",
			}
			if len(parts) >= 3 {
				share.Comment = parts[2]
			}
			shares = append(shares, share)
		}
	}
	return shares, nil
}

// ListShares 列出主机的 SMB 共享（简化版，无需认证）
func ListShares(host string) ([]ShareInfo, error) {
	// 尝试匿名访问
	return ListRemoteShares(host, "guest", "")
}
