package main

import (
	"fmt"
	"os"
	"rnctl/pkg/auth"
	"rnctl/pkg/zerotier"
	"time"
)

func main() {
	fmt.Println("=== [rnctl 核心链路测试] ===")

	// 1. 加载或初始化配置 (Network ID)
	// 如果本地没文件，这里会跳出交互让你输入
	cfg, err := auth.LoadOrRequestConfig()
	if err != nil {
		fmt.Printf("[错误] 无法加载配置: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("[配置] 目标 Network ID: %s\n", cfg.NetworkID)

	// 2. 获取 ZeroTier 网络状态
	fmt.Println("[ZT] 正在获取网络列表...")
	jsonData, err := zerotier.ListNetworks()
	if err != nil {
		fmt.Printf("[错误] 调用 zerotier-cli 失败 (请检查是否安装了 sudo): %v\n", err)
		return
	}

	// 3. 解析子网 (使用你的保底方案逻辑)
	fmt.Println("[解析] 正在提取虚拟子网段...")
	subnet, err := zerotier.GetSubnet(jsonData)
	if err != nil {
		fmt.Printf("[错误] 解析网段失败: %v\n", err)
		fmt.Println("[提示] 请确保 ZeroTier 已加入网络并分配了 IP 地址。")
		return
	}
	fmt.Printf("[解析] 成功获取网段: %s\n", subnet)

	// 4. 生成待扫描 IP 列表 (使用那个递增逻辑)
	fmt.Println("[准备] 正在生成 IP 待扫列表...")
	ips, err := zerotier.GenerateIPs(subnet)
	if err != nil {
		fmt.Printf("[错误] 生成 IP 失败: %v\n", err)
		return
	}
	fmt.Printf("[准备] 已生成 %d 个目标 IP\n", len(ips))

	// 5. 启动并发扫描 (Worker Pool)
	fmt.Printf("[扫描] 开始探测端口 9527 (并发数: 100)...\n")
	s := &zerotier.Scanner{
		TargetPort: 9527,
		Timeout:    1 * time.Second, // 1秒超时
		Workers:    100,
	}

	start := time.Now()
	activeIPs := s.Scan(ips) // 这里调用的就是你刚写的那个 Scan
	duration := time.Since(start)

	// 6. 输出结果
	fmt.Println("\n=== [测试结果] ===")
	fmt.Printf("总扫描耗时: %v\n", duration)
	fmt.Printf("在线设备数量: %d\n", len(activeIPs))

	if len(activeIPs) > 0 {
		for _, ip := range activeIPs {
			fmt.Printf("[+] 发现活跃服务: %s\n", ip)
		}
	} else {
		fmt.Println("[-] 未发现开启 9527 端口的设备。")
		fmt.Println("[提示] 如果你想测试，可以在本机开一个终端运行: nc -lk 9527")
	}
}
