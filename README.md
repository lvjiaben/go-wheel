# Go Admin Framework

一个基于Golang的后台管理框架，提供完整的后台管理功能。

## 目录结构

```
admin/
├── api/                # API接口定义
│   ├── v1/            # API版本1
│   └── swagger/       # Swagger文档
├── cmd/               # 命令行入口
│   └── server/        # 服务器启动入口
├── configs/           # 配置文件
│   ├── config.yaml    # 主配置文件
│   └── i18n/          # 国际化文件
├── deployments/       # 部署相关
│   ├── docker/        # Docker配置
│   └── k8s/           # Kubernetes配置
├── docs/             # 文档
│   ├── api/          # API文档
│   └── guide/        # 使用指南
├── internal/         # 内部代码
│   ├── app/          # 应用层
│   │   ├── controller/  # 控制器
│   │   ├── service/     # 服务层
│   │   ├── model/       # 数据模型
│   │   └── repository/  # 数据访问层
│   ├── pkg/          # 内部包
│   │   ├── auth/     # 认证授权
│   │   ├── cache/    # 缓存
│   │   ├── cron/     # 定时任务
│   │   ├── i18n/     # 国际化
│   │   ├── logger/   # 日志
│   │   ├── queue/    # 消息队列
│   │   └── validator/ # 验证器
│   └── server/       # 服务器
├── pkg/              # 公共包
│   ├── errors/       # 错误处理
│   ├── middleware/   # 中间件
│   └── utils/        # 工具函数
├── scripts/          # 脚本
│   ├── build/        # 构建脚本
│   └── deploy/       # 部署脚本
└── tests/            # 测试
    ├── api/          # API测试
    └── e2e/          # 端到端测试
```

## 功能特性

- 完整的后台管理功能
- 基于JWT的认证授权
- 多语言支持
- 消息队列
- 定时任务
- 缓存支持
- 日志记录
- 参数验证
- Swagger API文档

## 快速开始

1. 安装依赖
```bash
go mod tidy
```

2. 启动服务
```bash
go run cmd/server/main.go
```

3. 访问API文档
```
http://localhost:8081/swagger/index.html
```

## 开发指南

1. 添加新的API
   - 在`api/v1`目录下定义API接口
   - 在`internal/app/controller`中实现控制器
   - 在`internal/app/service`中实现业务逻辑
   - 在`internal/app/model`中定义数据模型
   - 在`internal/app/repository`中实现数据访问

2. 添加新的中间件
   - 在`pkg/middleware`目录下创建新的中间件
   - 在`internal/server`中注册中间件

3. 添加新的工具函数
   - 在`pkg/utils`目录下添加工具函数
   - 在`pkg/errors`中定义错误类型

## 部署

1. Docker部署
```bash
docker build -t go-admin .
docker run -p 8081:8081 go-admin
```

2. Kubernetes部署
```bash
kubectl apply -f deployments/k8s/
```

## 贡献指南

1. Fork项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License
