package storage

import (
	"fmt"
	"os"
	"os/user"
	"rnctl/pkg/auth"
	"strings"
)

type MountOptions struct {
	RemotePath string
	LocalPath  string
	Username   string
	Password   string
	DeviceName string // 设备名称（可选）
	ShareName  string // 共享名称（可选）
}

func createCreditFile(user, pswd string) (string, error) {
	tempFile, err := os.CreateTemp("", "rnctl_credit_")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	content := fmt.Sprintf("username=%s\npassword=%s\n", user, pswd)
	if _, err = tempFile.WriteString(content); err != nil {
		return "", err
	}
	return tempFile.Name(), nil
}

func GetOwnerID() (string, string) {
	sudoUID := os.Getenv("SUDO_UID")
	sudoGID := os.Getenv("SUDO_GID")
	if sudoUID != "" && sudoGID != "" {
		return sudoUID, sudoGID
	}
	usr, err := user.Current()
	if err != nil {
		return "1000", "1000"
	}
	return usr.Uid, usr.Gid
}
func MountSamba(opt MountOptions) error {
	// 获取 sudo 管理器
	sudoMgr := auth.GetGlobalSudoManager()
	if sudoMgr == nil {
		return fmt.Errorf("sudo 管理器未初始化，请先输入密码")
	}

	if err := os.MkdirAll(opt.LocalPath, 0755); err != nil {
		return fmt.Errorf("创建挂载点失败: %v", err)
	}

	creditFile, err := createCreditFile(opt.Username, opt.Password)
	if err != nil {
		return err
	}
	defer os.Remove(creditFile)

	usrid, gid := GetOwnerID()

	// 构造 mount 命令参数
	mountArgs := []string{
		"mount",
		"-t", "cifs",
		opt.RemotePath,
		opt.LocalPath,
		"-o", fmt.Sprintf("credentials=%s,soft,intr,uid=%s,gid=%s", creditFile, usrid, gid),
	}

	// 使用 sudo 管理器执行命令
	if _, err := sudoMgr.ExecuteWithSudo(mountArgs[0], mountArgs[1:]...); err != nil {
		return fmt.Errorf("挂载失败: %v", err)
	}

	// 挂载成功后添加到文件管理器侧边栏
	// 生成显示标签：优先使用"设备名-共享名"，否则使用远程路径
	label := opt.RemotePath
	if opt.DeviceName != "" && opt.ShareName != "" {
		label = fmt.Sprintf("%s-%s", opt.DeviceName, opt.ShareName)
	} else if opt.ShareName != "" {
		label = opt.ShareName
	}

	if err := AddToSidebar(label, opt.LocalPath); err != nil {
		// 侧边栏添加失败不影响挂载，只记录错误
		fmt.Printf("警告：添加到侧边栏失败: %v\n", err)
	}

	return nil
}

func UnmountSamba(localPath string, force bool) error {
	// 获取 sudo 管理器
	sudoMgr := auth.GetGlobalSudoManager()
	if sudoMgr == nil {
		return fmt.Errorf("sudo 管理器未初始化")
	}

	umountArgs := []string{"umount"}
	if force {
		umountArgs = append(umountArgs, "-l")
	}
	umountArgs = append(umountArgs, localPath)

	// 使用 sudo 管理器执行卸载
	if _, err := sudoMgr.ExecuteWithSudo(umountArgs[0], umountArgs[1:]...); err != nil {
		return fmt.Errorf("卸载失败：%v", err)
	}

	// 从侧边栏移除
	_ = RemoveFromSidebar(localPath)

	// 删除挂载点
	_ = os.Remove(localPath)
	return nil
}

// MountInfo 挂载信息
type MountInfo struct {
	RemotePath string `json:"remotePath"`
	LocalPath  string `json:"localPath"`
	Type       string `json:"type"`
}

// ListMounts 列出所有 CIFS 挂载
func ListMounts() ([]MountInfo, error) {
	// 直接读取 /proc/mounts 文件，避免 mount 命令可能的卡顿
	content, err := os.ReadFile("/proc/mounts")
	if err != nil {
		return nil, fmt.Errorf("读取挂载信息失败: %v", err)
	}

	var mounts []MountInfo
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		// 格式: //server/share /mnt/point cifs options 0 0
		parts := strings.Fields(line)
		if len(parts) >= 3 && parts[2] == "cifs" {
			mounts = append(mounts, MountInfo{
				RemotePath: parts[0],
				LocalPath:  parts[1],
				Type:       "cifs",
			})
		}
	}

	return mounts, nil
}
