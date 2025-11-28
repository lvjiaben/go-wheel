# 功能概览

## 内置功能

### 用户认证

- 账号密码登录
- 手机验证码登录
- 图形验证码
- JWT Token 认证
- 自动刷新 Token

### 权限控制

- 路由权限 - 根据用户角色动态生成路由
- 按钮权限 - 细粒度的按钮级权限控制
- 接口权限 - 后端接口权限校验

### 布局系统

- 多种布局模式（侧边栏、顶部导航、混合）
- 主题切换（亮色、暗色、跟随系统）
- 响应式设计
- 多标签页

### 国际化

- 中英文切换
- 组件国际化
- 路由国际化
- 后端消息国际化

## 页面功能

### 管理后台 (Admin)

| 页面 | 路径 | 说明 |
|------|------|------|
| 仪表盘 | /dashboard | 数据概览 |
| 用户管理 | /system/user | 用户增删改查 |
| 角色管理 | /system/role | 角色权限分配 |
| 菜单管理 | /system/menu | 动态菜单配置 |
| 配置管理 | /system/config | 系统参数配置 |
| 附件管理 | /system/attachment | 文件管理 |
| 代码生成 | /devtools/codegen | 可视化代码生成 |

### 用户端 (User)

| 页面 | 路径 | 说明 |
|------|------|------|
| 登录 | /login | 用户登录 |
| 注册 | /register | 用户注册 |
| 首页 | /home | 用户首页 |
| 个人中心 | /home/userinfo | 个人信息管理 |

## API 接口

### 认证接口

```typescript
// 登录
loginApi(data: LoginParams): Promise<LoginResult>

// 获取用户信息
getUserInfoApi(): Promise<UserInfo>

// 获取菜单
getMenusApi(): Promise<Menu[]>

// 登出
logoutApi(): Promise<void>
```

### 用户接口

```typescript
// 修改密码
changePwdApi(data: ChangePwdParams): Promise<void>

// 修改手机号
changeMobileApi(data: ChangeMobileParams): Promise<void>

// 发送验证码
sendSmsApi(data: SendSmsParams): Promise<void>
```

## 自定义组件

项目内置了多个自定义组件，详见 [自定义组件](./components/array-form-table.md) 章节：

- **ArrayFormTable** - 数组表单表格
- **AttachmentInput** - 附件输入框
- **RichEditor** - 富文本编辑器
- **TableSelect** - 表格选择器
- **ImageCaptcha** - 图形验证码

## 适配器模式

项目使用适配器模式统一组件接口：

```typescript
// adapter/component/index.ts
export {
  Button,
  message,
  notification,
  // ...
} from 'ant-design-vue'
```

这样可以方便地切换 UI 组件库，只需修改适配器即可。

## 状态管理

使用 Pinia 进行状态管理：

```typescript
// stores/user.ts
export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: null,
    token: '',
  }),
  actions: {
    async login(params) {
      // ...
    },
  },
})
```

## 更多功能

更多功能请参考 [Vben Admin 官方文档](https://doc.vben.pro/)。

