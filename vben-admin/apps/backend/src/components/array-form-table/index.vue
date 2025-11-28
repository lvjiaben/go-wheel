<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { IconifyIcon } from '@vben/icons';
import { $t } from '@vben/locales';

import { Button, Input, Table } from 'ant-design-vue';

interface Props {
  value?: string | Record<string, string>;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  value: '',
  disabled: false,
});

const emit = defineEmits<{
  'update:value': [value: string];
  change: [value: string];
}>();

interface ArrayItem {
  label: string;
  value: string;
}

// 解析值
const parseValue = (val: string | Record<string, string>): ArrayItem[] => {
  if (!val) return [];
  
  // 如果已经是对象
  if (typeof val === 'object' && !Array.isArray(val)) {
    return Object.entries(val).map(([key, value]) => ({
      label: key,
      value: value as string,
    }));
  }
  
  // 如果是字符串，尝试解析
  if (typeof val === 'string') {
    try {
      const parsed = JSON.parse(val);
      if (Array.isArray(parsed)) {
        return parsed;
      }
      if (typeof parsed === 'object') {
        return Object.entries(parsed).map(([key, value]) => ({
          label: key,
          value: value as string,
        }));
      }
    } catch {
      // 解析失败
    }
  }
  
  return [];
};

// 序列化值
const stringifyValue = (items: ArrayItem[]): string => {
  const obj: Record<string, string> = {};
  items.forEach((item) => {
    if (item.label) {
      obj[item.label] = item.value;
    }
  });
  return JSON.stringify(obj);
};

// 本地数据
const localData = ref<ArrayItem[]>(parseValue(props.value));

// 监听外部值变化
watch(
  () => props.value,
  (newValue) => {
    localData.value = parseValue(newValue);
  },
);

// 更新值
const updateValue = () => {
  const stringValue = stringifyValue(localData.value);
  emit('update:value', stringValue);
  emit('change', stringValue);
};

// 添加行
const handleAdd = () => {
  localData.value.push({ label: '', value: '' });
  updateValue();
};

// 删除行
const handleDelete = (index: number) => {
  localData.value.splice(index, 1);
  updateValue();
};

// 更新行数据
const handleUpdate = (index: number, field: 'label' | 'value', val: string) => {
  if (localData.value[index]) {
    localData.value[index][field] = val;
    updateValue();
  }
};

// 拖拽排序相关
const dragIndex = ref<number | null>(null);

const handleDragStart = (index: number) => {
  dragIndex.value = index;
};

const handleDragOver = (e: DragEvent) => {
  e.preventDefault();
};

const handleDrop = (index: number) => {
  if (dragIndex.value === null || dragIndex.value === index) return;

  const dragItem = localData.value[dragIndex.value];
  if (!dragItem) return;

  localData.value.splice(dragIndex.value, 1);
  localData.value.splice(index, 0, dragItem);

  dragIndex.value = null;
  updateValue();
};

// 表格列定义
const columns = computed(() => [
  {
    title: $t('common.components.arrayFormTable.labelColumn'),
    dataIndex: 'label',
    width: '40%',
  },
  {
    title: $t('common.components.arrayFormTable.valueColumn'),
    dataIndex: 'value',
    width: '40%',
  },
  {
    title: $t('common.components.arrayFormTable.actionColumn'),
    dataIndex: 'action',
    width: '20%',
  },
]);
</script>

<template>
  <div class="array-form-table">
    <Table
      :columns="columns"
      :data-source="localData"
      :pagination="false"
      bordered
      size="small"
    >
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.dataIndex === 'label'">
          <div
            class="flex items-center gap-2"
            :draggable="!disabled"
            @dragover="handleDragOver"
            @dragstart="handleDragStart(index)"
            @drop="handleDrop(index)"
          >
            <IconifyIcon
              v-if="!disabled"
              icon="mdi:drag"
              class="size-4 cursor-move text-gray-400"
            />
            <Input
              :value="record.label"
              :disabled="disabled"
              size="small"
              @update:value="(val) => handleUpdate(index, 'label', val)"
            />
          </div>
        </template>
        <template v-else-if="column.dataIndex === 'value'">
          <Input
            :value="record.value"
            :disabled="disabled"
            size="small"
            @update:value="(val) => handleUpdate(index, 'value', val)"
          />
        </template>
        <template v-else-if="column.dataIndex === 'action'">
          <Button
            danger
            :disabled="disabled"
            size="small"
            type="link"
            @click="handleDelete(index)"
          >
            <IconifyIcon icon="mdi:delete" class="size-4" />
          </Button>
        </template>
      </template>
    </Table>

    <Button
      v-if="!disabled"
      class="mt-2"
      size="small"
      type="dashed"
      @click="handleAdd"
    >
      <IconifyIcon icon="mdi:plus" class="mr-1 size-4" />
      {{ $t('common.components.arrayFormTable.addRow') }}
    </Button>
  </div>
</template>

<style scoped>
.array-form-table :deep(.ant-table) {
  font-size: 12px;
}
</style>

