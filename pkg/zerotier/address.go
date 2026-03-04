package zerotier

import (
	"fmt"
	"net"
)

func GetBoardcastAddr(subnet string) (string, error) {
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return "", fmt.Errorf("解析失败: %v", err)
	}
	ip := ipNet.IP.To4()
	if ip == nil {
		return "", fmt.Errorf("仅支持IPV4广播")
	}

	mask := ipNet.Mask
	boardcast := make(net.IP, len(ip))
	for i := range ip {
		boardcast[i] = ip[i] | ^mask[i]
	}
	return boardcast.String(), nil
}

func GenerateIPs(cidr string) ([]string, error) {
	defIP, ipNetwork, err := net.ParseCIDR(cidr)

	if err != nil {
		return nil, err
	}

	var ipList []string

	curIP := make(net.IP, len(defIP))
	copy(curIP, defIP)

	for ipNetwork.Contains(curIP) {
		ipList = append(ipList, curIP.String())
		// IP递增逻辑，万一分配的不是/24能保证健全性
		for i := len(curIP) - 1; i >= 0; i-- {
			curIP[i]++
			if curIP[i] > 0 {
				break
			}
		}
	}
	if len(ipList) > 2 {
		return ipList[1 : len(ipList)-1], nil
	}
	return ipList, nil
}
