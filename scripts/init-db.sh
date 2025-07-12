#!/bin/bash

# FoodCook 数据库初始化脚本

echo "🗄️  初始化 FoodCook 数据库..."

# 检查MySQL是否运行
if ! mysqladmin ping -h127.0.0.1 -uroot -proot --silent; then
    echo "❌ MySQL 未运行，请先启动 MySQL 服务"
    echo "   运行: brew services start mysql"
    exit 1
fi

# 创建数据库
echo "📦 创建数据库..."
mysql -u root -p -e "DROP DATABASE IF EXISTS foodcook; CREATE DATABASE foodcook CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入初始数据
echo "📋 导入初始数据..."
mysql -u root -p foodcook < scripts/init.sql

echo "✅ 数据库初始化完成！"
echo ""
echo "📊 数据库信息:"
echo "   数据库名: foodcook"
echo "   字符集: utf8mb4"
echo "   初始用户: root/root123"
echo "   示例数据: 8个菜品分类、18种食材、8道菜品、2条用餐记录"
echo ""
echo "🚀 现在可以运行 'make start' 启动应用" 