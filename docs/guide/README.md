# Go Admin Framework 使用指南

## 环境要求

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- Redis 6.0 或更高版本
- Docker（可选）
- Kubernetes（可选）

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/your-username/go-admin.git
cd go-admin
```

### 2. 配置环境

1. 复制配置文件模板
```bash
cp configs/config.yaml.example configs/config.yaml
```

2. 修改配置文件
```yaml
app:
  name: go-admin
  mode: development
  port: 8081
  version: 1.0.0

mysql:
  host: localhost
  port: 3306
  user: root
  password: your-password
  dbname: go_admin
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: localhost
  port: 6379
  password: your-password
  db: 0
```

### 3. 初始化数据库

```bash
mysql -u root -p < scripts/sql/init.sql
```

### 4. 运行应用

```bash
# 开发模式
go run cmd/server/main.go

# 生产模式
go build -o main cmd/server/main.go
./main
```

## 开发指南

### 1. 创建新的API

1. 在`api/v1`目录下定义API接口
2. 在`internal/app/controller`中实现控制器
3. 在`internal/app/service`中实现业务逻辑
4. 在`internal/app/model`中定义数据模型
5. 在`internal/app/repository`中实现数据访问

### 2. 添加新的中间件

1. 在`pkg/middleware`目录下创建新的中间件
2. 在`internal/server`中注册中间件

### 3. 添加新的工具函数

1. 在`pkg/utils`目录下添加工具函数
2. 在`pkg/errors`中定义错误类型

## 测试

### 1. 单元测试

```bash
go test ./...
```

### 2. 集成测试

```bash
go test -tags=integration ./...
```

### 3. 端到端测试

```bash
go test -tags=e2e ./tests/e2e/...
```

## 部署

### 1. Docker部署

```bash
# 构建镜像
./scripts/build/build.sh

# 运行容器
docker run -p 8081:8081 go-admin:1.0.0
```

### 2. Kubernetes部署

```bash
# 部署应用
./scripts/deploy/deploy.sh
```

## 常见问题

### 1. 数据库连接失败

检查MySQL服务是否正常运行，以及配置文件中的数据库连接信息是否正确。

### 2. Redis连接失败

检查Redis服务是否正常运行，以及配置文件中的Redis连接信息是否正确。

### 3. 权限问题

确保应用有足够的权限访问所需的资源，如数据库、Redis等。

## 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License 