# FoodCook 快速开始指南

## 🚀 5分钟快速启动

### 1. 环境准备

确保您已安装以下软件：
- Go 1.24.1+
- MySQL 8.0+
- Redis 7.0+

### 2. 数据库设置

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE foodcook CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入初始数据
mysql -u root -p foodcook < scripts/init.sql
```

### 3. 一键启动

```bash
# 使用启动脚本（推荐）
./scripts/start.sh

# 或者手动启动
make dev
```

### 4. 验证安装

应用启动后，访问以下地址：

- **应用地址**: http://localhost:8080
- **API文档**: 查看 `docs/api.md`

## 📝 测试API

### 注册用户
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

### 获取菜品列表
```bash
curl http://localhost:8080/api/dishes
```

## 🐳 Docker 快速启动

如果您更喜欢使用 Docker：

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

## 📁 项目结构

```
foodcook/
├── cmd/foodcook/main.go          # 应用入口
├── internal/
│   ├── app/                      # 应用层
│   │   ├── handlers/             # HTTP处理器
│   │   ├── middleware/           # 中间件
│   │   └── routes/               # 路由配置
│   ├── domain/                   # 领域层
│   │   ├── models/               # 数据模型
│   │   └── repositories/         # 仓储接口
│   └── pkg/                      # 内部包
│       ├── config/               # 配置管理
│       ├── database/             # 数据库连接
│       ├── utils/                # 工具函数
│       └── errors/               # 错误处理
├── configs/                      # 配置文件
├── scripts/                      # 脚本文件
├── docs/                         # 文档
└── web/                          # 前端资源（待开发）
```

## 🔧 常用命令

```bash
# 构建应用
make build

# 运行测试
make test

# 格式化代码
make fmt

# 清理构建文件
make clean

# 使用Docker构建
make build-docker
```

## 🐛 故障排除

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查MySQL服务状态
   mysql -u root -p -e "SELECT 1;"
   ```

2. **Redis连接失败**
   ```bash
   # 检查Redis服务状态
   redis-cli ping
   ```

3. **端口被占用**
   ```bash
   # 检查端口占用
   lsof -i :8080
   ```

### 日志查看

```bash
# 查看应用日志
tail -f logs/foodcook.log

# 查看Docker日志
docker-compose logs -f app
```

## 📚 下一步

1. 查看 [API文档](docs/api.md) 了解所有接口
2. 查看 [部署指南](docs/deployment.md) 了解生产部署
3. 开发前端界面（Vue.js + Element Plus）
4. 添加更多功能（文件上传、图片处理等）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

本项目采用 MIT 许可证。 