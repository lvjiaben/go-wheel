# Go Admin Framework

A Golang-based admin framework providing comprehensive management functions and rich component support.

<p align="center">
  <a href="./docs/guide/README_EN.md">
    <img src="https://img.shields.io/badge/Docs-View_Documentation-blue.svg" alt="View Documentation">
  </a>
  <a href="./docs/guide/README.md">
    <img src="https://img.shields.io/badge/文档-中文文档-blue.svg" alt="中文文档">
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License">
  </a>
  <a href="https://golang.org">
    <img src="https://img.shields.io/badge/Golang-1.23+-orange.svg" alt="Golang">
  </a>
</p>

English | [简体中文](./README.md)

## Features

- **Authentication & Authorization**: Complete JWT authentication and RBAC permission management
- **Multi-language Support**: Internationalization based on i18n
- **Message Queue**: RabbitMQ integration with support for standard and delayed queues
- **Scheduled Tasks**: Cron-based task scheduler
- **Cache Support**: Redis cache management supporting various data structures
- **Health Check**: Service health status monitoring
- **Logging**: Structured log recording and management
- **Parameter Validation**: Request parameter validation
- **API Documentation**: Automatic API documentation with Swagger
- **Database**: Gorm integration supporting multiple databases

## Quick Start

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Configure Environment
Modify the `configs/config.yaml` configuration file to set database, Redis, and other connection information.

### 3. Start the Service
```bash
go run main.go
```

### 4. Access API Documentation
```
http://localhost:8080/swagger/index.html
```

## Core Components

The framework integrates various commonly used components to make development more convenient:

- **Redis Cache**: For data caching, counters, distributed locks, etc.
- **RabbitMQ Queue**: For asynchronous processing, task distribution, event broadcasting, etc.
- **Cron Scheduler**: For regularly executing tasks, data statistics, report generation, etc.
- **Health Check**: For monitoring service health status with customizable check strategies

<p align="center">
  <a href="./docs/guide/components_en.md">
    <img src="https://img.shields.io/badge/View-Component_Docs-green.svg" alt="View Component Docs">
  </a>
</p>

## Directory Structure

```
admin/
├── app/               # Application code
│   ├── api/           # API controllers
│   │   ├── v1/        # API version 1
│   │   └── v2/        # API version 2
│   ├── backend/       # Backend services
│   │   ├── controller/ # Controllers
│   │   ├── service/    # Service layer
│   │   └── model/      # Data models
│   └── generate/      # Code generation
├── configs/           # Configuration files
│   └── config.yaml    # Main config file
├── deployments/       # Deployment-related
├── docs/              # Documentation
│   ├── api/           # API documentation
│   └── guide/         # User guides
├── pkg/               # Public packages
│   ├── auth/          # Authentication
│   ├── cache/         # Cache
│   ├── cron/          # Scheduled tasks
│   ├── db/            # Database
│   ├── logger/        # Logging
│   ├── mq/            # Message queue
│   └── validator/     # Validator
├── routes/            # Route definitions
├── scripts/           # Scripts
└── tests/             # Tests
    ├── cron/          # Scheduled task tests
    ├── health/        # Health check tests
    ├── rabbitmq/      # Message queue tests
    └── redis/         # Redis cache tests
```

## Development Guide

### Adding a New API

1. Define API controllers in the `app/api/v1` directory
2. Implement business logic in `app/backend/service`
3. Define data models in `app/backend/model`
4. Register routes in the `routes` directory

### Using Framework Components

All core components are used via dependency injection, no manual initialization required:

```go
// Using Redis in a service
type UserService struct {
    redisService *service.RedisService
}

func NewUserService(redisService *service.RedisService) *UserService {
    return &UserService{
        redisService: redisService,
    }
}

// Using Redis cache
func (s *UserService) GetUserByID(id string) (*User, error) {
    cacheKey := "user:" + id
    // Try to get from cache
    data, err := s.redisService.Get(cacheKey)
    // ...
}
```

## Deployment

### Docker Deployment
```bash
docker build -t go-admin .
docker run -p 8080:8080 go-admin
```

### Kubernetes Deployment
```bash
kubectl apply -f deployments/k8s/
```

## Documentation Center

<p align="center">
  <a href="./docs/guide/README_EN.md">
    <img src="https://img.shields.io/badge/User_Guide-Guide-blue.svg" alt="User Guide">
  </a>
  <a href="./docs/api/README_EN.md">
    <img src="https://img.shields.io/badge/API_Docs-API-blue.svg" alt="API Docs">
  </a>
  <a href="./docs/guide/components_en.md">
    <img src="https://img.shields.io/badge/Component_Docs-Components-blue.svg" alt="Component Docs">
  </a>
  <a href="./docs/guide/configuration_en.md">
    <img src="https://img.shields.io/badge/Configuration-Config-blue.svg" alt="Configuration">
  </a>
</p>

## License

[MIT License](./LICENSE) 