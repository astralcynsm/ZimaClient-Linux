package zerotier

import (
	"encoding/json"
	"fmt"
	"net"
)

type NetworkInfo struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Name   string `json:"name"`
	Routes []struct {
		Target string `json:"target"`
	} `json:"routes"`
	AssignedAddresses []string `json:"assignedAddresses"`
	// route在json输出中有几率不会返回，AssignedAddress作为保底手段
}

// routes抽取
func GetSubnet(jsonData string) (string, error) {
	var networks []NetworkInfo
	if err := json.Unmarshal([]byte(jsonData), &networks); err != nil {
		return "", fmt.Errorf("获取网段失败: %v", err)
	}
	if len(networks) == 0 {
		return "", fmt.Errorf("未加入任何网络")
	}
	netw := networks[0]
	for _, rou := range netw.Routes {
		if rou.Target != "0.0.0.0/0" && rou.Target != "::/0" {
			return rou.Target, nil
		}
	}
	// 保底方案，从AssignedAddress抽取
	if len(netw.AssignedAddresses) > 0 {
		rawAddr := netw.AssignedAddresses[0]

		_, ipNet, err := net.ParseCIDR(rawAddr)
		if err == nil {
			return ipNet.String(), nil
		}
	}

	return "", fmt.Errorf("无法获取子网信息")
}
