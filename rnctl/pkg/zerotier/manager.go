package zerotier

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunCommand(args ...string) (string, error) {
	cmd := exec.Command("sudo", append([]string{"zerotier-cli"}, args...)...)
	// 将sudo以及后面的zerotier-cli以及join等下级命令append起来并且用将切片里的每一项都当成一个独立的参数传给zerotier-cli和sudo，嵌套
	var stdout bytes.Buffer // 缓冲区域
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error: %v, stderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func GetInfo() (string, error) {
	return RunCommand("status")
}

func JoinNetwork(nwid string) (string, error) {
	return RunCommand("join", nwid)
}

func ListNetworks() (string, error) {
	return RunCommand("listnetworks", "-j")
}
