package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SmbCredential struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type Config struct {
	NetworkID        string          `json:"network_id"`
	MachineID        string          `json:"machine_id"`
	SambaCredentials []SmbCredential `json:"credentials"`
}

func GetConfigPath() (string, error) {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		return filepath.Join("/home", sudoUser, ".config", "rnctl", "config.json"), nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "rnctl", "config.json"), nil
}

func LoadOrRequestConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	configFile, err := os.ReadFile(path)
	if err == nil {
		var cfg Config
		if err := json.Unmarshal(configFile, &cfg); err == nil {
			return &cfg, nil
		}
	}

	if os.IsNotExist(err) {
		fmt.Println("无法检测到已有信息，请手动输入你的信息")
		var newCfg Config

		fmt.Print("ZeroTier Network ID: ")
		fmt.Scanln(&newCfg.NetworkID)

		fmt.Print("Machine ID (Can be skipped): ")
		fmt.Scanln(&newCfg.MachineID)

		if err := SaveConfig(&newCfg); err != nil {
			return nil, fmt.Errorf("Failed to save config: %v", err)
		}

		return &newCfg, nil
	}

	return nil, err
}

// 将结果保存到硬盘
func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (cre *Config) AddCredentials(host, username, rawPassword string) error {
	encrPass, err := Encrypt(rawPassword)
	if err != nil {
		return fmt.Errorf("加密失败:%v", err)
	}
	// 去重
	found := false
	for count, credit := range cre.SambaCredentials {
		if credit.Host == host {
			cre.SambaCredentials[count].Username = username
			cre.SambaCredentials[count].Password = encrPass
			found = true
			break
		}
	}
	if !found {
		newCred := SmbCredential{Host: host, Username: username, Password: encrPass}
		cre.SambaCredentials = append(cre.SambaCredentials, newCred)
	}
	return SaveConfig(cre)
}

// 获取并解密密码
func (cre *Config) GetCredentials(host string) (string, string, error) { // 返回username, password
	for _, credit := range cre.SambaCredentials {
		if credit.Host == host {
			decrPass, err := Decrypt(credit.Password)
			if err != nil {
				return "", "", fmt.Errorf("解密失败: %v", err)
			}
			return credit.Username, decrPass, nil
		}
	}
	return "", "", fmt.Errorf("没有在%s找到对应的凭证", host)
}

// SaveToVault 保存键值对到配置文件（通用方法）
func SaveToVault(key, value string) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// 读取现有配置
	var data map[string]string
	configFile, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(configFile, &data)
	} else {
		data = make(map[string]string)
	}

	// 更新值
	data[key] = value

	// 保存
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0644)
}

// LoadFromVault 从配置文件加载键值对
func LoadFromVault(key string) (string, error) {
	path, err := GetConfigPath()
	if err != nil {
		return "", err
	}

	configFile, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var data map[string]string
	if err := json.Unmarshal(configFile, &data); err != nil {
		return "", err
	}

	value, ok := data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}
