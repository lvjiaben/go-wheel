# Go Admin Framework User Guide

<p align="center">
  <a href="../../README_EN.md">
    <img src="https://img.shields.io/badge/Back-Home-blue.svg" alt="Back to Home">
  </a>
  <a href="./README.md">
    <img src="https://img.shields.io/badge/切换-中文-blue.svg" alt="中文">
  </a>
</p>

Welcome to the Go Admin Framework! This guide will help you understand how to use the various features and components of the framework.

## Table of Contents

- [Quick Start](#quick-start)
- [Basic Concepts](#basic-concepts)
- [Component Usage](./components_en.md)
  - [Redis Cache](./components/redis_en.md)
  - [RabbitMQ Message Queue](./components/rabbitmq_en.md)
  - [Cron Scheduler](./components/cron_en.md)
  - [Health Check](./components/health_en.md)
- [Best Practices](./best_practices_en.md)
- [Configuration](./configuration_en.md)
- [FAQ](./faq_en.md)

## Quick Start

### Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/yourusername/go-admin.git
cd go-admin
go mod tidy
```

### Configuration

Copy the example configuration file and modify it:

```bash
cp configs/config.example.yaml configs/config.yaml
```

Edit the `config.yaml` file to set database, Redis, and other connection information.

### Start the Service

```bash
go run main.go
```

Visit [http://localhost:8080](http://localhost:8080) to access the admin interface.

## Basic Concepts

The Go Admin Framework uses a clear layered architecture:

1. **Controller Layer**: Handles HTTP requests, validates parameters, calls the Service layer
2. **Service Layer**: Implements business logic, calls the Model layer and other services
3. **Model Layer**: Defines data structures and database interactions

### Directory Structure

```
app/
├── api/              # API controllers
│   ├── v1/           # API version 1
│   └── v2/           # API version 2
└── backend/          # Backend services
    ├── controller/   # Controllers
    ├── service/      # Service layer
    └── model/        # Data models
```

### Core Components

The Go Admin Framework integrates various commonly used components that are initialized when the framework starts. You only need to inject them where needed:

- **RedisService**: Redis cache service
- **MQService**: RabbitMQ message queue service
- **CronService**: Scheduled task service
- **HealthService**: Health check service

All these services are used via dependency injection, with no manual initialization required. For detailed usage, please refer to the [Component Usage Documentation](./components_en.md).

## Development Workflow

### 1. Create a New API

Create a new controller in the `app/api/v1` directory:

```go
// app/api/v1/user_controller.go
package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/backend/service"
)

type UserController struct {
    userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// GetUser gets user information
// @Summary Get user information
// @Description Get user information by user ID
// @Tags User Management
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} model.User
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
    id := ctx.Param("id")
    user, err := c.userService.GetUserByID(id)
    if err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(200, user)
}
```

### 2. Implement Business Logic

Create the corresponding service in the `app/backend/service` directory:

```go
// app/backend/service/user_service.go
package service

import (
    "encoding/json"
    "time"
    "github.com/lvjiaben/go-wheel/app/backend/model"
)

type UserService struct {
    redisService *RedisService
}

func NewUserService(redisService *RedisService) *UserService {
    return &UserService{
        redisService: redisService,
    }
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
    // Try to get from cache
    cacheKey := "user:" + id
    userData, err := s.redisService.Get(cacheKey)
    
    // Cache miss, get from database
    if err != nil {
        // Get user data from database
        user := &model.User{
            ID: id,
            // ... other fields
        }
        
        // Store user data in cache
        jsonData, _ := json.Marshal(user)
        s.redisService.Set(cacheKey, string(jsonData), 30*time.Minute)
        
        return user, nil
    }
    
    // Parse cached data
    var user model.User
    json.Unmarshal([]byte(userData), &user)
    
    return &user, nil
}
```

### 3. Register Routes

Register API routes in the `routes` directory:

```go
// routes/api.go
package routes

import (
    "github.com/gin-gonic/gin"
    v1 "github.com/lvjiaben/go-wheel/app/api/v1"
)

func RegisterAPIRoutes(r *gin.Engine, userController *v1.UserController) {
    api := r.Group("/api")
    {
        v1Group := api.Group("/v1")
        {
            users := v1Group.Group("/users")
            {
                users.GET("/:id", userController.GetUser)
                // Other user-related routes
            }
        }
    }
}
```

## Next Steps

- Check the [Component Usage Documentation](./components_en.md) to learn how to use the various components provided by the framework
- Check the [Configuration Guide](./configuration_en.md) to learn how to configure the framework
- Check the [Best Practices](./best_practices_en.md) to learn how to better use the framework

If you have any questions, please check the [FAQ](./faq_en.md) or submit an Issue. 