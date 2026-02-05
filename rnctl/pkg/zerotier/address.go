package zerotier

import (
	"net"
)

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
