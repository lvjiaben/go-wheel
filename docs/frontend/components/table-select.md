# TableSelect 表格选择器

远程数据选择组件，支持分页、搜索、自定义渲染。

## 组件注册

组件已通过 `adapter/component` 全局注册，在表单中直接使用组件名即可，**无需手动 import**：

```typescript
{
  component: 'TableSelect',
  fieldName: 'userId',
  // ...
}
```

## 基础用法

```vue
<script setup lang="ts">
import { ref } from 'vue';
import TableSelect from '@/components/table-select/index.vue';

const userId = ref('');
const config = {
  api: '/backend/user/list',
  labelField: 'username',
  valueField: 'id',
};
</script>

<template>
  <TableSelect v-model:value="userId" :config="config" />
</template>
```

## 多选模式

```vue
<template>
  <TableSelect v-model:value="userIds" :config="config" multiple />
</template>
```

## 带图片和描述

```vue
<script setup lang="ts">
const config = {
  api: '/backend/user/list',
  labelField: 'username',
  valueField: 'id',
  imageField: 'avatar',
  descField: 'email',
  pageSize: 20,
};
</script>

<template>
  <TableSelect v-model:value="userId" :config="config" />
</template>
```

## 在表单中使用

```vue
<script setup lang="ts">
import { useVbenForm } from '@vben/common-ui';
import TableSelect from '@/components/table-select/index.vue';

const [Form] = useVbenForm({
  schema: [
    {
      fieldName: 'user_id',
      label: '选择用户',
      component: TableSelect,
      componentProps: {
        config: JSON.stringify({
          api: '/backend/user/list',
          labelField: 'username',
          valueField: 'id',
        }),
      },
    },
  ],
});
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| value | `string \| string[] \| number` | - | 绑定值 |
| config | `string \| RemoteConfig` | - | 配置对象或 JSON 字符串 |
| multiple | `boolean` | `false` | 是否多选 |
| placeholder | `string` | - | 占位文本 |
| disabled | `boolean` | `false` | 是否禁用 |

## Config 配置

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| api | `string` | - | 数据接口地址（必填） |
| searchField | `string` | `'search'` | 搜索参数字段名 |
| labelField | `string` | `'name'` | 显示文本字段 |
| valueField | `string` | `'id'` | 值字段 |
| imageField | `string` | `'image'` | 图片字段 |
| descField | `string` | `'description'` | 描述字段 |
| pageSize | `number` | `10` | 每页条数 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:value | `(value: string \| string[] \| number)` | 值更新时触发 |
| change | `(value: string \| string[] \| number)` | 值改变时触发 |

## 功能特性

### 远程搜索

输入关键词后自动搜索，支持防抖（300ms）。

### 分页加载

当数据总数超过 pageSize 时，下拉框底部显示分页器。

### 选项缓存

已选中的选项会被缓存，切换分页后仍能正确显示标签。

### 自定义渲染

选项支持显示图片和描述信息：
- 左侧显示图片（如果配置了 imageField）
- 右侧显示标签和描述

## API 响应格式

接口需返回以下格式：

```json
{
  "list": [
    { "id": 1, "name": "选项1", "image": "url", "description": "描述" }
  ],
  "total": 100
}
```

## 国际化

组件使用以下国际化 key：

```json
{
  "common.components.tableSelect.placeholder": "请选择",
  "common.components.tableSelect.noData": "暂无数据",
  "common.components.tableSelect.fetchError": "获取数据失败"
}
```

## 使用场景

- 用户选择器
- 商品选择器
- 分类选择器
- 任何需要远程分页搜索的选择场景

