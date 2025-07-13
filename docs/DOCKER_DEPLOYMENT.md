# Docker 部署指南

## 概述

FoodCook 支持使用 Docker 进行容器化部署，提供完整的开发和生产环境支持。

## 架构说明

### 容器架构
- **app**: 主应用容器（包含前端和后端）
- **mysql**: MySQL 数据库
- **redis**: Redis 缓存
- **phpMyAdmin**: 数据库管理工具（可选）

### 网络架构
- 前端通过 Nginx 代理到后端 API
- 所有服务通过 Docker 网络通信
- 支持健康检查和自动重启

## 快速开始

### 1. 环境要求

- Docker 20.10+
- Docker Compose 2.0+
- 至少 2GB 可用内存
- 至少 5GB 可用磁盘空间

### 2. 一键启动

```bash
# 克隆项目
git clone <repository-url>
cd foodcook

# 一键启动
./docker-start.sh
```

### 3. 手动启动

```bash
# 构建并启动服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

## 开发环境部署

### 配置文件

开发环境使用 `docker-compose.yml`，包含以下特性：

- 前端和后端集成部署
- 数据库自动初始化
- 可选的管理工具
- 健康检查
- 日志收集

### 启动命令

```bash
# 启动所有服务
docker-compose up -d

# 启动特定服务
docker-compose up -d app mysql redis

# 启动包含管理工具
docker-compose --profile tools up -d
```

### 访问地址

- **前端应用**: http://localhost
- **后端API**: http://localhost/api
- **数据库**: localhost:3306
- **Redis**: localhost:6379
- **phpMyAdmin**: http://localhost:8081

## 生产环境部署

### 1. 环境变量配置

```bash
# 复制环境变量模板
cp env.example .env

# 编辑环境变量
vim .env
```

### 2. 生产环境配置

生产环境使用 `docker-compose.prod.yml`，包含以下特性：

- 资源限制
- 安全配置
- 环境变量支持
- 健康检查
- 自动重启

### 3. 启动生产环境

```bash
# 启动生产环境
docker-compose -f docker-compose.prod.yml up -d --build

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps
```

### 4. 环境变量说明

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| DB_USER | 数据库用户名 | root |
| DB_PASSWORD | 数据库密码 | root |
| DB_NAME | 数据库名 | foodcook |
| REDIS_PASSWORD | Redis密码 | 空 |
| JWT_SECRET | JWT密钥 | your-secret-key-change-in-production |
| APP_MODE | 应用模式 | release |

## 服务管理

### 常用命令

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f app      # 应用日志
docker-compose logs -f mysql    # 数据库日志
docker-compose logs -f redis    # Redis日志

# 进入容器
docker-compose exec app sh      # 进入应用容器
docker-compose exec mysql mysql -u root -proot foodcook  # 进入数据库

# 清理数据
docker-compose down -v          # 停止并删除数据卷
```

### 健康检查

所有服务都配置了健康检查：

```bash
# 查看健康状态
docker-compose ps

# 手动检查健康状态
docker-compose exec app wget -q -O- http://localhost/health
docker-compose exec mysql mysqladmin ping -h localhost -u root -proot
docker-compose exec redis redis-cli ping
```

## 数据持久化

### 数据卷

- **mysql_data**: MySQL 数据持久化
- **redis_data**: Redis 数据持久化
- **app_logs**: 应用日志持久化
- **uploads**: 上传文件持久化

### 备份和恢复

```bash
# 备份数据库
docker-compose exec mysql mysqldump -u root -proot foodcook > backup.sql

# 恢复数据库
docker-compose exec -T mysql mysql -u root -proot foodcook < backup.sql

# 备份上传文件
tar -czf uploads_backup.tar.gz uploads/

# 恢复上传文件
tar -xzf uploads_backup.tar.gz
```

## 监控和日志

### 日志管理

```bash
# 查看实时日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app

# 查看错误日志
docker-compose logs --tail=100 app | grep ERROR
```

### 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看容器详细信息
docker-compose exec app top
```

## 故障排查

### 常见问题

1. **服务启动失败**
   ```bash
   # 查看详细日志
   docker-compose logs app
   
   # 检查端口占用
   netstat -tulpn | grep :80
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库状态
   docker-compose exec mysql mysqladmin ping -h localhost -u root -proot
   
   # 查看数据库日志
   docker-compose logs mysql
   ```

3. **前端无法访问**
   ```bash
   # 检查nginx配置
   docker-compose exec app nginx -t
   
   # 检查前端文件
   docker-compose exec app ls -la /var/www/html
   ```

### 调试模式

```bash
# 以调试模式启动
docker-compose up

# 进入容器调试
docker-compose exec app sh
```

## 安全配置

### 生产环境安全建议

1. **修改默认密码**
   - 数据库密码
   - Redis密码
   - JWT密钥

2. **网络安全**
   - 使用防火墙限制端口访问
   - 配置SSL证书
   - 使用反向代理

3. **容器安全**
   - 定期更新镜像
   - 使用非root用户运行
   - 限制容器资源使用

### SSL配置

```bash
# 生成SSL证书
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx.key -out nginx.crt

# 配置nginx SSL
# 修改 docker/nginx.conf 添加SSL配置
```

## 扩展部署

### 负载均衡

```bash
# 使用nginx作为负载均衡器
# 配置多个应用实例
```

### 数据库集群

```bash
# 配置MySQL主从复制
# 使用Redis集群
```

## 总结

Docker 部署提供了完整的容器化解决方案，支持开发和生产环境，具有以下优势：

- **简单部署**: 一键启动所有服务
- **环境一致**: 开发和生产环境一致
- **易于维护**: 容器化隔离，便于管理
- **可扩展**: 支持水平扩展和负载均衡
- **安全可靠**: 包含健康检查和自动重启

通过合理的配置和管理，可以构建稳定可靠的 FoodCook 应用部署环境。 