package zerotier

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"rnctl/pkg/utils"
	"strings"
	"time"
)

//go:embed packages/*
var ztPackages embed.FS

func EnsureInstalled() error {
	if _, err := exec.LookPath("zerotier-cli"); err == nil {
		return nil
	}

	fmt.Println("未检测到 ZeroTier，准备安装...")

	fmt.Println("正在尝试从官方脚本在线安装 (超时: 5s)...")
	err := tryOnlineInstall(5 * time.Second)
	if err == nil {
		fmt.Println("在线安装成功！")
		return waitForIdentity()
	}

	fmt.Printf("在线安装不可用 (%v)，正在尝试本地离线安装...\n", err)
	return installLocalPackage()
}

func tryOnlineInstall(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 构造命令
	installCmd := "curl -s https://install.zerotier.com | bash"
	cmd := exec.CommandContext(ctx, "sudo", "bash", "-c", installCmd)

	return cmd.Run()
}

// installLocalPackage 解析发行版并从内存释放包进行安装
func installLocalPackage() error {
	osID := detectOS()
	var pkgName string
	var installCmd []string

	switch osID {
	case "debian", "ubuntu", "pop", "linuxmint", "kali":
		pkgName = "zerotier-one_amd64.deb"
		installCmd = []string{"dpkg", "-i"}
	case "centos", "rocky", "almalinux", "rhel", "fedora", "amzn":
		pkgName = "zerotier-one_amd64.rpm"
		installCmd = []string{"rpm", "-ivh"}
	default:
		return fmt.Errorf("不支持的发行版离线安装: %s", osID)
	}

	// 从 embed 读取二进制流
	data, err := ztPackages.ReadFile("packages/" + pkgName)
	if err != nil {
		return fmt.Errorf("内存中找不到对应的离线包: %v", err)
	}

	// 释放到临时文件
	tmpFile := filepath.Join("/tmp", pkgName)
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	fmt.Printf("正在执行本地安装 (%s)...\n", pkgName)
	_, err = utils.RunCommand(installCmd[0], append(installCmd[1:], tmpFile)...)
	if err != nil {
		return err
	}

	// 强制拉起服务
	utils.RunCommand("systemctl", "enable", "--now", "zerotier-one")
	return waitForIdentity()
}

// 读取 /etc/os-release 获取系统标识
func detectOS() string {
	b, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID=") {
			return strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
	}
	return "unknown"
}

func waitForIdentity() error {
	fmt.Print("正在初始化ZeroTier身份信息...")
	for i := 0; i < 10; i++ {
		if _, err := os.Stat("/var/lib/zerotier-one/identity.secret"); err == nil {
			fmt.Println(" Done.")
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("ZeroTier初始化超时")
}
