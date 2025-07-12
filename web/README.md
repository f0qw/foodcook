# FoodCook 前端项目

这是 FoodCook 家庭菜单管理系统的前端项目，基于 Vue 3 + Element Plus 构建。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vue Router 4** - 官方路由管理器
- **Pinia** - Vue 状态管理库
- **Element Plus** - Vue 3 组件库
- **Axios** - HTTP 客户端
- **Vite** - 构建工具

## 项目结构

```
web/
├── public/                 # 静态资源
├── src/
│   ├── components/         # 公共组件
│   ├── router/            # 路由配置
│   ├── stores/            # 状态管理
│   ├── utils/             # 工具函数
│   ├── views/             # 页面组件
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── index.html             # HTML 模板
├── package.json           # 依赖配置
├── vite.config.js         # Vite 配置
└── README.md              # 项目文档
```

## 功能特性

### 1. 用户认证
- 用户注册和登录
- JWT Token 认证
- 路由守卫保护
- 用户信息管理

### 2. 菜品管理
- 菜品列表展示
- 菜品搜索和筛选
- 菜品详情查看
- 菜品增删改查
- 分类管理

### 3. 食材管理
- 食材列表管理
- 食材信息编辑
- 库存管理
- 价格管理

### 4. 用餐记录
- 用餐记录创建
- 历史记录查看
- 日期范围筛选
- 用餐感想记录
- 图片上传

### 5. 个人中心
- 个人信息展示
- 统计数据展示
- 最近用餐记录
- 信息编辑

## 页面说明

### 首页 (Home.vue)
- 菜品展示网格
- 搜索功能
- 菜品详情弹窗
- 快速创建用餐记录

### 登录页面 (Login.vue)
- 用户登录表单
- 表单验证
- 错误提示

### 注册页面 (Register.vue)
- 用户注册表单
- 密码确认验证
- 邮箱格式验证

### 菜品管理 (Dishes.vue)
- 菜品列表表格
- 搜索和筛选
- 分页功能
- 增删改查操作

### 食材管理 (Ingredients.vue)
- 食材列表管理
- 单位选择
- 价格和库存管理

### 用餐记录 (MealRecords.vue)
- 卡片式记录展示
- 日期范围筛选
- 菜品选择
- 感想记录

### 个人中心 (Profile.vue)
- 用户信息展示
- 统计数据
- 最近记录
- 信息编辑

## 状态管理

### Auth Store (stores/auth.js)
管理用户认证状态：
- 用户信息
- 登录状态
- Token 管理
- 登录/注册/登出操作

### Dishes Store (stores/dishes.js)
管理菜品数据：
- 菜品列表
- 分页信息
- 搜索功能
- CRUD 操作

### Meal Records Store (stores/mealRecords.js)
管理用餐记录：
- 记录列表
- 分页信息
- 筛选功能
- 记录操作

## API 接口

### 认证接口
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/profile` - 获取用户信息

### 菜品接口
- `GET /api/dishes` - 获取菜品列表
- `GET /api/dishes/:id` - 获取菜品详情
- `POST /api/dishes` - 创建菜品
- `PUT /api/dishes/:id` - 更新菜品
- `DELETE /api/dishes/:id` - 删除菜品
- `GET /api/dishes/search` - 搜索菜品

### 食材接口
- `GET /api/ingredients` - 获取食材列表
- `POST /api/ingredients` - 创建食材
- `PUT /api/ingredients/:id` - 更新食材
- `DELETE /api/ingredients/:id` - 删除食材

### 用餐记录接口
- `GET /api/meal-records` - 获取用餐记录
- `GET /api/meal-records/:id` - 获取记录详情
- `POST /api/meal-records` - 创建用餐记录
- `DELETE /api/meal-records/:id` - 删除用餐记录

## 开发指南

### 环境要求
- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖
```bash
cd web
npm install
```

### 开发模式
```bash
npm run dev
```

### 构建生产版本
```bash
npm run build
```

### 预览构建结果
```bash
npm run preview
```

### 代码检查
```bash
npm run lint
```

## 配置说明

### Vite 配置 (vite.config.js)
- 开发服务器代理配置
- 路径别名配置
- 构建优化配置

### 路由配置 (router/index.js)
- 路由定义
- 路由守卫
- 页面标题设置

### API 配置 (utils/api.js)
- Axios 实例配置
- 请求/响应拦截器
- 错误处理
- Token 管理

## 组件使用

### Element Plus 组件
项目使用了 Element Plus 的以下组件：
- ElButton - 按钮
- ElInput - 输入框
- ElForm - 表单
- ElTable - 表格
- ElCard - 卡片
- ElDialog - 对话框
- ElPagination - 分页
- ElSelect - 选择器
- ElDatePicker - 日期选择器
- ElImage - 图片
- ElTag - 标签
- ElAvatar - 头像
- ElDescriptions - 描述列表
- ElEmpty - 空状态

### 自定义组件
可以根据需要创建自定义组件，放在 `src/components/` 目录下。

## 样式指南

### CSS 类命名
- 使用 BEM 命名规范
- 组件样式使用 scoped
- 全局样式放在 App.vue 中

### 响应式设计
- 使用 Element Plus 的栅格系统
- 移动端适配
- 断点设置

## 部署说明

### 开发环境
1. 确保后端服务运行在 `http://localhost:8080`
2. 前端开发服务器运行在 `http://localhost:5173`
3. 配置代理转发 API 请求

### 生产环境
1. 构建前端项目：`npm run build`
2. 将 `dist` 目录部署到 Web 服务器
3. 配置 Nginx 反向代理
4. 设置 API 基础路径

### Docker 部署
```bash
# 构建镜像
docker build -t foodcook-frontend .

# 运行容器
docker run -d -p 80:80 foodcook-frontend
```

## 常见问题

### 1. 跨域问题
开发环境下已配置代理，生产环境需要配置 Nginx 反向代理。

### 2. Token 过期
Token 过期会自动跳转到登录页面，重新登录即可。

### 3. 图片加载失败
检查图片链接是否正确，或使用默认占位图。

### 4. 路由刷新 404
需要配置服务器将所有路由指向 index.html。

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 基础功能实现
- 用户认证系统
- 菜品管理功能
- 用餐记录功能

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License 