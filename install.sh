#!/bin/bash
# grow-check 安装脚本
# 用法: curl -sL https://raw.githubusercontent.com/openclaw-coding/grow-check/main/install.sh | bash

set -e

REPO="openclaw-coding/grow-check"
BINARY="grow-check"
INSTALL_DIR="/usr/local/bin"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🌱 Installing grow-check...${NC}"

# 检测操作系统
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo "Detected: $OS $ARCH"

# 从源码编译（暂时）
if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go 1.21+ first.${NC}"
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

echo -e "${YELLOW}Cloning repository...${NC}"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"
git clone "https://github.com/$REPO.git"
cd grow-check

echo -e "${YELLOW}Building...${NC}"
make build

echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
sudo cp ./build/grow-check "$INSTALL_DIR/grow-check"
sudo chmod +x "$INSTALL_DIR/grow-check"

# 清理
cd -
rm -rf "$TMP_DIR"

# 验证安装
if command -v grow-check &> /dev/null; then
    echo -e "${GREEN}✅ Installation successful!${NC}"
    echo ""
    echo "Quick Start:"
    echo "  1. cd your-project"
    echo "  2. grow-check init"
    echo "  3. grow-check learn --since=30d"
    echo ""
    echo "Documentation: https://github.com/$REPO#readme"
else
    echo -e "${RED}❌ Installation failed${NC}"
    exit 1
fi
