# ImageCaptcha 图形验证码

图形验证码输入组件，集成验证码生成和刷新功能。

## 组件注册

组件已通过 `adapter/component` 全局注册，在表单中直接使用组件名即可，**无需手动 import**：

```typescript
{
  component: 'ImageCaptcha',
  fieldName: 'captcha',
  // ...
}
```

## 基础用法

```vue
<script setup lang="ts">
import { ref } from 'vue';
import ImageCaptcha from '@/components/form/image-captcha.vue';
import { getCaptchaApi } from '@/api/core/auth';

const captchaValue = ref({ id: '', code: '' });
</script>

<template>
  <ImageCaptcha v-model="captchaValue" :captcha-api="getCaptchaApi" />
</template>
```

## 在登录表单中使用

```vue
<script setup lang="ts">
import { useVbenForm } from '@vben/common-ui';
import ImageCaptcha from '@/components/form/image-captcha.vue';
import { getCaptchaApi } from '@/api/core/auth';

const [Form, { validate }] = useVbenForm({
  schema: [
    {
      fieldName: 'username',
      label: '用户名',
      component: 'Input',
    },
    {
      fieldName: 'password',
      label: '密码',
      component: 'InputPassword',
    },
    {
      fieldName: 'captcha',
      label: '验证码',
      component: ImageCaptcha,
      componentProps: {
        captchaApi: getCaptchaApi,
      },
    },
  ],
});

const handleLogin = async () => {
  const values = await validate();
  // values.captcha = { id: 'xxx', code: 'xxxx' }
  await loginApi({
    username: values.username,
    password: values.password,
    captcha_id: values.captcha.id,
    captcha_code: values.captcha.code,
  });
};
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| modelValue | `CaptchaValue \| string` | - | 绑定值 |
| captchaApi | `() => Promise<any>` | - | 获取验证码的 API 函数 |

## CaptchaValue 类型

```typescript
interface CaptchaValue {
  id: string;   // 验证码 ID
  code: string; // 用户输入的验证码
}
```

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: CaptchaValue)` | 值更新时触发 |

## Expose 方法

| 方法名 | 参数 | 返回值 | 说明 |
|--------|------|--------|------|
| getCaptchaId | - | `string` | 获取当前验证码 ID |
| refresh | - | `void` | 刷新验证码 |

```vue
<script setup lang="ts">
import { ref } from 'vue';
import ImageCaptcha from '@/components/form/image-captcha.vue';

const captchaRef = ref();

// 手动刷新验证码
const refreshCaptcha = () => {
  captchaRef.value?.refresh();
};

// 获取验证码 ID
const getCaptchaId = () => {
  return captchaRef.value?.getCaptchaId();
};
</script>

<template>
  <ImageCaptcha ref="captchaRef" v-model="captchaValue" :captcha-api="getCaptchaApi" />
</template>
```

## API 响应格式

captchaApi 需返回以下格式：

```json
{
  "captcha_id": "uuid-string",
  "captcha_data": "data:image/png;base64,..."
}
```

## 功能特性

### 自动加载

组件挂载时自动调用 captchaApi 生成验证码。

### 点击刷新

点击验证码图片可刷新获取新的验证码。

### 值同步

输入验证码时，自动将 id 和 code 组合成对象输出。

## 国际化

组件使用以下国际化 key：

```json
{
  "page.auth.imageCaptcha": "图形验证码",
  "page.auth.imageCaptchaTip": "请输入验证码"
}
```

## 样式说明

- 验证码图片宽度：120px
- 验证码图片高度：32px
- 输入框使用 large 尺寸
- 最大输入长度：6 位

## 使用场景

- 登录页面
- 注册页面
- 找回密码
- 敏感操作验证

