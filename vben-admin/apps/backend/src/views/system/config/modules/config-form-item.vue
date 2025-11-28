<script lang="ts" setup>
import type { ConfigApi } from '#/api/system/config';

import { computed } from 'vue';

import { IconifyIcon } from '@vben/icons';

import {
  Button,
  Checkbox,
  DatePicker,
  Input,
  Radio,
  Select,
  Switch,
  TimePicker,
} from 'ant-design-vue';

import ArrayFormTable from '#/components/array-form-table/index.vue';
import AttachmentInput from '#/components/attachment-input/index.vue';
import RichEditor from '#/components/rich-editor/index.vue';
import TableSelect from '#/components/table-select/index.vue';

const props = defineProps<{
  config: ConfigApi.Config;
  value: any;
}>();

const emit = defineEmits<{
  'update:value': [value: any];
  delete: [];
}>();

// 本地值
const localValue = computed({
  get: () => props.value,
  set: (val: any) => {
    emit('update:value', val);
  },
});

// 解析 variable 字段为选项数组
const parseOptions = (variable: string) => {
  if (!variable) return [];
  try {
    const parsed = JSON.parse(variable);

    // 如果是数组格式
    if (Array.isArray(parsed)) {
      return parsed.map((item) => {
        // 如果数组元素已经是对象格式 {label, value}
        if (typeof item === 'object' && item !== null && 'label' in item && 'value' in item) {
          return item;
        }
        // 如果是简单字符串，转换为 {label, value} 格式
        return {
          label: String(item),
          value: String(item),
        };
      });
    }

    // 如果是对象格式 {key: label}
    if (typeof parsed === 'object' && parsed !== null) {
      return Object.entries(parsed).map(([key, label]) => ({
        label: label as string,
        value: key,
      }));
    }
  } catch {
    // 解析失败，返回空数组
  }
  return [];
};

// Switch 值转换 (1/0 <-> boolean)
const switchValue = computed({
  get: () => props.value === '1' || props.value === 1 || props.value === true,
  set: (val: boolean) => {
    emit('update:value', val ? '1' : '0');
  },
});

// Checkbox 多选值处理
const checkboxValue = computed({
  get: () => {
    if (!props.value) return [];
    if (typeof props.value === 'string') {
      return props.value.split(',').filter(Boolean);
    }
    return props.value;
  },
  set: (val: string[]) => {
    emit('update:value', val.join(','));
  },
});

// Selects 多选值处理
const selectsValue = computed({
  get: () => {
    if (!props.value) return [];
    if (typeof props.value === 'string') {
      // 如果是空字符串，返回空数组
      if (props.value === '') return [];
      return props.value.split(',').filter(Boolean);
    }
    if (Array.isArray(props.value)) {
      return props.value;
    }
    return [];
  },
  set: (val: string[]) => {
    emit('update:value', val.join(','));
  },
});

// DatetimeRange 值处理
const datetimeRangeValue = computed({
  get: () => {
    if (!props.value) return undefined;
    if (typeof props.value === 'string') {
      // 如果是空字符串，返回 undefined
      if (props.value === '') return undefined;
      // 尝试解析为数组
      try {
        const parsed = JSON.parse(props.value);
        if (Array.isArray(parsed)) {
          return parsed;
        }
      } catch {
        // 解析失败，返回 undefined
      }
    }
    if (Array.isArray(props.value)) {
      return props.value;
    }
    return undefined;
  },
  set: (val: any) => {
    if (Array.isArray(val)) {
      emit('update:value', JSON.stringify(val));
    } else {
      emit('update:value', '');
    }
  },
});

// TableSelect 值处理（确保是字符串或数组）
const tableSelectValue = computed({
  get: () => {
    // 如果是数字，转换为字符串
    if (typeof props.value === 'number') {
      return String(props.value);
    }
    return props.value;
  },
  set: (val: any) => {
    emit('update:value', val);
  },
});

// TableSelects 多选值处理
const tableSelectsValue = computed({
  get: () => {
    if (!props.value) return [];
    if (typeof props.value === 'string') {
      if (props.value === '') return [];
      // 尝试解析 JSON 数组
      try {
        const parsed = JSON.parse(props.value);
        if (Array.isArray(parsed)) {
          return parsed.map(String); // 确保都是字符串
        }
      } catch {
        // 如果不是 JSON，按逗号分隔
        return props.value.split(',').filter(Boolean);
      }
    }
    if (Array.isArray(props.value)) {
      return props.value.map(String); // 确保都是字符串
    }
    return [];
  },
  set: (val: any[]) => {
    if (Array.isArray(val)) {
      emit('update:value', JSON.stringify(val));
    } else {
      emit('update:value', '');
    }
  },
});

// 获取选项列表
const options = computed(() => parseOptions(props.config.variable));
</script>

