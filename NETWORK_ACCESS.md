# 网络访问配置指南

## 问题描述

用户需要在局域网内的其他设备（如手机）访问FoodCook应用，但遇到了跨域访问问题。

## 解决方案

### 1. 前端配置（Vite）

修改 `web/vite.config.js` 文件：

```javascript
export default defineConfig({
  // ... 其他配置
  server: {
    port: 3000,
    host: '0.0.0.0', // 绑定到所有网络接口
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
})
```

### 2. 后端CORS配置

修改 `internal/pkg/config/config.go` 文件中的默认CORS配置：

```go
viper.SetDefault("cors.allowed_origins", []string{"*"})
```

或者创建 `configs/config.yaml` 配置文件：

```yaml
cors:
  allowed_origins:
    - "*"  # 开发环境允许所有来源
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "*"
  allow_credentials: true
```

### 3. 防火墙配置

确保防火墙允许相关端口：

```bash
# macOS (如果使用防火墙)
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/bin/node
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblock /usr/bin/node

# 或者临时关闭防火墙（仅开发环境）
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --setglobalstate off
```

### 4. 网络访问测试

#### 获取本机IP地址
```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
# 或者
ipconfig getifaddr en0  # macOS
```

#### 测试前端访问
```bash
curl -X GET "http://192.168.2.184:3000" -I
```

#### 测试后端API访问
```bash
curl -X GET "http://192.168.2.184:8080/api/dishes" -H "Content-Type: application/json"
```

## 访问地址

- **前端应用**: `http://192.168.2.184:3000`
- **后端API**: `http://192.168.2.184:8080`

## 安全注意事项

### 开发环境
- 允许所有来源的CORS请求（`"*"`）
- 仅用于局域网内测试

### 生产环境
- 限制CORS来源为特定域名
- 启用HTTPS
- 配置适当的防火墙规则

```yaml
# 生产环境CORS配置示例
cors:
  allowed_origins:
    - "https://yourdomain.com"
    - "https://www.yourdomain.com"
  allowed_methods:
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
    - "OPTIONS"
  allowed_headers:
    - "Content-Type"
    - "Authorization"
    - "X-Requested-With"
  allow_credentials: true
```

## 常见问题

### 1. 403 Forbidden 错误
- 检查CORS配置是否正确
- 确认防火墙设置
- 验证IP地址是否正确

### 2. 连接超时
- 检查网络连接
- 确认服务是否正在运行
- 验证端口是否被占用

### 3. 跨域问题
- 重启后端服务以应用新的CORS配置
- 检查浏览器控制台错误信息
- 确认前端代理配置

## 启动命令

```bash
# 启动后端服务
go run cmd/foodcook/main.go

# 启动前端服务
cd web && npm run dev
```

## 验证步骤

1. 在本机访问 `http://localhost:3000` 确认前端正常
2. 在本机访问 `http://localhost:8080/api/dishes` 确认后端正常
3. 在局域网设备访问 `http://192.168.2.184:3000` 确认网络访问正常
4. 在局域网设备测试API调用确认CORS正常 