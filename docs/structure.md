rnctl/
├── main.go            # Wails 入口 (GUI Mode)
├── app.go             # Wails 逻辑绑定 (调用 pkg)
├── wails.json
├── go.mod
├── pkg/               # 【内核层：纯逻辑，无 UI 依赖】
│   ├── zerotier/
│   │   ├── manager.go # Join/Leave/Status (CLI 调用)
│   │   └── scanner.go # Worker Pool 端口扫描 (TCP 9527)
            encoding_json.go # 负责名词解析
│   ├── storage/
│   │   ├── mounter.go # mount.cifs 封装, umount 逻辑
│   │   └── check.go   # cifs-utils 环境检查
│   ├── auth/
│   │   ├── crypto.go  # AES-GCM 加解密实现
│   │   ├── machine.go # 获取 machine-id 并生成密钥
│   │   └── vault.go   # 配置文件/凭据的读写 (JSON/SQLite)
│   └── utils/
│       ├── system.go  # 执行系统命令的封装 (处理 sudo 等)
│       └── logger.go  # 统一日志处理
├── cmd/
│   └── rnctl-cli/     # 【CLI 层：轻量级命令行】
│       ├── main.go    # CLI 入口
│       └── cmd/       # Cobra 命令定义
│           ├── root.go
│           ├── zt.go      # 调用 pkg/zerotier
│           └── storage.go # 调用 pkg/storage
├── frontend/          # 【GUI 前端：Vue/React】
│   └── ... 
└── build/             # 编译产物图标等
