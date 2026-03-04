package auth

import (
	"crypto/sha256" // machine-id Hash成32位字节存放
	"fmt"
	"os"
	"strings"
)

func GetMachineKey() (string, error) {
	paths := []string{"/etc/machine-id", "/var/lib/dbus/machine-id"}
	for _, machineFile := range paths {
		if content, err := os.ReadFile(machineFile); err == nil {
			return strings.TrimSpace(string(content)), nil
		}
	}
	return "", fmt.Errorf("无法获取machine-id，请确认系统是否为Linux")
}

func GetEncryptionKey() ([]byte, error) {
	machineID, err := GetMachineKey()
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256([]byte(machineID))
	return hash[:], nil
}
