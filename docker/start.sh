#!/bin/sh

# 启动nginx
echo "🚀 启动 Nginx..."
nginx -g "daemon off;" &

# 等待nginx启动
sleep 2

# 启动后端应用
echo "🚀 启动后端应用..."
./foodcook &

# 等待后端启动
sleep 3

echo "✅ 所有服务已启动"
echo "🌐 前端地址: http://localhost"
echo "🔧 后端API: http://localhost/api"
echo "🔑 初始账户: root / password"

# 保持容器运行
wait 