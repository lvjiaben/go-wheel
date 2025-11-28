# 前端介绍

前端基于 [Vben Admin](https://doc.vben.pro/) 开发，采用 Vue 3 + TypeScript + Vite 5 技术栈。

## 项目结构

```
vben-admin/
├── apps/                      # 应用目录
│   ├── admin/                 # 管理后台应用
│   │   ├── src/
│   │   │   ├── api/          # API 接口
│   │   │   ├── views/        # 页面组件
│   │   │   ├── router/       # 路由配置
│   │   │   └── locales/      # 国际化
│   │   └── index.html
│   │
│   └── user/                  # 用户端应用
│       ├── src/
│       │   ├── api/          # API 接口
│       │   │   └── core/     # 核心接口
│       │   ├── components/   # 自定义组件 ⭐
│       │   ├── views/        # 页面组件
│       │   ├── router/       # 路由配置
│       │   └── locales/      # 国际化
│       └── index.html
│
├── packages/                  # 公共包
│   ├── @core/                # 核心功能
│   ├── effects/              # 副作用
│   ├── icons/                # 图标
│   ├── locales/              # 国际化
│   └── stores/               # 状态管理
│
└── internal/                  # 内部工具
```

## 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| Vue | 3.4+ | 渐进式 JavaScript 框架 |
| TypeScript | 5.0+ | JavaScript 超集 |
| Vite | 5.0+ | 下一代前端构建工具 |
| Ant Design Vue | 4.0+ | UI 组件库 |
| Pinia | 2.0+ | 状态管理 |
| Vue Router | 4.0+ | 路由管理 |
| VueUse | 10.0+ | Vue 组合式 API 工具集 |

## 应用说明

### Admin 应用

管理后台应用，包含：

- 用户管理
- 角色管理
- 菜单管理
- 配置管理
- 附件管理
- 代码生成

### User 应用

用户端应用，包含：

- 用户登录/注册
- 个人中心
- 修改密码
- 修改手机号

## 开发命令

```bash
# 安装依赖
pnpm install

# 启动管理后台
pnpm dev:admin

# 启动用户端
pnpm dev:user

# 构建管理后台
pnpm build:admin

# 构建用户端
pnpm build:user

# 代码检查
pnpm lint

# 类型检查
pnpm typecheck
```

## 环境变量

### 开发环境 `.env.development`

```env
# API 地址
VITE_GLOB_API_URL=/api

# 是否开启 Mock
VITE_USE_MOCK=false
```

### 生产环境 `.env.production`

```env
# API 地址
VITE_GLOB_API_URL=https://api.example.com

# 是否开启压缩
VITE_BUILD_COMPRESS=gzip
```

## 更多功能

本项目基于 Vben Admin 开发，更多功能请参考官方文档：

- [Vben Admin 官方文档](https://doc.vben.pro/)
- [组件文档](https://doc.vben.pro/components/)
- [主题配置](https://doc.vben.pro/guide/essentials/settings.html)
- [权限控制](https://doc.vben.pro/guide/in-depth/access.html)

