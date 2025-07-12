# 家庭菜单web应用

## 目标
用户首先点击首页可以流式查看存储在mysql数据库中的所有菜品，这里展示时需要可以方便看到该菜所属的不同分类。
展示的信息主要是菜名和图片和价格和制作链接。可以点击对应的菜跳转弹框展示该菜的详情信息，展示制作的材料，和各个材料的价格。
用户还可以勾选首页展示的若干菜品，并点击下单。点击下单后可以输入饭后感想和图片，同时点击下单后还需要展示该顿饭总价格。

首页还需要有个地方可以跳转展示历史就餐情况。可以展示出历史所有的下单信息。包含时间、以及当时下单的总价格、输入的饭后感想和图片，
并且查看时同样可以对历史下单信息点击再次下单。

注意需要有个管理页面可以录入菜品信息。

## 功能特性

### 主要功能
- **菜品管理**: 添加、编辑、删除菜品，包含图片、描述、制作链接、所属分类等信息
- **食材管理**: 记录菜品所需食材及价格信息
- **用餐记录**: 选择菜品创建用餐记录，添加感想和图片
- **历史查看**: 瀑布流展示历史用餐记录，像博客一样展示

### 用户功能
- **用户注册/登录**: 支持用户账户管理
- **菜品选择**: 在首页勾选想要吃的菜品
- **用餐记录**: 创建包含日期、感想、图片的用餐记录
- **记录展示**: 查看历史用餐记录

## 技术架构

### 后端技术栈
- **语言**: Go 1.24.1
- **Web框架**: Gin (高性能HTTP Web框架)
- **数据库**: MySQL 8.0+
- **ORM**: GORM (Go语言的ORM库)
- **认证**: JWT (JSON Web Token)
- **文件上传**: 本地存储 + 图片处理
- **配置管理**: Viper
- **日志**: Logrus
- **验证**: Go-playground/validator

### 前端技术栈
- **框架**: Vue.js 3 + TypeScript
- **UI组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router
- **HTTP客户端**: Axios
- **构建工具**: Vite

### 数据库设计
```sql
-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 菜品分类表
CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 菜品表
CREATE TABLE dishes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    price DECIMAL(10,2) NOT NULL,
    cooking_link VARCHAR(255),
    category_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- 食材表
CREATE TABLE ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 菜品食材关联表
CREATE TABLE dish_ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    dish_id BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);

-- 用餐记录表
CREATE TABLE meal_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    thoughts TEXT,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 用餐记录菜品关联表
CREATE TABLE meal_record_dishes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    meal_record_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    quantity INT DEFAULT 1,
    FOREIGN KEY (meal_record_id) REFERENCES meal_records(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);
```

## 项目结构

```
foodcook/
├── cmd/
│   └── foodcook/
│       └── main.go                 # 应用程序入口点
├── internal/
│   ├── app/
│   │   ├── handlers/               # HTTP处理器
│   │   │   ├── auth.go
│   │   │   ├── dishes.go
│   │   │   ├── ingredients.go
│   │   │   └── meal_records.go
│   │   ├── middleware/             # 中间件
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logging.go
│   │   └── routes/                 # 路由配置
│   │       └── routes.go
│   ├── domain/                     # 领域模型
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── dish.go
│   │   │   ├── ingredient.go
│   │   │   └── meal_record.go
│   │   └── repositories/           # 仓储接口
│   │       ├── user_repository.go
│   │       ├── dish_repository.go
│   │       └── meal_record_repository.go
│   └── pkg/
│       ├── database/               # 数据库连接
│       │   └── mysql.go
│       ├── config/                 # 配置管理
│       │   └── config.go
│       ├── utils/                  # 工具函数
│       │   ├── jwt.go
│       │   ├── password.go
│       │   └── file_upload.go
│       └── errors/                 # 错误处理
│           └── errors.go
├── pkg/                           # 可重用包
│   └── validator/
│       └── validator.go
├── web/                          # 前端资源
│   ├── dist/                     # 构建后的前端文件
│   ├── src/
│   │   ├── components/
│   │   ├── views/
│   │   ├── stores/
│   │   └── router/
│   ├── package.json
│   └── vite.config.ts
├── configs/                      # 配置文件
│   ├── config.yaml
│   └── config.prod.yaml
├── scripts/                      # 脚本文件
│   ├── build.sh
│   └── deploy.sh
├── docs/                         # 文档
│   ├── api.md
│   └── deployment.md
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## 安装和部署

### 环境要求
- Go 1.24.1+
- MySQL 8.0+
- Node.js 18+ (用于前端构建)

### 本地开发环境搭建

1. **克隆项目**
```bash
git clone <repository-url>
cd foodcook
```

2. **安装Go依赖**
```bash
go mod tidy
```

3. **配置数据库**
```bash
# 创建数据库
目前我数据库已经创建,数据库名是foodcook,
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=foodcook

