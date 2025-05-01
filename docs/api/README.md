# API文档

## 概述

本文档描述了Go Admin Framework的API接口规范。

## 基础信息

- 基础URL: `http://localhost:8081/api/v1`
- 认证方式: JWT Bearer Token
- 响应格式: JSON

## 通用响应格式

```json
{
    "code": 0,           // 状态码，0表示成功
    "message": "success", // 响应消息
    "data": {}           // 响应数据
}
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 认证失败 |
| 1003 | 权限不足 |
| 1004 | 资源不存在 |
| 1005 | 服务器错误 |

## API列表

### 认证相关

#### 登录
- 请求路径: `/auth/login`
- 请求方法: POST
- 请求参数:
```json
{
    "username": "string",
    "password": "string"
}
```
- 响应示例:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "token": "string",
        "expires_at": "2024-01-01T00:00:00Z"
    }
}
```

### 用户管理

#### 获取用户列表
- 请求路径: `/users`
- 请求方法: GET
- 请求参数:
  - page: 页码
  - page_size: 每页数量
- 响应示例:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "total": 100,
        "list": [
            {
                "id": 1,
                "username": "string",
                "email": "string",
                "created_at": "2024-01-01T00:00:00Z"
            }
        ]
    }
}
```

## 注意事项

1. 所有API请求都需要在Header中携带Authorization字段
2. 分页参数page从1开始
3. 时间格式统一使用ISO 8601标准
4. 所有请求和响应都需要使用UTF-8编码 