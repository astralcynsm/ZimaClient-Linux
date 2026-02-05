package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	NetworkID string `json:"network_id"`
	MachineID string `json:"machine_id"`
}

func GetConfigPath() (string, error) {
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
		fmt.Println("Cannot detect a existing file, please input your own information")
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
