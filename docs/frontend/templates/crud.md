# 前端 CRUD 模板

基于 Vben Admin 的标准 CRUD 页面模板。

## 目录结构

```
vben-admin/apps/backend/src/
├── api/
│   └── article.ts          # API 接口
└── views/
    └── article/
        ├── index.vue       # 列表页面
        ├── data.ts         # 表格和表单配置
        └── form.vue        # 表单弹窗
```

## 1. API 接口

```typescript
// api/article.ts
import { requestClient } from '#/api/request';

export interface ArticleRecord {
  id: number;
  title: string;
  content: string;
  status: number;
  sort: number;
  created_at: string;
}

// 列表
export function getArticleList(params?: object) {
  return requestClient.get('/article', { params });
}

// 创建
export function createArticle(data: Partial<ArticleRecord>) {
  return requestClient.post('/article', data);
}

// 更新
export function updateArticle(data: Partial<ArticleRecord>) {
  return requestClient.put('/article', data);
}

// 删除
export function deleteArticle(ids: number[]) {
  return requestClient.delete('/article', { data: { ids } });
}
```

## 2. 数据配置

```typescript
// views/article/data.ts
import type { VxeGridProps } from '#/adapter/vxe-table';
import type { VbenFormProps } from '#/adapter/form';

// 表格列配置
export const columns: VxeGridProps['columns'] = [
  { type: 'checkbox', width: 50 },
  { field: 'id', title: 'ID', width: 80 },
  { field: 'title', title: '标题', minWidth: 200 },
  { 
    field: 'status', 
    title: '状态', 
    width: 100,
    slots: { default: 'status' }
  },
  { field: 'sort', title: '排序', width: 80 },
  { field: 'created_at', title: '创建时间', width: 180 },
  { 
    field: 'action', 
    title: '操作', 
    width: 150, 
    fixed: 'right',
    slots: { default: 'action' }
  },
];

// 搜索表单配置
export const searchFormSchema: VbenFormProps['schema'] = [
  {
    component: 'Input',
    fieldName: 'search',
    label: '关键词',
    componentProps: {
      placeholder: '请输入标题',
    },
  },
  {
    component: 'Select',
    fieldName: 'filter[status]',
    label: '状态',
    componentProps: {
      options: [
        { label: '全部', value: '' },
        { label: '启用', value: 1 },
        { label: '禁用', value: 0 },
      ],
    },
  },
];

// 表单配置
export const formSchema: VbenFormProps['schema'] = [
  {
    component: 'Input',
    fieldName: 'title',
    label: '标题',
    rules: 'required',
    componentProps: {
      placeholder: '请输入标题',
    },
  },
  {
    component: 'RichEditor',  // 自定义组件，无需 import
    fieldName: 'content',
    label: '内容',
    rules: 'required',
  },
  {
    component: 'RadioGroup',
    fieldName: 'status',
    label: '状态',
    defaultValue: 1,
    componentProps: {
      options: [
        { label: '启用', value: 1 },
        { label: '禁用', value: 0 },
      ],
    },
  },
  {
    component: 'InputNumber',
    fieldName: 'sort',
    label: '排序',
    defaultValue: 0,
  },
];
```

## 3. 列表页面

```vue
<!-- views/article/index.vue -->
<script setup lang="ts">
import { ref } from 'vue';
import { Page, useVbenModal } from '@vben/common-ui';
import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { Button, Tag, message, Popconfirm } from 'ant-design-vue';
import { getArticleList, deleteArticle } from '#/api/article';
import { columns, searchFormSchema } from './data';
import Form from './form.vue';

const [FormModal, formModalApi] = useVbenModal({
  connectedComponent: Form,
});

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns,
    proxyConfig: {
      ajax: {
        query: async ({ page }) => {
          const params = {
            page: page.currentPage,
            page_size: page.pageSize,
            ...gridApi.formApi.form.values,
          };
          return await getArticleList(params);
        },
      },
    },
    pagerConfig: {},
    toolbarConfig: {
      search: true,
    },
  },
  formOptions: {
    schema: searchFormSchema,
    submitOnChange: true,
  },
});

// 新增
const handleAdd = () => {
  formModalApi.setData({ mode: 'create' });
  formModalApi.open();
};

// 编辑
const handleEdit = (row: any) => {
  formModalApi.setData({ mode: 'edit', record: row });
  formModalApi.open();
};

// 删除
const handleDelete = async (ids: number[]) => {
  await deleteArticle(ids);
  message.success('删除成功');
  gridApi.reload();
};

// 表单提交成功回调
const handleSuccess = () => {
  gridApi.reload();
};
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #toolbar-tools>
        <Button type="primary" @click="handleAdd">新增</Button>
      </template>

      <template #status="{ row }">
        <Tag :color="row.status === 1 ? 'green' : 'red'">
          {{ row.status === 1 ? '启用' : '禁用' }}
        </Tag>
      </template>

      <template #action="{ row }">
        <Button type="link" size="small" @click="handleEdit(row)">编辑</Button>
        <Popconfirm title="确定删除？" @confirm="handleDelete([row.id])">
          <Button type="link" size="small" danger>删除</Button>
        </Popconfirm>
      </template>
    </Grid>

    <FormModal @success="handleSuccess" />
  </Page>
</template>
```

## 4. 表单弹窗

```vue
<!-- views/article/form.vue -->
<script setup lang="ts">
import { ref, computed } from 'vue';
import { useVbenModal } from '@vben/common-ui';
import { useVbenForm } from '#/adapter/form';
import { message } from 'ant-design-vue';
import { createArticle, updateArticle } from '#/api/article';
import { formSchema } from './data';

const emit = defineEmits(['success']);

const mode = ref<'create' | 'edit'>('create');
const recordId = ref<number>();

const [Form, formApi] = useVbenForm({
  schema: formSchema,
  showDefaultActions: false,
});

const [Modal, modalApi] = useVbenModal({
  onOpenChange: async (isOpen) => {
    if (isOpen) {
      const data = modalApi.getData<{ mode: string; record?: any }>();
      mode.value = data?.mode as 'create' | 'edit';

      if (data?.record) {
        recordId.value = data.record.id;
        await formApi.setValues(data.record);
      } else {
        recordId.value = undefined;
        await formApi.resetForm();
      }
    }
  },
  onConfirm: async () => {
    const values = await formApi.getValues();

    if (mode.value === 'create') {
      await createArticle(values);
      message.success('创建成功');
    } else {
      await updateArticle({ ...values, id: recordId.value });
      message.success('更新成功');
    }

    emit('success');
    modalApi.close();
  },
});

const title = computed(() => mode.value === 'create' ? '新增文章' : '编辑文章');
</script>

<template>
  <Modal :title="title">
    <Form />
  </Modal>
</template>
```

## 5. 使用自定义组件

在表单中直接使用组件名，无需 import：

```typescript
// 富文本编辑器
{ component: 'RichEditor', fieldName: 'content', label: '内容' }

// 附件选择
{ component: 'AttachmentInput', fieldName: 'image', label: '图片' }

// 表格选择器
{
  component: 'TableSelect',
  fieldName: 'user_id',
  label: '用户',
  componentProps: {
    api: getUserList,
    columns: [{ field: 'id', title: 'ID' }, { field: 'username', title: '用户名' }],
  }
}

// 数组表单
{ component: 'ArrayFormTable', fieldName: 'items', label: '配置项' }
```
```

