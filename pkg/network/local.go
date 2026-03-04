package network

import (
	"net"
)

// GetActiveSubnets 获取所有活跃网卡的 IPv4 CIDR 段
func GetActiveSubnets() ([]string, error) {
	var subnets []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		// 过滤：必须是开启状态，且不能是回环地址(127.0.0.1)
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := i.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				// 只处理 IPv4
				if ipNet.IP.To4() != nil {
					subnets = append(subnets, ipNet.String())
				}
			}
		}
	}
	return subnets, nil
}
