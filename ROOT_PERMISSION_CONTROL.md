# Root 权限控制功能文档

## 功能概述

本次更新为 FoodCook 应用添加了基于角色的权限控制功能，只有 root 用户才能访问菜品管理和食材管理功能。

## 实现的功能

### 1. 用户角色系统

#### 后端实现
- **用户模型** (`internal/domain/models/user.go`)
  - 添加 `Role` 字段，支持 `user` 和 `root` 两种角色
  - 定义角色常量 `RoleUser` 和 `RoleRoot`

- **数据库初始化** (`internal/pkg/database/init.go`)
  - 自动创建 root 用户（用户名：root，密码：root123）
  - 设置用户角色为 `root`

#### 主要特性
- ✅ 用户模型支持角色字段
- ✅ 自动创建 root 用户
- ✅ 角色常量定义

### 2. 权限验证中间件

#### 后端实现
- **角色中间件** (`internal/app/middleware/role_middleware.go`)
  - `RootMiddleware()` 检查用户是否为 root 用户
  - 验证用户认证状态和角色权限
  - 返回相应的错误信息

#### 主要特性
- ✅ 用户认证状态验证
- ✅ 用户角色验证
- ✅ 友好的错误提示

### 3. API 权限控制

#### 后端实现
- **路由配置** (`internal/app/routes/routes.go`)
  - 菜品管理 API：只有 root 用户可以创建、更新、删除
  - 食材管理 API：只有 root 用户可以创建、更新、删除
  - 菜品查看 API：所有用户都可以访问

#### 权限分配
```
菜品管理：
- GET /api/dishes - 所有用户可查看
- GET /api/dishes/:id - 所有用户可查看
- GET /api/dishes/search - 所有用户可搜索
- POST /api/dishes - 仅 root 用户
- PUT /api/dishes/:id - 仅 root 用户
- DELETE /api/dishes/:id - 仅 root 用户

食材管理：
- GET /api/ingredients - 所有用户可查看
- POST /api/ingredients - 仅 root 用户
- PUT /api/ingredients/:id - 仅 root 用户
- DELETE /api/ingredients/:id - 仅 root 用户
```

### 4. 前端权限控制

#### 前端实现
- **认证 Store** (`web/src/stores/auth.js`)
  - 添加 `isRoot` 计算属性
  - 根据用户角色判断权限

- **路由守卫** (`web/src/router/index.js`)
  - 添加 `requiresRoot` 元信息
  - 路由级别的权限验证

- **页面权限检查**
  - 菜品管理页面 (`web/src/views/Dishes.vue`)
  - 食材管理页面 (`web/src/views/Ingredients.vue`)
  - 非 root 用户显示权限不足提示

- **导航菜单** (`web/src/views/Home.vue`)
  - 根据用户角色显示/隐藏管理菜单

#### 主要特性
- ✅ 前端路由权限验证
- ✅ 页面级权限检查
- ✅ 友好的权限不足提示
- ✅ 动态菜单显示

## 用户体验

### 1. Root 用户
- 可以访问所有功能
- 可以管理菜品和食材
- 可以看到完整的管理菜单

### 2. 普通用户
- 可以查看菜品和食材列表
- 可以创建用餐记录
- 无法访问管理功能
- 看到权限不足提示

### 3. 未登录用户
- 只能查看首页
- 需要登录才能使用其他功能

## 技术实现细节

### 1. 后端权限验证
```go
// 角色中间件
func RootMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 验证用户认证
        // 验证用户角色
        // 检查是否为 root 用户
    }
}

// 路由配置
dishes.POST("", middleware.AuthMiddleware(), middleware.RootMiddleware(), dishHandler.Create)
```

### 2. 前端权限检查
```javascript
// 认证 Store
const isRoot = computed(() => user.value?.role === 'root')

// 路由守卫
if (to.meta.requiresRoot && !authStore.isRoot) {
  next('/')
  return
}

// 页面权限检查
<div v-if="!authStore.isRoot" class="permission-denied">
  <el-result title="权限不足" sub-title="只有 root 用户才能访问此功能" />
</div>
```

## 文件结构

```
internal/
├── domain/models/
│   └── user.go                    # 用户模型（新增角色字段）
├── app/middleware/
│   └── role_middleware.go         # 角色验证中间件（新增）
├── app/routes/
│   └── routes.go                  # 路由配置（新增权限验证）
└── pkg/database/
    └── init.go                    # 数据库初始化（新增 root 用户）

web/src/
├── stores/
│   └── auth.js                    # 认证 Store（新增 isRoot）
├── router/
│   └── index.js                   # 路由配置（新增权限验证）
└── views/
    ├── Home.vue                   # 首页（新增菜单权限控制）
    ├── Dishes.vue                 # 菜品管理（新增权限检查）
    └── Ingredients.vue            # 食材管理（新增权限检查）
```

## 测试建议

### 1. Root 用户测试
- [ ] 使用 root 用户登录（用户名：root，密码：root123）
- [ ] 验证可以访问菜品管理
- [ ] 验证可以访问食材管理
- [ ] 验证可以执行所有管理操作

### 2. 普通用户测试
- [ ] 使用普通用户登录
- [ ] 验证无法访问管理页面
- [ ] 验证看到权限不足提示
- [ ] 验证菜单中不显示管理选项

### 3. API 权限测试
- [ ] 测试菜品管理 API 权限
- [ ] 测试食材管理 API 权限
- [ ] 验证非 root 用户无法执行管理操作

### 4. 前端权限测试
- [ ] 测试路由守卫功能
- [ ] 测试页面权限检查
- [ ] 测试菜单显示逻辑

## 安全考虑

1. **后端验证**：所有权限验证都在后端进行，前端验证仅用于用户体验
2. **角色检查**：每次管理操作都会验证用户角色
3. **错误处理**：提供友好的错误信息，不泄露敏感信息
4. **数据库安全**：用户角色信息存储在数据库中，确保数据一致性

## 总结

本次更新成功实现了基于角色的权限控制：

1. **完善的权限系统**：支持用户和 root 两种角色
2. **多层权限验证**：后端 API、前端路由、页面组件三层验证
3. **良好的用户体验**：友好的权限不足提示和动态菜单
4. **安全性保障**：后端主导的权限验证确保系统安全

这样的设计确保了只有 root 用户才能管理菜品和食材，同时保持了系统的安全性和可用性。 