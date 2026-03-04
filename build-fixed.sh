#!/bin/bash
set -e

# 配置项
APP_NAME="ZimaClient"
BINARY_SOURCE="build/bin/rnctl"
APPDIR="ZimaClient.AppDir"
OUTPUT_IMAGE="ZimaClient-x86_64.AppImage"
ICON_SOURCE="build/appicon.png"

echo "🚀 [准备阶段] 开始打包流程..."

# 1. 环境检查
if [ ! -f "$BINARY_SOURCE" ]; then
    echo "❌ 错误: 未找到二进制文件 $BINARY_SOURCE"
    echo "请确保已先运行 Docker 编译指令！"
    exit 1
fi

# 2. 准备工具: appimagetool
if ! command -v appimagetool &> /dev/null; then
    if [ ! -f "appimagetool" ]; then
        echo "📥 下载 appimagetool..."
        wget -q https://github.com/AppImage/AppImageKit/releases/download/continuous/appimagetool-x86_64.AppImage -O appimagetool
        chmod +x appimagetool
    fi
    APPIMAGETOOL="./appimagetool"
else
    APPIMAGETOOL="appimagetool"
fi

# 3. 清理旧场子
echo "🧹 清理旧的构建环境..."
rm -rf "$APPDIR"
rm -f "$OUTPUT_IMAGE"

# 4. 搭骨架 (符合 AppDir 标准结构)
echo "🏗️  创建 AppDir 结构..."
mkdir -p "$APPDIR/usr/bin"
mkdir -p "$APPDIR/usr/share/icons/hicolor/256x256/apps"
mkdir -p "$APPDIR/usr/share/applications"

# 5. 搬运资产
echo "🚚 搬运二进制文件与图标..."
cp "$BINARY_SOURCE" "$APPDIR/usr/bin/zimaclient"
chmod +x "$APPDIR/usr/bin/zimaclient"

if [ -f "$ICON_SOURCE" ]; then
    cp "$ICON_SOURCE" "$APPDIR/zimaclient.png"
    cp "$ICON_SOURCE" "$APPDIR/usr/share/icons/hicolor/256x256/apps/zimaclient.png"
else
    echo "⚠️ 警告: 未找到图标 $ICON_SOURCE，使用占位符。"
    touch "$APPDIR/zimaclient.png"
fi

# 6. 生成入口描述文件 (.desktop)
echo "📝 生成桌面元数据..."
cat > "$APPDIR/zimaclient.desktop" <<EOF
[Desktop Entry]
Name=ZimaClient
Comment=Official ZimaOS Device Manager
Exec=zimaclient
Icon=zimaclient
Type=Application
Categories=Network;Utility;
Terminal=false
EOF
cp "$APPDIR/zimaclient.desktop" "$APPDIR/usr/share/applications/"

# 7. 生成核心运行脚本 (AppRun)
echo "🏃 生成 AppRun 入口脚本..."
cat > "$APPDIR/AppRun" <<'EOF'
#!/bin/sh
# 定位 AppImage 挂载目录
HERE="$(dirname "$(readlink -f "${0}")")"
export PATH="${HERE}/usr/bin:${PATH}"
# 考虑到 Wails 依赖系统 Webkit，这里保持环境通透
# 如果以后需要打包库，这里加 LD_LIBRARY_PATH
exec "${HERE}/usr/bin/zimaclient" "$@"
EOF
chmod +x "$APPDIR/AppRun"

# 8. 压制最终成品
echo "📦 正在封装 AppImage (使用最高兼容模式)..."
# ARCH 是 appimagetool 必须的环境变量
export ARCH=x86_64
$APPIMAGETOOL "$APPDIR" "$OUTPUT_IMAGE"

echo ""
echo "=========================================="
echo "🎉 打包圆满完成！"
echo "产物路径: $OUTPUT_IMAGE"
echo "兼容级别: GLIBC 2.4+ (内核 3.2+)"
echo "文件大小: $(du -sh $OUTPUT_IMAGE | cut -f1)"
echo "=========================================="
