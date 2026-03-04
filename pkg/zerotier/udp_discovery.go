package zerotier

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

type ZimaDevice struct {
	DeviceName  string `json:"device_name"`
	DeviceModel string `json:"device_model"`
	OSVersion   string `json:"os_version"`
	Port        int    `json:"port"`
	RequestIP   string `json:"request_ip"`
	Initialized bool   `json:"initialized"`
	IP4         []struct {
		IP      string `json:"ip"`
		MAC     string `json:"mac"`
		Netcard string `json:"netcard"`
	} `json:"ip4"`
	ActualIP string `json:"actual_ip"` // 当前响应的IP

	// 多路径连接支持
	DeviceID    string   `json:"device_id"`     // 设备唯一标识（MAC地址）
	AllIPs      []string `json:"all_ips"`       // 所有可用IP
	LANIPs      []string `json:"lan_ips"`       // 局域网IP
	ZeroTierIPs []string `json:"zerotier_ips"`  // ZeroTier IP
	PreferredIP string   `json:"preferred_ip"`  // 优先使用的IP
}

func UDPSearch(boardcastAddress string, targetPort int, timeout time.Duration) ([]ZimaDevice, error) {
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP listener: %w", err)
	}
	defer conn.Close()

	dst, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", boardcastAddress, targetPort))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	probes := [][]byte{
		[]byte("ZIMA"),
		[]byte("ZIMA\n"),
	}

	// 改进了点，脉冲式发包，不然这个屎连接udp根本特么扫不到只能回退tcp那我写这个的意义何在
	go func() {
		for i := 0; i < 5; i++ { // 增加到5次
			for _, p := range probes {
				conn.WriteTo(p, dst)
			}
			time.Sleep(300 * time.Millisecond) // 减少间隔
		}
	}()

	// 去重，防止多重探针多次触发重复响应
	var foundDevice []ZimaDevice
	uniqueCheck := make(map[string]bool)
	conn.SetReadDeadline(time.Now().Add(timeout))
	buffer := make([]byte, 4096)

	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			break
		}

		var dev ZimaDevice
		if err := json.Unmarshal(buffer[:n], &dev); err != nil {
			continue
		}
		hostIP := addr.(*net.UDPAddr).IP.String()
		dev.ActualIP = hostIP

		// 填充多路径信息
		dev = enrichDeviceInfo(dev, hostIP)

		if !uniqueCheck[hostIP] {
			uniqueCheck[hostIP] = true
			foundDevice = append(foundDevice, dev)
		}
	}
	return foundDevice, nil
}

// enrichDeviceInfo 填充设备的多路径信息
func enrichDeviceInfo(dev ZimaDevice, actualIP string) ZimaDevice {
	// 使用第一个MAC地址作为设备唯一标识
	if len(dev.IP4) > 0 {
		dev.DeviceID = dev.IP4[0].MAC
	}

	// 收集所有IP
	var allIPs []string
	var lanIPs []string
	var ztIPs []string

	// 添加当前响应的IP
	allIPs = append(allIPs, actualIP)

	// 从IP4列表中提取所有IP
	for _, ipInfo := range dev.IP4 {
		if ipInfo.IP != "" && ipInfo.IP != actualIP {
			allIPs = append(allIPs, ipInfo.IP)
		}

		// 判断是LAN还是ZeroTier IP
		if isZeroTierIP(ipInfo.IP) {
			ztIPs = append(ztIPs, ipInfo.IP)
		} else if isPrivateIP(ipInfo.IP) {
			lanIPs = append(lanIPs, ipInfo.IP)
		}
	}

	// 判断当前响应IP的类型
	if isZeroTierIP(actualIP) {
		ztIPs = append(ztIPs, actualIP)
	} else if isPrivateIP(actualIP) {
		lanIPs = append(lanIPs, actualIP)
	}

	// 去重
	dev.AllIPs = uniqueStrings(allIPs)
	dev.LANIPs = uniqueStrings(lanIPs)
	dev.ZeroTierIPs = uniqueStrings(ztIPs)

	// 优先使用LAN IP
	if len(lanIPs) > 0 {
		dev.PreferredIP = lanIPs[0]
	} else if len(ztIPs) > 0 {
		dev.PreferredIP = ztIPs[0]
	} else {
		dev.PreferredIP = actualIP
	}

	return dev
}

// isZeroTierIP 判断是否为ZeroTier IP（通常在10.x.x.x或172.x.x.x范围）
func isZeroTierIP(ip string) bool {
	// ZeroTier通常使用特定的IP段
	// 这里简化判断：如果是10.147.x.x或10.209.x.x等常见ZT段
	return strings.HasPrefix(ip, "10.147.") ||
		strings.HasPrefix(ip, "10.209.") ||
		strings.HasPrefix(ip, "172.2")
}

// isPrivateIP 判断是否为私有IP
func isPrivateIP(ip string) bool {
	return strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.")
}

// uniqueStrings 字符串数组去重
func uniqueStrings(strs []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
