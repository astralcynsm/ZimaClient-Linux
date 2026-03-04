#!/bin/bash
# postinstall.sh - 安装后脚本

# 安装 sudoers 配置（允许 zerotier-cli 无需密码）
if [ -f /usr/share/rnctl/rnctl-sudoers ]; then
    cp /usr/share/rnctl/rnctl-sudoers /etc/sudoers.d/rnctl
    chmod 440 /etc/sudoers.d/rnctl
    echo "已安装 sudoers 配置"
fi

echo "rnctl 安装完成！"
echo "运行 'rnctl' 启动应用"
