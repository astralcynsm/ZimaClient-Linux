#!/bin/bash

# rnctl sudoers 配置安装脚本

echo "正在安装 rnctl sudoers 配置..."

# 复制 sudoers 文件
sudo cp polkit/rnctl-sudoers /etc/sudoers.d/rnctl

# 设置正确的权限
sudo chmod 440 /etc/sudoers.d/rnctl

# 验证配置
if sudo visudo -c -f /etc/sudoers.d/rnctl; then
    echo "✓ sudoers 配置已成功安装"
    echo "✓ mount、umount 和 zerotier-cli 命令现在可以无密码运行"
else
    echo "✗ sudoers 配置验证失败，请检查文件内容"
    exit 1
fi
