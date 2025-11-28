# RichEditor 富文本编辑器

基于 Quill 的富文本编辑器组件，集成附件管理实现图片上传。

## 组件注册

组件已通过 `adapter/component` 全局注册，在表单中直接使用组件名即可，**无需手动 import**：

```typescript
{
  component: 'RichEditor',
  fieldName: 'content',
  // ...
}
```

## 基础用法

```vue
<script setup lang="ts">
import { ref } from 'vue';
import RichEditor from '@/components/rich-editor/index.vue';

const content = ref('<p>Hello World</p>');
</script>

<template>
  <RichEditor v-model:value="content" />
</template>
```

## 自定义高度

```vue
<template>
  <RichEditor v-model:value="content" :min-height="500" />
</template>
```

## 禁用状态

```vue
<template>
  <RichEditor v-model:value="content" disabled />
</template>
```

## 在表单中使用

```vue
<script setup lang="ts">
import { useVbenForm } from '@vben/common-ui';
import RichEditor from '@/components/rich-editor/index.vue';

const [Form] = useVbenForm({
  schema: [
    {
      fieldName: 'content',
      label: '文章内容',
      component: RichEditor,
      componentProps: {
        minHeight: 400,
        placeholder: '请输入文章内容',
      },
    },
  ],
});
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| value | `string` | `''` | 绑定值（HTML 字符串） |
| placeholder | `string` | - | 占位文本 |
| minHeight | `number` | `300` | 编辑器最小高度（px） |
| disabled | `boolean` | `false` | 是否禁用 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:value | `(value: string)` | 内容更新时触发 |
| change | `(value: string)` | 内容改变时触发 |

## 工具栏功能

编辑器内置完整的工具栏：

| 功能 | 说明 |
|------|------|
| 文字样式 | 加粗、斜体、下划线、删除线 |
| 引用代码 | 引用块、代码块 |
| 标题 | H1-H6 标题 |
| 列表 | 有序列表、无序列表 |
| 上下标 | 上标、下标 |
| 缩进 | 增加/减少缩进 |
| 文字方向 | RTL 支持 |
| 字号 | 小、正常、大、超大 |
| 颜色 | 文字颜色、背景颜色 |
| 字体 | 字体选择 |
| 对齐 | 左对齐、居中、右对齐、两端对齐 |
| 清除格式 | 清除所有格式 |
| 媒体 | 链接、图片、视频 |

## 图片上传

点击工具栏的图片按钮，会打开附件管理弹窗：

1. 从已上传的附件中选择图片
2. 或上传新图片
3. 支持多选，一次插入多张图片
4. 图片插入到当前光标位置

## 国际化

组件使用以下国际化 key：

```json
{
  "common.components.richEditor.placeholder": "请输入内容..."
}
```

## 样式定制

组件使用 Quill 的 Snow 主题，可通过 CSS 变量定制：

```css
.rich-editor-wrapper {
  /* 边框颜色 */
  border-color: #d9d9d9;
}

.rich-editor-wrapper:hover {
  /* 悬停边框颜色 */
  border-color: #4096ff;
}
```

## 使用场景

- 文章内容编辑
- 商品描述编辑
- 公告内容编辑
- 邮件模板编辑