<template>
  <div class="flex items-start gap-2">
    <div class="flex-1">
      <div class="mb-1 font-medium">{{ config.name }}</div>

      <!-- array: 数组表单 -->
      <ArrayFormTable
        v-if="config.type === 'array'"
        v-model:value="localValue"
      />

      <!-- input: 普通输入框 -->
      <div v-else-if="config.type === 'input'" class="flex items-center gap-2">
        <Input v-model:value="localValue" class="flex-1" />
      </div>

      <!-- text: 文本域 -->
      <div v-else-if="config.type === 'text'" class="flex items-center gap-2">
        <Input.TextArea v-model:value="localValue" :rows="4" class="flex-1" />
      </div>

      <!-- file: 文件上传 -->
      <div v-else-if="config.type === 'file'" class="flex items-center gap-2">
        <AttachmentInput v-model:value="localValue" :multiple="false" class="flex-1" />
      </div>

      <!-- image: 单图上传 -->
      <div v-else-if="config.type === 'image'" class="flex items-center gap-2">
        <AttachmentInput v-model:value="localValue" :multiple="false" class="flex-1" />
      </div>

      <!-- images: 多图上传 -->
      <div v-else-if="config.type === 'images'" class="flex items-center gap-2">
        <AttachmentInput v-model:value="localValue" :multiple="true" class="flex-1" />
      </div>

      <!-- editor: 富文本编辑器 -->
      <div v-else-if="config.type === 'editor'" class="flex items-center gap-2">
        <RichEditor v-model:value="localValue" class="flex-1" />
      </div>

      <!-- date: 日期选择 -->
      <div v-else-if="config.type === 'date'" class="flex items-center gap-2">
        <DatePicker
          v-model:value="localValue"
          class="flex-1"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
        />
      </div>

      <!-- time: 时间选择 -->
      <div v-else-if="config.type === 'time'" class="flex items-center gap-2">
        <TimePicker
          v-model:value="localValue"
          class="flex-1"
          format="HH:mm:ss"
          value-format="HH:mm:ss"
        />
      </div>

      <!-- datetime: 日期时间选择 -->
      <div v-else-if="config.type === 'datetime'" class="flex items-center gap-2">
        <DatePicker
          v-model:value="localValue"
          class="flex-1"
          format="YYYY-MM-DD HH:mm:ss"
          show-time
          value-format="YYYY-MM-DD HH:mm:ss"
        />
      </div>

      <!-- datetimerange: 日期时间范围 -->
      <div v-else-if="config.type === 'datetimerange'" class="flex items-center gap-2">
        <DatePicker.RangePicker
          v-model:value="datetimeRangeValue"
          class="flex-1"
          format="YYYY-MM-DD HH:mm:ss"
          show-time
          value-format="YYYY-MM-DD HH:mm:ss"
        />
      </div>

      <!-- select: 单选下拉框 -->
      <div v-else-if="config.type === 'select'" class="flex items-center gap-2">
        <Select v-model:value="localValue" :options="options" class="flex-1" />
      </div>

      <!-- selects: 多选下拉框 -->
      <div v-else-if="config.type === 'selects'" class="flex items-center gap-2">
        <Select v-model:value="selectsValue" :options="options" class="flex-1" mode="multiple" />
      </div>

      <!-- switch: 开关 -->
      <div v-else-if="config.type === 'switch'" class="flex items-center gap-2">
        <Switch v-model:checked="switchValue" />
      </div>

      <!-- checkbox: 多选框 -->
      <div v-else-if="config.type === 'checkbox'" class="flex items-center gap-2">
        <Checkbox.Group v-model:value="checkboxValue" :options="options" />
      </div>

      <!-- radio: 单选框 -->
      <div v-else-if="config.type === 'radio'" class="flex items-center gap-2">
        <Radio.Group v-model:value="localValue" :options="options" />
      </div>

      <!-- tableselect: 单选远程搜索下拉框 -->
      <div v-else-if="config.type === 'tableselect'" class="flex items-center gap-2">
        <TableSelect v-model:value="tableSelectValue" :config="config.variable" :multiple="false" class="flex-1" />
      </div>

      <!-- tableselects: 多选远程搜索下拉框 -->
      <div v-else-if="config.type === 'tableselects'" class="flex items-center gap-2">
        <TableSelect v-model:value="tableSelectsValue" :config="config.variable" :multiple="true" class="flex-1" />
      </div>

      <!-- 默认: 普通输入框 -->
      <div v-else class="flex items-center gap-2">
        <Input v-model:value="localValue" class="flex-1" />
      </div>

      <!-- 提示信息 -->
      <div v-if="config.tip" class="mt-1 text-xs text-gray-500">
        {{ config.tip }}
      </div>
    </div>

    <!-- 删除按钮 -->
    <Button size="small" type="text" class="delete-btn" @click="emit('delete')">
      <IconifyIcon icon="mdi:close" class="size-3" />
    </Button>
  </div>
</template>

<style scoped>
.delete-btn {
  color: #999;
  opacity: 0.6;
  transition: all 0.2s;
}

.delete-btn:hover {
  color: #ff4d4f;
  opacity: 1;
}
</style>
