# AttachmentInput 附件输入

附件选择输入组件，集成附件管理弹窗，支持图片预览。

## 组件注册

组件已通过 `adapter/component` 全局注册，在表单中直接使用组件名即可，**无需手动 import**：

```typescript
{
  component: 'AttachmentInput',
  fieldName: 'image',
  // ...
}
```

## 基础用法

```vue
<script setup lang="ts">
import { ref } from 'vue';
import AttachmentInput from '@/components/attachment-input/index.vue';

const imageUrl = ref('');
</script>

<template>
  <AttachmentInput v-model:value="imageUrl" />
</template>
```

## 多选模式

```vue
<template>
  <AttachmentInput v-model:value="images" multiple />
</template>
```

多选时，值以英文逗号分隔：`url1,url2,url3`

## 仅按钮模式

隐藏输入框，只显示选择按钮：

```vue
<template>
  <AttachmentInput v-model:value="imageUrl" :show-input="false" />
</template>
```

## 在表单中使用

```vue
<script setup lang="ts">
import { useVbenForm } from '@vben/common-ui';
import AttachmentInput from '@/components/attachment-input/index.vue';

const [Form] = useVbenForm({
  schema: [
    {
      fieldName: 'avatar',
      label: '头像',
      component: AttachmentInput,
      componentProps: {
        placeholder: '请选择头像图片',
      },
    },
    {
      fieldName: 'gallery',
      label: '图片集',
      component: AttachmentInput,
      componentProps: {
        multiple: true,
      },
    },
  ],
});
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| value | `string` | `''` | 绑定值（v-model:value） |
| modelValue | `string` | `''` | 绑定值（v-model） |
| multiple | `boolean` | `false` | 是否多选 |
| placeholder | `string` | - | 输入框占位文本 |
| showPreview | `boolean` | `true` | 是否显示图片预览 |
| showInput | `boolean` | `true` | 是否显示输入框 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:modelValue | `(value: string)` | v-model 更新 |
| update:value | `(value: string)` | v-model:value 更新 |
| change | `(value: string)` | 值改变时触发 |

## 功能特性

### 图片预览

- 自动识别图片类型（jpg、png、gif、webp 等）
- 图片显示缩略图预览
- 非图片文件显示文件扩展名
- 点击预览图可在新窗口打开原图

### 删除功能

鼠标悬停在预览图上时，显示删除按钮，点击可移除该文件。

### 附件选择器

点击「选择文件」按钮打开附件管理弹窗：
- 支持分组筛选
- 支持搜索
- 支持上传新文件
- 支持多选

## 国际化

组件使用以下国际化 key：

```json
{
  "system.attachment.inputPlaceholder": "请输入或选择附件",
  "system.attachment.selectFile": "选择文件"
}
```

## 使用场景

- 用户头像上传
- 商品图片管理
- 文章封面图选择
- 多图上传场景

