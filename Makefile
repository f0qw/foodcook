.PHONY: build run test clean migrate

# 构建应用
build:
	go build -o bin/foodcook cmd/foodcook/main.go

# 运行应用
run:
	go run cmd/foodcook/main.go

# 运行测试
test:
	go test ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -cover ./...

# 清理构建文件
clean:
	rm -rf bin/
	go clean

# 格式化代码
fmt:
	go fmt ./...

# 运行代码检查
lint:
	golangci-lint run

# 数据库迁移（需要实现）
migrate:
	go run cmd/foodcook/main.go migrate

# 安装依赖
deps:
	go mod tidy
	go mod download

# 开发模式运行
dev:
	go run cmd/foodcook/main.go

# 构建Docker镜像
build-docker:
	docker build -t foodcook .

# 运行Docker容器
run-docker:
	docker run -p 8080:8080 foodcook 