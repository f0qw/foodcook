# 购物车功能实现文档

## 功能概述

本次更新为 FoodCook 应用添加了购物车功能，让用户可以更方便地创建用餐记录。主要改进包括：

1. **购物车功能**：用户可以在首页将菜品添加到购物车
2. **批量创建用餐记录**：从购物车中批量创建用餐记录
3. **用餐记录编辑**：支持编辑用餐感想和图片链接

## 实现的功能

### 1. 购物车功能

#### 前端实现
- **购物车 Store** (`web/src/stores/cart.js`)
  - 管理购物车中的菜品
  - 提供添加、移除、清空购物车功能
  - 计算总价格和菜品数量

- **首页购物车界面** (`web/src/views/Home.vue`)
  - 每个菜品卡片添加"加入购物车"按钮
  - 购物车图标显示菜品数量
  - 购物车对话框显示已选菜品
  - 支持从购物车直接创建用餐记录

#### 主要特性
- ✅ 点击菜品卡片查看详情
- ✅ 点击"加入购物车"按钮添加菜品
- ✅ 购物车图标显示菜品数量徽章
- ✅ 购物车对话框显示菜品列表和总价
- ✅ 支持从购物车移除菜品
- ✅ 从购物车创建用餐记录

### 2. 用餐记录编辑功能

#### 后端实现
- **API 接口** (`internal/app/routes/routes.go`)
  - 添加 `PUT /api/meal-records/:id` 路由

- **处理器** (`internal/app/handlers/meal_record.go`)
  - 添加 `Update` 方法
  - 支持更新用餐感想和图片链接
  - 权限验证（只能编辑自己的记录）

- **仓储层** (`internal/infrastructure/repositories/mysql_meal_record_repository.go`)
  - `Update` 方法已存在，支持更新用餐记录

#### 前端实现
- **API 工具** (`web/src/utils/api.js`)
  - 添加 `mealRecordsAPI.update` 方法

- **Store** (`web/src/stores/mealRecords.js`)
  - 添加 `updateMealRecord` 方法

- **用餐记录页面** (`web/src/views/MealRecords.vue`)
  - 每个记录卡片添加"编辑"按钮
  - 编辑对话框支持修改用餐感想和图片链接
  - 实时更新记录列表

#### 主要特性
- ✅ 点击"编辑"按钮打开编辑对话框
- ✅ 支持编辑用餐感想
- ✅ 支持编辑图片链接
- ✅ 权限验证（只能编辑自己的记录）
- ✅ 编辑成功后自动刷新列表

## 用户体验改进

### 1. 首页购物流程
1. 用户浏览菜品列表
2. 点击感兴趣的菜品查看详情
3. 点击"加入购物车"按钮
4. 继续浏览其他菜品或查看购物车
5. 在购物车中确认菜品后创建用餐记录

### 2. 用餐记录管理
1. 在用餐记录页面查看所有记录
2. 点击"编辑"按钮修改记录
3. 添加或修改用餐感想和图片
4. 保存更新

## 技术实现细节

### 1. 购物车状态管理
```javascript
// 购物车 Store
const cartStore = useCartStore()

// 添加菜品到购物车
cartStore.addToCart(dish)

// 从购物车移除菜品
cartStore.removeFromCart(dishId)

// 检查菜品是否在购物车中
cartStore.isInCart(dishId)

// 获取购物车总价
cartStore.totalPrice
```

### 2. 用餐记录更新
```javascript
// 更新用餐记录
const updateRecord = async () => {
  await mealRecordsStore.updateMealRecord(recordId, {
    thoughts: editForm.thoughts,
    image_url: editForm.image_url
  })
}
```

### 3. 后端 API
```go
// 更新用餐记录
func (h *MealRecordHandler) Update(c *gin.Context) {
    // 验证用户权限
    // 更新用餐感想和图片链接
    // 返回更新后的记录
}
```

## 文件结构

```
web/src/
├── stores/
│   ├── cart.js                    # 购物车状态管理
│   └── mealRecords.js             # 用餐记录状态管理（新增更新方法）
├── views/
│   ├── Home.vue                   # 首页（新增购物车功能）
│   └── MealRecords.vue            # 用餐记录页面（新增编辑功能）
└── utils/
    └── api.js                     # API 工具（新增更新方法）

internal/
├── app/
│   ├── handlers/
│   │   └── meal_record.go         # 用餐记录处理器（新增 Update 方法）
│   └── routes/
│       └── routes.go              # 路由配置（新增 PUT 路由）
└── infrastructure/repositories/
    └── mysql_meal_record_repository.go  # MySQL 仓储（Update 方法已存在）
```

## 测试建议

### 1. 购物车功能测试
- [ ] 添加菜品到购物车
- [ ] 从购物车移除菜品
- [ ] 购物车数量显示
- [ ] 从购物车创建用餐记录
- [ ] 购物车清空功能

### 2. 用餐记录编辑测试
- [ ] 编辑用餐感想
- [ ] 编辑图片链接
- [ ] 权限验证（不能编辑他人记录）
- [ ] 编辑成功后列表更新

### 3. 用户体验测试
- [ ] 界面响应流畅
- [ ] 错误提示友好
- [ ] 操作反馈及时
- [ ] 数据一致性

## 总结

本次更新显著改善了用户体验：

1. **简化了用餐记录创建流程**：用户不再需要一次性选择所有菜品，可以逐步添加到购物车
2. **增强了用餐记录管理**：支持编辑用餐感想和图片，让记录更加丰富
3. **提升了界面交互性**：购物车功能让用户操作更加直观和便捷

这些改进让 FoodCook 应用更加用户友好，符合现代 Web 应用的使用习惯。 