我的redis服务器已经在本地127.0.0.1:6379启动
用户名是root没有设置密码
# 运行数据库迁移
go run cmd/foodcook/main.go migrate
```

4. **配置环境变量**
```bash
cp configs/config.yaml configs/config.local.yaml
# 编辑配置文件，设置数据库连接等信息
```

5. **启动后端服务**
```bash
go run cmd/foodcook/main.go
```

6. **启动前端开发服务器**
```bash
cd web
npm install
npm run dev
```

### Docker部署

1. **使用Docker Compose启动**
```bash
docker-compose up -d
```

2. **构建生产镜像**
```bash
make build-docker
```

### 生产环境部署

1. **构建二进制文件**
```bash
make build
```

2. **配置生产环境**
```bash
# 设置环境变量
export DB_HOST=your-db-host
export DB_PORT=3306
export DB_USER=your-db-user
export DB_PASSWORD=your-db-password
export DB_NAME=foodcook
export JWT_SECRET=your-jwt-secret
```

3. **启动服务**
```bash
./bin/foodcook
```

## API文档

### 认证相关
- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/profile` - 获取用户信息

### 菜品管理
- `GET /api/dishes` - 获取菜品列表
- `GET /api/dishes/:id` - 获取菜品详情
- `POST /api/dishes` - 创建菜品
- `PUT /api/dishes/:id` - 更新菜品
- `DELETE /api/dishes/:id` - 删除菜品

### 食材管理
- `GET /api/ingredients` - 获取食材列表
- `POST /api/ingredients` - 创建食材
- `PUT /api/ingredients/:id` - 更新食材
- `DELETE /api/ingredients/:id` - 删除食材

### 用餐记录
- `GET /api/meal-records` - 获取用餐记录列表
- `POST /api/meal-records` - 创建用餐记录
- `GET /api/meal-records/:id` - 获取用餐记录详情
- `DELETE /api/meal-records/:id` - 删除用餐记录

## 开发指南

### 代码规范
- 遵循Go官方代码规范
- 使用`gofmt`格式化代码
- 使用`golint`检查代码质量
- 编写单元测试，覆盖率不低于80%

### 提交规范
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式调整
refactor: 代码重构
test: 测试相关
chore: 构建过程或辅助工具的变动
```

### 测试
```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...

# 运行基准测试
go test -bench=. ./...
```

## 性能优化

### 数据库优化
- 使用索引优化查询性能
- 实现数据库连接池
- 使用读写分离（可选）

### 缓存策略
- 使用Redis缓存热门菜品
- 实现图片CDN加速
- 前端资源压缩和缓存

### 并发处理
- 使用goroutine处理并发请求
- 实现请求限流
- 优化文件上传处理

## 安全考虑

### 认证授权
- 使用JWT进行无状态认证
- 实现基于角色的访问控制(RBAC)
- 密码加密存储

### 数据安全
- 防止SQL注入
- 输入数据验证和清理
- 文件上传安全检查

### API安全
- 实现API限流
- 使用HTTPS
- 设置适当的CORS策略

## 监控和日志

### 日志管理
- 结构化日志记录
- 日志级别配置
- 日志轮转

### 性能监控
- 使用pprof进行性能分析
- 监控关键指标
- 错误追踪和报警

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

本项目采用MIT许可证 - 查看[LICENSE](LICENSE)文件了解详情

## 联系方式

- 项目维护者: [您的姓名]
- 邮箱: [您的邮箱]
- 项目地址: [项目GitHub地址]