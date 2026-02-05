# 1.项目概述

简述目标：构建一个轻量级的网络管理组件，支持自动发现局域网设备，通过Samba挂载与扩展存储、并且提供远程访问能力。
## 双模式运行
源码库的结构大致如下：
`pkg/`：核心逻辑库。包含ZeroTier状态管理、Worker Pool端口扫描、Samba挂载、AES加密存储。不依赖任何UI框架。

`cmd/rnctl-cli`:Cobra实现的轻量级CLI工具。
`cmd/rnctl-gui`:Wails实现的桌面GUI程序。

# 2.核心功能设计

## 网络扫描和远程连接
- 实际需求是基于`ZeroTier`的虚拟组网工具，针对`ZeroTier`给定的虚拟网卡的子网下，给定的9527端口是否有接口。核心逻辑如下：
    1. 调用本地`ZeroTier`接口指定加入的网络
    2. 轮询网络状态，等待接口分配IP地址（比如10.126.126.1/2/3等）。同时为了保证扫描效率，使用`Go`的特性`Goroutine`进行并发扫描，具体则是使用`Worker Pool`并发进行扫描该网段。
- 容错机制：目前的方案是底层`ZeroTier` + 表层的管理工具，核心问题实际上是ZeroTier的管理。针对容错，首先需要保`ZeroTier`的状态。`ZeroTier`有可能会卡在`REQUESTING_CONFIGURATION`，这种情况下程序需要能够处理这种情况。同样，`ZeroTier`有可能会因为P2P打洞失败而切换到走中继服务器，这个时候程序需要增加动态扫描的超时时间，防止因为高延迟导致的误判或者漏扫。
- 本地缓存/记忆：扫描出的设备，使用`ZeroTier`的`NODE ID`地址作为“唯一标识符”，并且在本地缓存已经发现的设备IP与Node ID的映射关系，下次启动的时候优先探测历史IP，可以提升响应速度。
- 安全性：针对Samba挂载所需要的用户名或密码，以及可能需要的ZeroTier ID等等。采用AES-256-GCM算法加密，加密密钥则通过读取Linux宿主机的`/etc/machine-id`作为种子生成加密密钥，确保数据库即便被拷贝到其他机器之后也无法解密。

## Samba挂载
- 基于Linux的`mount.cifs`工具。在这之前需要检查宿主机的`cifs-utils`依赖是否齐全，不需要的话需要先执行安装。
- 关于稳定性和容错：
    1. 默认开启`soft`和`intr`模式。当服务器断开的时候，允许系统在超时后返回错误，避免让访问该目录的进程进入D状态而导致死机。
    2. 挂载前自动检测挂载点是否存在，如果挂载失败或者挂载点不存在，程序必须负责清理产生的空目录。
    3. 如果因为网络突发中断等等问题导致的问题，一旦发现挂载失效，马上执行umount，确保本地文件系统树稳定。
    4. 程序启动的时候，解析`/proc/self/mounts`识别当前已经挂载的资源，防止重复挂载。
    5. 挂载超时、挂载权限等等问题也需要在代码中一并加入容错机制并且写好健全的反馈。

## GUI设计
采用Go下的`Wails`框架，将整个程序打包为一个跨平台可用的。将`pkg/zerotier`和`pkg/storage`的核心方法绑定到`Wails App`实例，利用Wails事件总线，将扫描进度、ZeroTier 状态变化实时推送到前端，避免前端频繁轮询。

## API/命令行设计
暂定程序名为`rnctl`(`Remote Network Control`)。命令行应用程序为`rnctl-cli`。

### 全局参数
这些参数对所有命令都有效。
- `json`：(bool)所有输出结果以JSON格式显现，方便被其它脚本或程序调用。
- `-v, --verbose`：(bool)开启详细输出。

### 子命令
#### 1. `rnctl-cli zt`
逻辑：封装宿主机的`zerotier-cli`操作，负责虚拟组网的接入和状态查询。

下级动作：
- `join <NetworkID>`：加入指定的网络。
    - 参数：`-t, --timeout`，（duration）等待入网成功的超时时间，默认30s。
- `leave <NetworkID`：离开指定网络。
- `status`：当前状态。显示`Node ID`、当前加入的网络ID、分配的IP地址、连接状态（ONLINE,OFFLINE,RELAY）。
- `scan`：探测虚拟子网下的的服务，自动获取`ZeroTier`网卡的范围，并发探测给定的端口。
  - 参数：`-t, --timeout`，（duration）端口握手超时时间。默认1s。
  - `-w, --worker`：（int）并发工作池大小。默认100.
  - `-p, --port`：（int）手动指定端口，默认为9527。

注：如果用户直接输入了`rnctl-cli zt join`而没有任何参数，程序会自动进入Survey询问：
```
? ZeroTier Network ID:
? 加入后自动扫描？（Y/n）
```


#### 2. `rnctl-cli storage`
逻辑：依赖检查 -> /proc 解析 -> 挂载执行 -> 僵尸清理。

下级动作：
- `mount`：执行挂载。
- `umount`：执行卸载。
- `list`：列出当前系统中已挂载的Samba资源。

参数：
-  `-s, --source`: (string) 远程路径，如 //10.147.x.x/Data。
-  `-t, --target`: (string) 本地挂载点，如 /mnt/remote。
-  `-u, --user`: (string) Samba 用户名。
-  `-p, --password`: (string) Samba 密码。
-  `--save`: (bool) 是否加密保存该凭据至本地。默认为`true`。
-  `--force`: (bool) 在卸载时，若目录被占用，强制执行卸载。


