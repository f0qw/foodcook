# FoodCook 家庭菜单管理系统

一个现代化的家庭菜单管理系统，帮助您管理家庭菜品、记录用餐体验、管理食材库存。

## ✨ 主要功能

- 🍽️ **菜品管理**: 添加、编辑、删除菜品，包含图片、描述、制作链接
- 🥬 **食材管理**: 记录菜品所需食材及价格信息
- 📝 **用餐记录**: 选择菜品创建用餐记录，添加感想和图片
- 🛒 **购物车**: 将菜品加入购物车，批量创建用餐记录
- 👤 **用户系统**: 用户注册/登录，角色权限管理
- 📊 **历史查看**: 瀑布流展示历史用餐记录

## 🚀 快速开始

### 环境要求
- Go 1.24+
- MySQL 8.0+
- Node.js 18+
- Redis (可选，用于缓存)

### 方式一：一键启动 (推荐)

```bash
# 克隆项目
git clone <repository-url>
cd foodcook

# 初始化数据库
make init-db

# 一键启动所有服务
./start.sh
```

### 方式二：使用Makefile

```bash
# 初始化数据库
make init-db

# 启动所有服务
make start

# 仅启动后端服务
make start-backend

# 仅启动前端服务
make start-frontend

# 检查服务状态
make status
```

### 方式三：手动启动

```bash
# 1. 启动数据库服务
brew services start mysql
brew services start redis

# 2. 初始化数据库
make init-db

# 3. 启动后端服务
make run

# 4. 新开终端，启动前端服务
cd web
npm install
npm run dev
```

### 访问应用
- **前端应用**: http://localhost:3000
- **后端API**: http://localhost:8080

### 初始账户
- **Root用户**: 用户名 `root`，密码 `root123`
- **普通用户**: 注册新账户

## 🏗️ 技术架构

### 后端
- **语言**: Go 1.24+
- **框架**: Gin
- **数据库**: MySQL + GORM
- **认证**: JWT
- **缓存**: Redis

### 前端
- **框架**: Vue 3 + TypeScript
- **UI**: Element Plus
- **状态管理**: Pinia
- **构建工具**: Vite

## 📁 项目结构

```
foodcook/
├── cmd/foodcook/main.go          # 应用入口
├── internal/                     # 后端代码
│   ├── app/                      # 应用层
│   ├── domain/                   # 领域层
│   └── infrastructure/           # 基础设施层
├── web/                          # 前端代码
├── configs/                      # 配置文件
├── scripts/                      # 脚本文件
│   ├── init.sql                  # 数据库初始化脚本
│   └── init-db.sh                # 数据库初始化工具
├── bin/                          # 构建输出目录
├── start.sh                      # 一键启动脚本
├── Makefile                      # 构建和管理命令
└── docs/                         # 文档
```

## 🔧 开发命令

```bash
# 数据库操作
make init-db          # 初始化数据库

# 服务管理
make start            # 启动所有服务
make start-backend    # 仅启动后端
make start-frontend   # 仅启动前端
make status           # 检查服务状态

# 后端开发
make run              # 启动后端服务
make build            # 构建后端
make test             # 运行测试
make fmt              # 格式化代码
make deps             # 安装依赖

# 前端开发
cd web
npm run dev           # 启动开发服务器
npm run build         # 构建生产版本
```

## 🔐 权限控制

- **Root用户**: 可以管理菜品和食材
- **普通用户**: 可以查看菜品、创建用餐记录
- **未登录用户**: 只能查看首页

## 📚 核心API

### 认证
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/profile` - 获取用户信息

### 菜品管理
- `GET /api/dishes` - 获取菜品列表
- `POST /api/dishes` - 创建菜品 (root用户)
- `PUT /api/dishes/:id` - 更新菜品 (root用户)
- `DELETE /api/dishes/:id` - 删除菜品 (root用户)

### 用餐记录
- `GET /api/meal-records` - 获取用餐记录
- `POST /api/meal-records` - 创建用餐记录
- `PUT /api/meal-records/:id` - 更新用餐记录

## 🐳 Docker 部署

```bash
# 使用 Docker Compose
docker-compose up -d

# 构建镜像
make build-docker
```

## 📖 详细文档

- [快速使用指南](QUICK_GUIDE.md) - 详细的功能使用说明
- [API文档](docs/API.md) - 完整的API接口文档
- [部署指南](docs/DEPLOYMENT.md) - 生产环境部署说明
- [项目总结](PROJECT_SUMMARY.md) - 项目概述和技术架构

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License