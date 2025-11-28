<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { IconifyIcon } from '@vben/icons';
import { $t } from '@vben/locales';

import { Button, Input } from 'ant-design-vue';

import AttachmentSelector from '#/views/system/attachment/selector.vue';

const props = withDefaults(
  defineProps<{
    value?: string;
    modelValue?: string;
    multiple?: boolean;
    placeholder?: string;
    showPreview?: boolean;
    showInput?: boolean;
  }>(),
  {
    showPreview: true,
    showInput: true,
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: string];
  'update:value': [value: string];
  change: [value: string];
}>();

const selectorRef = ref<InstanceType<typeof AttachmentSelector>>();

// 支持 v-model 和 v-model:value
const localValue = computed({
  get: () => props.modelValue ?? props.value ?? '',
  set: (val: string) => {
    emit('update:modelValue', val);
    emit('update:value', val);
    emit('change', val);
  },
});

const openSelector = () => {
  selectorRef.value?.open();
};

const handleConfirm = (urls: string[]) => {
  if (props.multiple) {
    // 多选：英文逗号分隔
    localValue.value = urls.join(',');
  } else {
    // 单选：直接赋值
    localValue.value = urls[0] || '';
  }
};

// 监听外部值变化
watch(
  () => props.value,
  (val) => {
    if (val !== undefined) {
      localValue.value = val;
    }
  },
);

// 计算图片URL列表
const imageUrls = computed(() => {
  const value = localValue.value;
  if (!value) return [];
  return value.split(',').filter((url) => url.trim());
});

// 判断是否是图片
const isImageUrl = (url: string) => {
  const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg', '.bmp'];
  return imageExts.some((ext) => url.toLowerCase().includes(ext));
};

// 移除图片
const removeImage = (index: number) => {
  const urls = imageUrls.value.filter((_, i) => i !== index);
  localValue.value = urls.join(',');
};

// 在新窗口打开
const openInNewTab = (url: string) => {
  window.open(url, '_blank');
};
</script>

<template>
  <div>
    <!-- 输入框和按钮 -->
    <div v-if="showInput" class="flex items-center gap-2">
      <Input
        v-model:value="localValue"
        :placeholder="placeholder || $t('system.attachment.inputPlaceholder')"
      />
      <Button @click="openSelector">
        <IconifyIcon icon="mdi:upload" class="size-4" />
        {{ $t('system.attachment.selectFile') }}
      </Button>
    </div>

    <!-- 只显示按钮，不显示输入框 -->
    <div v-else class="flex items-center">
      <Button @click="openSelector">
        <IconifyIcon icon="mdi:upload" class="size-4" />
        {{ $t('system.attachment.selectFile') }}
      </Button>
    </div>

    <!-- 图片预览区域 -->
    <div v-if="showPreview && imageUrls.length > 0" class="mt-2 flex flex-wrap gap-2">
      <div
        v-for="(url, index) in imageUrls"
        :key="index"
        class="group relative h-16 w-16 flex-shrink-0 overflow-hidden rounded border bg-gray-50"
      >
        <!-- 图片预览 -->
        <img
          v-if="isImageUrl(url)"
          :alt="`preview-${index}`"
          class="h-full w-full cursor-pointer"
          :src="url"
          style="object-fit: cover;"
          @click="() => openInNewTab(url)"
        />
        <!-- 非图片文件显示文件图标 -->
        <div
          v-else
          class="flex h-full w-full cursor-pointer items-center justify-center text-xs font-medium text-gray-600"
          @click="openInNewTab(url)"
        >
          {{ url.split('.').pop()?.toUpperCase() || 'FILE' }}
        </div>

        <!-- 删除按钮 -->
        <div
          class="absolute right-0.5 top-0.5 hidden cursor-pointer rounded-full bg-red-500 p-0.5 text-white shadow-sm transition-all hover:bg-red-600 group-hover:block"
          @click.stop="removeImage(index)"
        >
          <IconifyIcon icon="mdi:close" class="size-3" />
        </div>
      </div>
    </div>

    <AttachmentSelector
      ref="selectorRef"
      :multiple="multiple"
      @confirm="handleConfirm"
    />
  </div>
</template>

