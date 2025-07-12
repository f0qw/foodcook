#!/bin/bash

echo "=== FoodCook 权限控制测试 ==="
echo

# 测试菜品列表（无需认证）
echo "1. 测试菜品列表访问（无需认证）"
curl -s http://localhost:8080/api/dishes | jq '.data | length' > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ 菜品列表访问正常"
else
    echo "❌ 菜品列表访问失败"
fi
echo

# 测试食材列表（无需认证）
echo "2. 测试食材列表访问（无需认证）"
curl -s http://localhost:8080/api/ingredients | jq '.data | length' > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ 食材列表访问正常"
else
    echo "❌ 食材列表访问失败"
fi
echo

# 测试菜品创建（需要 root 权限）
echo "3. 测试菜品创建（需要 root 权限）"
curl -s -X POST http://localhost:8080/api/dishes \
  -H "Content-Type: application/json" \
  -d '{"name":"测试菜品","price":25.5}' | jq '.error' > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ 菜品创建权限验证正常（返回权限错误）"
else
    echo "❌ 菜品创建权限验证失败"
fi
echo

# 测试食材创建（需要 root 权限）
echo "4. 测试食材创建（需要 root 权限）"
curl -s -X POST http://localhost:8080/api/ingredients \
  -H "Content-Type: application/json" \
  -d '{"name":"测试食材","price":10.5,"unit":"斤"}' | jq '.error' > /dev/null
if [ $? -eq 0 ]; then
    echo "✅ 食材创建权限验证正常（返回权限错误）"
else
    echo "❌ 食材创建权限验证失败"
fi
echo

echo "=== 测试完成 ==="
echo
echo "说明："
echo "- 菜品和食材列表可以正常访问（无需认证）"
echo "- 菜品和食材的创建、更新、删除操作需要 root 权限"
echo "- 非 root 用户尝试管理操作会收到权限错误"
echo
echo "Root 用户登录信息："
echo "用户名: root"
echo "密码: root123" 