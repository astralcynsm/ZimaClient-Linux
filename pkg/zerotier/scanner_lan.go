package zerotier

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

// UDPUnicastScan 针对单 IP 发送探针并等待回信
func (s *Scanner) udpUnicastProbe(ip string, timeout time.Duration) (*ZimaDevice, error) {
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip, s.TargetPort))
	// 开启一个随机本地端口接收数据
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// 发送暗号
	conn.WriteTo([]byte("ZIMA"), addr)
	conn.WriteTo([]byte("ZIMA\n"), addr)

	// 设置单点探测超时（因为是局域网，不需要像广播等那么久，300ms 足够）
	conn.SetReadDeadline(time.Now().Add(timeout))

	buf := make([]byte, 2048)
	n, remoteAddr, err := conn.ReadFrom(buf)
	if err != nil {
		return nil, err
	}

	var dev ZimaDevice
	if err := json.Unmarshal(buf[:n], &dev); err == nil {
		dev.ActualIP = remoteAddr.(*net.UDPAddr).IP.String()
		// 填充多路径信息
		dev = enrichDeviceInfo(dev, dev.ActualIP)
		return &dev, nil
	}
	return nil, fmt.Errorf("decode error")
}

// 统一的并发扫描器
func (s *Scanner) SmartScanLAN(ips []string) []ZimaDevice {
	var found []ZimaDevice
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 任务通道
	jobs := make(chan string, len(ips))

	// 启动 Worker Pool
	for w := 0; w < s.Workers; w++ {
		go func() {
			for ip := range jobs {
				// 1. 尝试 UDP 单播探测
				dev, err := s.udpUnicastProbe(ip, 300*time.Millisecond)
				if err == nil {
					mu.Lock()
					found = append(found, *dev)
					mu.Unlock()
				}
				wg.Done()
			}
		}()
	}

	// 派活
	for _, ip := range ips {
		wg.Add(1)
		jobs <- ip
	}
	close(jobs)
	wg.Wait()

	return found
}
