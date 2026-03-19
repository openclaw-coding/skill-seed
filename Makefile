.PHONY: all build install test clean fmt lint help

# 变量
BINARY_NAME=grow-check
MAIN_PATH=./cmd/grow-check
BUILD_DIR=./build
INSTALL_PATH=/usr/local/bin

# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 版本信息
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

all: clean deps build

## build: 编译二进制文件
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## install: 安装到系统
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Installation complete"

## deps: 下载依赖
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## test: 运行测试
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Test coverage:"
	$(GOCMD) tool cover -func=coverage.out | grep total

## test-html: 生成 HTML 测试覆盖率报告
test-html: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	$(GOCLEAN)

## fmt: 格式化代码
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "Format complete"

## lint: 代码检查
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

## run: 本地运行（用于开发）
run:
	$(GOBUILD) -o /tmp/$(BINARY_NAME) $(MAIN_PATH)
	/tmp/$(BINARY_NAME) $(ARGS)

## docker-build: 构建 Docker 镜像（可选）
docker-build:
	docker build -t grow-check:$(VERSION) .

## help: 显示帮助信息
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
