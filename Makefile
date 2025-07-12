.PHONY: build run test clean dev fmt deps start start-backend start-frontend status init-db

# 构建应用
build:
	go build -o bin/foodcook cmd/foodcook/main.go

# 运行后端应用
run:
	go run cmd/foodcook/main.go

# 开发模式运行
dev:
	go run cmd/foodcook/main.go

# 初始化数据库
init-db:
	@echo "🗄️  初始化数据库..."
	@./scripts/init-db.sh

# 一键启动所有服务
start:
	@echo "🍽️  FoodCook 家庭菜单管理系统"
	@echo "=================================="
	@echo "启动所有服务..."
	@echo ""
	@echo "1. 检查服务..."
	@if ! brew services list | grep -q "mysql.*started"; then \
		echo "⚠️  启动MySQL服务..."; \
		brew services start mysql; \
		sleep 3; \
	else \
		echo "✅ MySQL服务已启动"; \
	fi
	@if ! brew services list | grep -q "redis.*started"; then \
		echo "⚠️  启动Redis服务..."; \
		brew services start redis; \
		sleep 2; \
	else \
		echo "✅ Redis服务已启动"; \
	fi
	@echo ""
	@echo "2. 启动后端服务 (端口: 8080)..."
	@make run &
	@echo "3. 启动前端服务 (端口: 3000)..."
	@cd web && npm run dev &
	@echo ""
	@echo "🎉 服务启动完成！"
	@echo "=================================="
	@echo "📱 前端地址: http://localhost:3000"
	@echo "🔧 后端地址: http://localhost:8080"
	@echo ""
	@echo "🛑 停止服务: 按 Ctrl+C"
	@echo "=================================="
	@wait

# 仅启动后端服务
start-backend:
	@echo "🚀 启动后端服务..."
	go run cmd/foodcook/main.go

# 仅启动前端服务
start-frontend:
	@echo "🚀 启动前端服务..."
	@cd web && npm run dev

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	rm -rf bin/
	go clean

# 格式化代码
fmt:
	go fmt ./...

# 安装依赖
deps:
	go mod tidy
	go mod download

# 检查服务状态
status:
	@echo "检查服务状态..."
	@echo "MySQL: $$(brew services list | grep mysql | awk '{print $$2}')"
	@echo "Redis: $$(brew services list | grep redis | awk '{print $$2}')"
	@echo "后端: $$(curl -s http://localhost:8080/api/dishes > /dev/null && echo "运行中" || echo "未运行")"
	@echo "前端: $$(curl -s http://localhost:3000 > /dev/null && echo "运行中" || echo "未运行")" 