# FoodCook API 文档

## 基础信息

- 基础URL: `http://localhost:8080/api`
- 认证方式: JWT Bearer Token
- 内容类型: `application/json`

## 认证相关

### 用户注册

**POST** `/auth/register`

请求体:
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

响应:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 用户登录

**POST** `/auth/login`

请求体:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

响应:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

### 获取用户信息

**GET** `/auth/profile`

需要认证头: `Authorization: Bearer <token>`

响应:
```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "created_at": "2024-01-01T00:00:00Z"
}
```

## 菜品管理

### 获取菜品列表

**GET** `/dishes`

查询参数:
- `offset`: 偏移量 (默认: 0)
- `limit`: 限制数量 (默认: 10)
- `category_id`: 分类ID (可选)

响应:
```json
{
  "data": [
    {
      "id": 1,
      "name": "麻婆豆腐",
      "description": "经典川菜，麻辣鲜香",
      "image_url": "https://example.com/image1.jpg",
      "price": 28.00,
      "cooking_link": "https://example.com/recipe1",
      "category_id": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "offset": 0,
  "limit": 10
}
```

### 获取菜品详情

**GET** `/dishes/{id}`

响应:
```json
{
  "id": 1,
  "name": "麻婆豆腐",
  "description": "经典川菜，麻辣鲜香",
  "image_url": "https://example.com/image1.jpg",
  "price": 28.00,
  "cooking_link": "https://example.com/recipe1",
  "category_id": 1,
  "category": {
    "id": 1,
    "name": "川菜",
    "description": "四川菜系，以麻辣著称"
  },
  "ingredients": [
    {
      "id": 1,
      "name": "猪肉",
      "price": 25.00,
      "unit": "斤",
      "quantity": 0.5
    }
  ]
}
```

### 创建菜品

**POST** `/dishes`

需要认证头: `Authorization: Bearer <token>`

请求体:
```json
{
  "name": "新菜品",
  "description": "菜品描述",
  "image_url": "https://example.com/image.jpg",
  "price": 30.00,
  "cooking_link": "https://example.com/recipe",
  "category_id": 1
}
```

### 更新菜品

**PUT** `/dishes/{id}`

需要认证头: `Authorization: Bearer <token>`

请求体:
```json
{
  "name": "更新后的菜品名",
  "price": 35.00
}
```

### 删除菜品

**DELETE** `/dishes/{id}`

需要认证头: `Authorization: Bearer <token>`

### 搜索菜品

**GET** `/dishes/search`

查询参数:
- `q`: 搜索关键词
- `offset`: 偏移量 (默认: 0)
- `limit`: 限制数量 (默认: 10)

## 食材管理

### 获取食材列表

**GET** `/ingredients`

查询参数:
- `offset`: 偏移量 (默认: 0)
- `limit`: 限制数量 (默认: 10)

### 创建食材

**POST** `/ingredients`

需要认证头: `Authorization: Bearer <token>`

请求体:
```json
{
  "name": "新食材",
  "price": 10.00,
  "unit": "斤"
}
```

### 更新食材

**PUT** `/ingredients/{id}`

需要认证头: `Authorization: Bearer <token>`

### 删除食材

**DELETE** `/ingredients/{id}`

需要认证头: `Authorization: Bearer <token>`

## 用餐记录

### 获取用餐记录列表

**GET** `/meal-records`

需要认证头: `Authorization: Bearer <token>`

查询参数:
- `offset`: 偏移量 (默认: 0)
- `limit`: 限制数量 (默认: 10)

响应:
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "total_price": 65.00,
      "thoughts": "今天的菜很好吃！",
      "image_url": "https://example.com/meal1.jpg",
      "created_at": "2024-01-01T12:00:00Z",
      "dishes": [
        {
          "id": 1,
          "name": "麻婆豆腐",
          "price": 28.00,
          "quantity": 1
        }
      ]
    }
  ],
  "total": 50,
  "offset": 0,
  "limit": 10
}
```

### 创建用餐记录

**POST** `/meal-records`

需要认证头: `Authorization: Bearer <token>`

请求体:
```json
{
  "dish_ids": [1, 2, 3],
  "thoughts": "今天的菜很好吃！",
  "image_url": "https://example.com/meal1.jpg"
}
```

### 获取用餐记录详情

**GET** `/meal-records/{id}`

需要认证头: `Authorization: Bearer <token>`

### 删除用餐记录

**DELETE** `/meal-records/{id}`

需要认证头: `Authorization: Bearer <token>`

## 错误响应

所有API在发生错误时都会返回以下格式:

```json
{
  "error": "错误描述信息"
}
```

常见HTTP状态码:
- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未认证
- `403`: 无权限
- `404`: 资源不存在
- `409`: 资源冲突
- `500`: 服务器内部错误 