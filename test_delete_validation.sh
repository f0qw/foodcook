#!/bin/bash

echo "=== 测试删除校验功能 ==="

# 等待服务启动
sleep 3

# 获取认证token
echo "获取认证token..."
TOKEN_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"root","password":"root123"}')

TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.token // empty')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ 获取认证token失败: $TOKEN_RESPONSE"
    exit 1
fi

echo "✅ 获取token成功"

echo ""
echo "1. 测试删除正在使用的食材..."
echo "首先获取一个食材ID"
INGREDIENT_RESPONSE=$(curl -s "http://localhost:8080/api/ingredients?limit=1")
INGREDIENT_ID=$(echo $INGREDIENT_RESPONSE | jq -r '.data[0].id // empty')

if [ -z "$INGREDIENT_ID" ]; then
    echo "没有找到食材，跳过食材删除测试"
else
    echo "尝试删除食材ID: $INGREDIENT_ID"
    DELETE_RESPONSE=$(curl -s -X DELETE "http://localhost:8080/api/ingredients/$INGREDIENT_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "删除响应: $DELETE_RESPONSE"
    
    if echo "$DELETE_RESPONSE" | jq -e '.error' > /dev/null; then
        echo "✅ 食材删除校验成功 - 返回错误信息"
    else
        echo "❌ 食材删除校验失败 - 应该返回错误"
    fi
fi

echo ""
echo "2. 测试删除正在使用的菜品..."
echo "首先获取一个菜品ID"
DISH_RESPONSE=$(curl -s "http://localhost:8080/api/dishes?limit=1")
DISH_ID=$(echo $DISH_RESPONSE | jq -r '.data[0].id // empty')

if [ -z "$DISH_ID" ]; then
    echo "没有找到菜品，跳过菜品删除测试"
else
    echo "尝试删除菜品ID: $DISH_ID"
    DELETE_RESPONSE=$(curl -s -X DELETE "http://localhost:8080/api/dishes/$DISH_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "删除响应: $DELETE_RESPONSE"
    
    if echo "$DELETE_RESPONSE" | jq -e '.error' > /dev/null; then
        echo "✅ 菜品删除校验成功 - 返回错误信息"
    else
        echo "❌ 菜品删除校验失败 - 应该返回错误"
    fi
fi

echo ""
echo "3. 测试删除未使用的食材..."
echo "创建一个新的测试食材"
CREATE_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/ingredients" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"测试食材","price":1.0,"unit":"个"}')

NEW_INGREDIENT_ID=$(echo $CREATE_RESPONSE | jq -r '.id // empty')

if [ -n "$NEW_INGREDIENT_ID" ]; then
    echo "创建了测试食材ID: $NEW_INGREDIENT_ID"
    echo "尝试删除这个未使用的食材"
    
    DELETE_RESPONSE=$(curl -s -X DELETE "http://localhost:8080/api/ingredients/$NEW_INGREDIENT_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "删除响应: $DELETE_RESPONSE"
    
    if echo "$DELETE_RESPONSE" | jq -e '.message' > /dev/null; then
        echo "✅ 未使用食材删除成功"
    else
        echo "❌ 未使用食材删除失败"
    fi
else
    echo "❌ 创建测试食材失败: $CREATE_RESPONSE"
fi

echo ""
echo "4. 测试删除未使用的菜品..."
echo "创建一个新的测试菜品"
CREATE_RESPONSE=$(curl -s -X POST "http://localhost:8080/api/dishes" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"测试菜品","description":"测试描述","price":10.0,"ingredients":[]}')

NEW_DISH_ID=$(echo $CREATE_RESPONSE | jq -r '.id // empty')

if [ -n "$NEW_DISH_ID" ]; then
    echo "创建了测试菜品ID: $NEW_DISH_ID"
    echo "尝试删除这个未使用的菜品"
    
    DELETE_RESPONSE=$(curl -s -X DELETE "http://localhost:8080/api/dishes/$NEW_DISH_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo "删除响应: $DELETE_RESPONSE"
    
    if echo "$DELETE_RESPONSE" | jq -e '.message' > /dev/null; then
        echo "✅ 未使用菜品删除成功"
    else
        echo "❌ 未使用菜品删除失败"
    fi
else
    echo "❌ 创建测试菜品失败: $CREATE_RESPONSE"
fi

echo ""
echo "=== 测试完成 ===" 