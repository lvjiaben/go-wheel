# ArrayFormTable 数组表单

用于编辑键值对数据的表格组件，支持拖拽排序。

## 组件注册

组件已通过 `adapter/component` 全局注册，在表单中直接使用组件名即可，**无需手动 import**：

```typescript
{
  component: 'ArrayFormTable',
  fieldName: 'items',
  // ...
}
```

## 基础用法

```vue
<script setup lang="ts">
import { ref } from 'vue';
import ArrayFormTable from '@/components/array-form-table/index.vue';

const formData = ref('{"key1":"value1","key2":"value2"}');
</script>

<template>
  <ArrayFormTable v-model:value="formData" />
</template>
```

## 在表单中使用

```vue
<script setup lang="ts">
import { useVbenForm } from '@vben/common-ui';
import ArrayFormTable from '@/components/array-form-table/index.vue';

const [Form] = useVbenForm({
  schema: [
    {
      fieldName: 'config',
      label: '配置项',
      component: ArrayFormTable,
    },
  ],
});
</script>
```

## Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| value | `string \| Record<string, string>` | `''` | 绑定值，支持 JSON 字符串或对象 |
| disabled | `boolean` | `false` | 是否禁用 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| update:value | `(value: string)` | 值更新时触发 |
| change | `(value: string)` | 值改变时触发 |

## 数据格式

### 输入格式

组件支持以下输入格式：

```typescript
// JSON 字符串
'{"key1":"value1","key2":"value2"}'

// 对象
{ key1: 'value1', key2: 'value2' }

// 数组格式
'[{"label":"key1","value":"value1"}]'
```

### 输出格式

组件始终输出 JSON 字符串格式：

```json
{"key1":"value1","key2":"value2"}
```

## 功能特性

### 拖拽排序

每行左侧有拖拽图标，可以通过拖拽调整行的顺序。

### 动态增删

- 点击底部「添加行」按钮添加新行
- 点击每行右侧的删除按钮删除该行

### 禁用状态

设置 `disabled` 属性后：
- 输入框变为只读
- 隐藏拖拽图标
- 隐藏添加按钮
- 删除按钮禁用

## 国际化

组件使用以下国际化 key：

```json
{
  "common.components.arrayFormTable.labelColumn": "键名",
  "common.components.arrayFormTable.valueColumn": "键值",
  "common.components.arrayFormTable.actionColumn": "操作",
  "common.components.arrayFormTable.addRow": "添加行"
}
```

## 使用场景

- 系统配置项编辑
- 环境变量配置
- 自定义字段配置
- 键值对数据管理

