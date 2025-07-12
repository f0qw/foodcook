# FoodCook 部署指南

## 本地开发环境

### 前置要求

- Go 1.24.1+
- MySQL 8.0+
- Redis 7.0+
- Node.js 18+ (用于前端开发)

### 快速开始

1. **克隆项目**
```bash
git clone <repository-url>
cd foodcook
```

2. **安装依赖**
```bash
make deps
```

3. **配置数据库**
```bash
# 确保MySQL服务已启动
mysql -u root -p -e "CREATE DATABASE foodcook CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行初始化脚本
mysql -u root -p foodcook < scripts/init.sql
```

4. **配置环境变量**
```bash
# 设置数据库连接信息
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=root
export DB_NAME=foodcook

# 设置Redis连接信息
export REDIS_HOST=127.0.0.1
export REDIS_PORT=6379

# 设置JWT密钥
export JWT_SECRET=your-secret-key-change-in-production
```

5. **启动应用**
```bash
make dev
```

应用将在 `http://localhost:8080` 启动。

## Docker 部署

### 使用 Docker Compose

1. **启动所有服务**
```bash
docker-compose up -d
```

2. **查看日志**
```bash
docker-compose logs -f app
```

3. **停止服务**
```bash
docker-compose down
```

### 手动构建 Docker 镜像

1. **构建镜像**
```bash
make build-docker
```

2. **运行容器**
```bash
make run-docker
```

## 生产环境部署

### 使用二进制文件

1. **构建应用**
```bash
make build
```

2. **配置生产环境**
```bash
# 创建生产配置文件
cp configs/config.yaml configs/config.prod.yaml

# 编辑配置文件，设置生产环境参数
vim configs/config.prod.yaml
```

3. **启动应用**
```bash
./bin/foodcook
```

### 使用 Systemd 服务

1. **创建服务文件**
```bash
sudo vim /etc/systemd/system/foodcook.service
```

服务文件内容:
```ini
[Unit]
Description=FoodCook Web Application
After=network.target mysql.service redis.service

[Service]
Type=simple
User=foodcook
WorkingDirectory=/opt/foodcook
ExecStart=/opt/foodcook/bin/foodcook
Restart=always
RestartSec=5
Environment=DB_HOST=localhost
Environment=DB_PORT=3306
Environment=DB_USER=foodcook
Environment=DB_PASSWORD=your-password
Environment=DB_NAME=foodcook
Environment=REDIS_HOST=localhost
Environment=REDIS_PORT=6379
Environment=JWT_SECRET=your-production-secret

[Install]
WantedBy=multi-user.target
```

2. **启用并启动服务**
```bash
sudo systemctl daemon-reload
sudo systemctl enable foodcook
sudo systemctl start foodcook
```

3. **查看服务状态**
```bash
sudo systemctl status foodcook
```

## 数据库迁移

### 自动迁移

应用启动时会自动创建数据库表结构。

### 手动迁移

1. **运行初始化脚本**
```bash
mysql -u root -p foodcook < scripts/init.sql
```

2. **检查数据库状态**
```bash
mysql -u root -p foodcook -e "SHOW TABLES;"
```

## 监控和日志

### 日志配置

应用使用结构化日志，日志级别可通过配置文件调整：

```yaml
app:
  mode: "release"  # debug, release
```

### 健康检查

应用提供健康检查端点：

```bash
curl http://localhost:8080/health
```

### 性能监控

使用 pprof 进行性能分析：

```bash
# CPU 分析
go tool pprof http://localhost:8080/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap
```

## 安全配置

### 生产环境安全建议

1. **更改默认密码**
   - 修改数据库密码
   - 修改Redis密码
   - 设置强JWT密钥

2. **网络安全**
   - 使用HTTPS
   - 配置防火墙
   - 限制数据库访问

3. **应用安全**
   - 定期更新依赖
   - 启用CORS保护
   - 实施速率限制

### 环境变量安全

```bash
# 生产环境环境变量示例
export DB_HOST=your-db-host
export DB_PORT=3306
export DB_USER=foodcook_user
export DB_PASSWORD=strong-password-here
export DB_NAME=foodcook
export REDIS_HOST=your-redis-host
export REDIS_PORT=6379
export REDIS_PASSWORD=redis-password
export JWT_SECRET=very-long-random-secret-key
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务状态
   - 验证连接参数
   - 确认网络连通性

2. **Redis连接失败**
   - 检查Redis服务状态
   - 验证连接参数
   - 确认网络连通性

3. **应用启动失败**
   - 检查配置文件
   - 查看错误日志
   - 验证依赖服务

### 日志查看

```bash
# 查看应用日志
tail -f logs/foodcook.log

# 查看系统日志
sudo journalctl -u foodcook -f
```

## 备份和恢复

### 数据库备份

```bash
# 创建备份
mysqldump -u root -p foodcook > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复备份
mysql -u root -p foodcook < backup_20240101_120000.sql
```

### 应用备份

```bash
# 备份应用文件
tar -czf foodcook_backup_$(date +%Y%m%d_%H%M%S).tar.gz \
  bin/foodcook \
  configs/ \
  uploads/ \
  scripts/
``` 