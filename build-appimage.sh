#!/bin/bash
set -e

echo "=== ZimaClient AppImage 打包脚本 ==="

# 检查依赖
if ! command -v appimage-builder &> /dev/null; then
    echo "错误: 未找到 appimage-builder"
    echo "请安装: pip3 install appimage-builder"
    exit 1
fi

# 清理旧的构建
echo "清理旧的构建文件..."
rm -rf AppDir
rm -f *.AppImage

# 构建应用
echo "构建 ZimaClient..."
wails build -clean -platform linux/amd64

# 创建 AppDir 结构
echo "创建 AppDir 结构..."
mkdir -p AppDir/usr/bin
mkdir -p AppDir/usr/share/applications
mkdir -p AppDir/usr/share/icons/hicolor/256x256/apps

# 复制二进制文件
echo "复制二进制文件..."
cp build/bin/rnctl AppDir/usr/bin/zimaclient
chmod +x AppDir/usr/bin/zimaclient

# 创建 desktop 文件
echo "创建 desktop 文件..."
cat > AppDir/usr/share/applications/zimaclient.desktop << 'EOF'
[Desktop Entry]
Name=ZimaClient
Comment=Zima Device Management Client
Exec=zimaclient
Icon=zimaclient
Type=Application
Categories=Network;RemoteAccess;
Terminal=false
StartupWMClass=ZimaClient
EOF

# 创建图标（使用简单的占位符，实际应该使用真实图标）
echo "创建图标..."
if [ -f "build/appicon.png" ]; then
    cp build/appicon.png AppDir/usr/share/icons/hicolor/256x256/apps/zimaclient.png
else
    # 创建一个简单的占位符图标
    convert -size 256x256 xc:#0057FF -gravity center -pointsize 72 -fill white -annotate +0+0 "ZC" AppDir/usr/share/icons/hicolor/256x256/apps/zimaclient.png 2>/dev/null || {
        echo "警告: 无法创建图标，请手动添加 AppDir/usr/share/icons/hicolor/256x256/apps/zimaclient.png"
    }
fi

# 构建 AppImage
echo "构建 AppImage..."
appimage-builder --recipe AppImageBuilder.yml --skip-test

# 重命名输出文件
if [ -f "ZimaClient-0.1.0-x86_64.AppImage" ]; then
    mv ZimaClient-0.1.0-x86_64.AppImage ZimaClient-latest-x86_64.AppImage
    echo "✓ 构建完成: ZimaClient-latest-x86_64.AppImage"
    ls -lh ZimaClient-latest-x86_64.AppImage
else
    echo "✓ 构建完成，请检查生成的 AppImage 文件"
    ls -lh *.AppImage 2>/dev/null || echo "未找到 AppImage 文件"
fi

echo ""
echo "=== 打包完成 ==="
echo "运行方式: ./ZimaClient-latest-x86_64.AppImage"
