#!/bin/bash

# FoodCook 快速启动脚本
# 用于启动前后端服务进行联调测试

echo "🍽️  FoodCook 家庭菜单管理系统"
echo "=================================="

# 检查必要的服务
check_services() {
    echo "检查必要服务..."
    
    # 检查MySQL
    if ! brew services list | grep -q "mysql.*started"; then
        echo "⚠️  MySQL服务未启动，正在启动..."
        brew services start mysql
        sleep 3
    else
        echo "✅ MySQL服务已启动"
    fi
    
    # 检查Redis
    if ! brew services list | grep -q "redis.*started"; then
        echo "⚠️  Redis服务未启动，正在启动..."
        brew services start redis
        sleep 2
    else
        echo "✅ Redis服务已启动"
    fi
}

# 启动后端服务
start_backend() {
    echo "启动后端服务..."
    cd "$(dirname "$0")"
    
    # 检查Go环境
    if ! command -v go &> /dev/null; then
        echo "❌ Go环境未安装，请先安装Go"
        exit 1
    fi
    
    # 启动后端
    echo "🚀 启动后端服务 (端口: 8080)..."
    make run &
    BACKEND_PID=$!
    
    # 等待后端启动
    echo "等待后端服务启动..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/api/dishes > /dev/null 2>&1; then
            echo "✅ 后端服务启动成功"
            break
        fi
        sleep 1
    done
    
    if [ $i -eq 30 ]; then
        echo "❌ 后端服务启动失败"
        exit 1
    fi
}

# 启动前端服务
start_frontend() {
    echo "启动前端服务..."
    cd "$(dirname "$0")/web"
    
    # 检查Node.js环境
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js环境未安装，请先安装Node.js"
        exit 1
    fi
    
    # 检查依赖
    if [ ! -d "node_modules" ]; then
        echo "📦 安装前端依赖..."
        npm install
    fi
    
    # 启动前端
    echo "🚀 启动前端服务 (端口: 3000)..."
    npm run dev &
    FRONTEND_PID=$!
    
    # 等待前端启动
    echo "等待前端服务启动..."
    for i in {1..30}; do
        if curl -s http://localhost:3000 > /dev/null 2>&1; then
            echo "✅ 前端服务启动成功"
            break
        fi
        sleep 1
    done
    
    if [ $i -eq 30 ]; then
        echo "❌ 前端服务启动失败"
        exit 1
    fi
}

# 显示服务信息
show_info() {
    echo ""
    echo "🎉 服务启动完成！"
    echo "=================================="
    echo "📱 前端地址: http://localhost:3000"
    echo "🔧 后端地址: http://localhost:8080"
    echo "📊 API文档: http://localhost:8080/api/docs"
    echo ""
    echo "📋 快速测试命令:"
    echo "  测试后端API: curl http://localhost:8080/api/dishes"
    echo "  测试前端页面: open http://localhost:3000"
    echo ""
    echo "🛑 停止服务: 按 Ctrl+C"
    echo "=================================="
}

# 清理函数
cleanup() {
    echo ""
    echo "🛑 正在停止服务..."
    
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null
        echo "✅ 后端服务已停止"
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null
        echo "✅ 前端服务已停止"
    fi
    
    echo "👋 再见！"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 主函数
main() {
    check_services
    start_backend
    start_frontend
    show_info
    
    # 保持脚本运行
    wait
}

# 运行主函数
main 