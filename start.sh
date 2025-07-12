#!/bin/bash

# FoodCook 快速启动脚本
# 调用 Makefile 的 start 命令

echo "🍽️  FoodCook 家庭菜单管理系统"
echo "=================================="

# 检查Makefile是否存在
if [ ! -f "Makefile" ]; then
    echo "❌ Makefile 不存在，请确保在项目根目录运行此脚本"
    exit 1
fi

# 调用Makefile的start命令
make start 