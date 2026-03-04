package zerotier

import (
	"net"
)

// GetLocalNetworks 获取所有本地网络的广播地址
func GetLocalNetworks() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var broadcasts []string
	for _, iface := range interfaces {
		// 跳过回环接口和未启用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			// 只处理IPv4
			ip4 := ipNet.IP.To4()
			if ip4 == nil {
				continue
			}

			// 计算广播地址（使用4字节的IPv4地址）
			mask := ipNet.Mask
			if len(mask) != 4 {
				continue
			}

			broadcast := make(net.IP, 4)
			for i := 0; i < 4; i++ {
				broadcast[i] = ip4[i] | ^mask[i]
			}

			broadcasts = append(broadcasts, broadcast.String())
		}
	}

	if len(broadcasts) == 0 {
		return []string{"255.255.255.255"}, nil
	}

	return broadcasts, nil
}

// GetLocalSubnets 获取所有本地子网的IP范围
func GetLocalSubnets() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var subnets []string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.To4() == nil {
				continue
			}

			subnets = append(subnets, ipNet.String())
		}
	}

	return subnets, nil
}

// GenerateIPRange 生成子网内的所有IP地址
func GenerateIPRange(cidr string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		// 跳过网络地址和广播地址
		if ip[3] == 0 || ip[3] == 255 {
			continue
		}
		ips = append(ips, ip.String())
	}

	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
