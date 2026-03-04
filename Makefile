.PHONY: help docker-build docker-image clean test

# 默认目标
help:
	@echo "ZimaClient 构建命令:"
	@echo ""
	@echo "  make docker-build    - 使用 Docker 构建应用"
	@echo "  make docker-image    - 只构建 Docker 镜像"
	@echo "  make local-build     - 本地构建（需要本地环境）"
	@echo "  make appimage        - 打包 AppImage"
	@echo "  make packages        - 打包 DEB/RPM"
	@echo "  make clean           - 清理构建产物"
	@echo "  make test            - 运行测试"
	@echo ""

# Docker 构建
docker-build: docker-image
	@echo "开始 Docker 构建..."
	docker run --rm \
		-v "$(PWD):/build" \
		-w /build \
		zimaclient-builder:ubuntu20.04 \
		wails build -clean -platform linux/amd64
	@echo "构建完成: build/bin/rnctl"

# 构建 Docker 镜像
docker-image:
	@echo "构建 Docker 镜像..."
	docker build -t zimaclient-builder:ubuntu20.04 -f Dockerfile.build .

# 本地构建
local-build:
	@echo "本地构建..."
	wails build -clean -platform linux/amd64

# 打包 AppImage
appimage: docker-build
	@echo "打包 AppImage..."
	./build-appimage-simple.sh

# 打包 DEB/RPM
packages: docker-build
	@echo "打包 DEB/RPM..."
	./build-packages.sh

# 清理
clean:
	@echo "清理构建产物..."
	rm -rf build/
	rm -rf frontend/node_modules/
	rm -rf frontend/dist/
	rm -f *.AppImage
	rm -f *.deb
	rm -f *.rpm
	@echo "清理完成"

# 测试
test:
	@echo "运行测试..."
	go test ./...

# 开发模式
dev:
	@echo "启动开发模式..."
	wails dev

# 检查依赖
check-deps:
	@echo "检查依赖..."
	@command -v docker >/dev/null 2>&1 || { echo "Docker 未安装"; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "Go 未安装"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "Node.js 未安装"; exit 1; }
	@echo "所有依赖已安装"
