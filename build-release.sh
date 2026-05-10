#!/bin/bash
set -e

APP_NAME="rnctl"
VERSION="1.0.0"
OUTPUT_DIR="dist"
BINARY_PATH="build/bin/rnctl"

echo "🎨 [1/4] 启动 Docker 进行老版 GLIBC 兼容性编译..."
docker run --rm -v "$(pwd):/build" -w /build \
  -e GOPROXY=https://goproxy.cn,direct \
  zimaclient-builder:ubuntu20.04 \
  wails build -platform linux/amd64

echo "📦 [2/4] 打包 AppImage..."

chmod +x ./build-fixed.sh
./build-fixed.sh
mkdir -p $OUTPUT_DIR
mv ZimaClient-x86_64.AppImage $OUTPUT_DIR/ZimaClient-${VERSION}-x86_64.AppImage

echo "🐘 [3/4] 使用 nfpm 打包 .deb 和 .rpm ..."

docker run --rm -v "$(pwd):/tmp" -w /tmp goreleaser/nfpm pkg --target /tmp/$OUTPUT_DIR/deb.deb --packager deb
docker run --rm -v "$(pwd):/tmp" -w /tmp goreleaser/nfpm pkg --target /tmp/$OUTPUT_DIR/rpm.rpm --packager rpm

mv $OUTPUT_DIR/deb.deb $OUTPUT_DIR/${APP_NAME}_${VERSION}_amd64.deb
mv $OUTPUT_DIR/rpm.rpm $OUTPUT_DIR/${APP_NAME}-${VERSION}.x86_64.rpm

echo "🗜️  [4/4] 打包 .tar.gz 绿色版..."
tar -czvf $OUTPUT_DIR/${APP_NAME}-${VERSION}-linux-x86_64.tar.gz \
    -C build/bin rnctl \
    -C ../.. README.md

echo ""
echo "==========================================="
echo "🎉 所有平台发布包已就绪！ (目录: $OUTPUT_DIR)"
ls -lh $OUTPUT_DIR
echo "==========================================="
