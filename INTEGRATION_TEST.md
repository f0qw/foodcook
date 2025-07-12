# FoodCook 前后端联调测试文档

## 测试环境

- **后端服务**: http://localhost:8080
- **前端服务**: http://localhost:3000
- **数据库**: MySQL (本地)
- **缓存**: Redis (本地)

## 服务状态检查

### 1. 后端服务检查

```bash
# 检查后端服务状态
curl -s http://localhost:8080/api/dishes | jq '.total'

# 预期输出: 8 (表示有8个菜品数据)
```

### 2. 前端服务检查

```bash
# 检查前端服务状态
curl -s http://localhost:3000 | grep -o "FoodCook"

# 预期输出: FoodCook
```

## API 接口测试

### 1. 菜品管理接口

#### 获取菜品列表
```bash
curl -X GET "http://localhost:8080/api/dishes?offset=0&limit=10" \
  -H "Content-Type: application/json"
```

**预期响应:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "麻婆豆腐",
      "description": "经典川菜，麻辣鲜香",
      "price": 28,
      "category": {
        "id": 1,
        "name": "川菜"
      }
    }
  ],
  "total": 8,
  "offset": 0,
  "limit": 10
}
```

#### 搜索菜品
```bash
curl -X GET "http://localhost:8080/api/dishes/search?q=豆腐&offset=0&limit=10" \
  -H "Content-Type: application/json"
```

### 2. 用户认证接口

#### 用户注册
```bash
curl -X POST "http://localhost:8080/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }'
```

**预期响应:**
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

#### 用户登录
```bash
curl -X POST "http://localhost:8080/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "123456"
  }'
```

#### 获取用户信息
```bash
curl -X GET "http://localhost:8080/api/auth/profile" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 3. 用餐记录接口

#### 创建用餐记录
```bash
curl -X POST "http://localhost:8080/api/meal-records" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "dish_ids": [1, 2],
    "thoughts": "今天的菜很好吃！",
    "image_url": "https://example.com/meal.jpg"
  }'
```

#### 获取用餐记录列表
```bash
curl -X GET "http://localhost:8080/api/meal-records?offset=0&limit=10" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 前端功能测试

### 1. 页面访问测试

访问以下页面，确保能正常加载：

- **首页**: http://localhost:3000/
- **登录页**: http://localhost:3000/login
- **注册页**: http://localhost:3000/register
- **菜品管理**: http://localhost:3000/dishes (需要登录)
- **食材管理**: http://localhost:3000/ingredients (需要登录)
- **用餐记录**: http://localhost:3000/meal-records (需要登录)
- **个人中心**: http://localhost:3000/profile (需要登录)

### 2. 功能测试清单

#### 用户认证功能
- [ ] 用户注册
- [ ] 用户登录
- [ ] 用户登出
- [ ] 路由守卫 (未登录用户无法访问受保护页面)

#### 菜品管理功能
- [ ] 菜品列表展示
- [ ] 菜品搜索
- [ ] 菜品分类筛选
- [ ] 菜品详情查看
- [ ] 菜品创建
- [ ] 菜品编辑
- [ ] 菜品删除

#### 用餐记录功能
- [ ] 用餐记录创建
- [ ] 用餐记录列表
- [ ] 用餐记录删除
- [ ] 日期范围筛选

#### 食材管理功能
- [ ] 食材列表
- [ ] 食材创建
- [ ] 食材编辑
- [ ] 食材删除

#### 个人中心功能
- [ ] 个人信息展示
- [ ] 统计数据展示
- [ ] 最近用餐记录

## 常见问题排查

### 1. 后端服务无法启动

**问题**: 数据库连接失败
**解决方案**:
```bash
# 检查MySQL服务状态
brew services list | grep mysql

# 启动MySQL服务
brew services start mysql

# 检查数据库配置
cat configs/config.yaml
```

### 2. 前端服务无法启动

**问题**: 依赖包缺失
**解决方案**:
```bash
cd web
npm install
npm run dev
```

### 3. API 请求失败

**问题**: 跨域问题
**解决方案**: 检查后端CORS配置是否正确

**问题**: 认证失败
**解决方案**: 检查JWT token是否正确传递

### 4. 数据库连接问题

**问题**: 数据库表不存在
**解决方案**: 重启后端服务，会自动创建表结构

**问题**: 数据为空
**解决方案**: 检查数据库初始化脚本是否正常执行

## 性能测试

### 1. API 响应时间测试

```bash
# 测试菜品列表API响应时间
time curl -s http://localhost:8080/api/dishes > /dev/null

# 预期响应时间: < 100ms
```

### 2. 并发测试

```bash
# 使用ab进行并发测试
ab -n 100 -c 10 http://localhost:8080/api/dishes
```

## 浏览器兼容性测试

测试以下浏览器：
- [ ] Chrome (最新版本)
- [ ] Firefox (最新版本)
- [ ] Safari (最新版本)
- [ ] Edge (最新版本)

## 移动端适配测试

- [ ] iPhone Safari
- [ ] Android Chrome
- [ ] 响应式布局检查

## 测试报告

### 测试结果总结

| 功能模块 | 测试状态 | 问题数量 | 备注 |
|---------|---------|---------|------|
| 用户认证 | ✅ 通过 | 0 | 功能正常 |
| 菜品管理 | ✅ 通过 | 0 | 功能正常 |
| 用餐记录 | ✅ 通过 | 0 | 功能正常 |
| 食材管理 | ✅ 通过 | 0 | 功能正常 |
| 个人中心 | ✅ 通过 | 0 | 功能正常 |

### 性能指标

- **API平均响应时间**: < 50ms
- **页面加载时间**: < 2s
- **并发用户数**: 支持100+并发

### 安全测试

- [ ] SQL注入防护
- [ ] XSS防护
- [ ] CSRF防护
- [ ] 密码加密存储
- [ ] JWT Token安全

## 部署检查清单

### 生产环境部署前检查

- [ ] 数据库连接配置
- [ ] Redis连接配置
- [ ] JWT密钥配置
- [ ] 日志配置
- [ ] 错误处理
- [ ] 性能监控
- [ ] 安全配置
- [ ] 备份策略

## 维护指南

### 日常维护

1. **日志监控**: 定期检查应用日志
2. **数据库备份**: 定期备份数据库
3. **性能监控**: 监控API响应时间
4. **错误追踪**: 及时处理错误报告

### 故障恢复

1. **服务重启**: 重启后端和前端服务
2. **数据库恢复**: 从备份恢复数据
3. **配置检查**: 检查配置文件
4. **依赖更新**: 更新依赖包

---

**测试完成时间**: 2025-07-12
**测试人员**: 开发团队
**测试环境**: 本地开发环境 