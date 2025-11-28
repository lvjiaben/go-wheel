<script lang="ts" setup>
import type { AttachmentApi } from '#/api/system/attachment';

import { computed, onMounted, ref } from 'vue';

import { Checkbox } from 'ant-design-vue';

const props = defineProps<{
  file: AttachmentApi.Attachment;
  selected?: boolean;
  multiple?: boolean;
}>();

const emit = defineEmits<{
  select: [file: AttachmentApi.Attachment];
  dblclick: [file: AttachmentApi.Attachment];
}>();

const canvasRef = ref<HTMLCanvasElement>();

// 判断是否是图片
const isImage = computed(() => {
  return props.file.mediatype?.startsWith('image/');
});

// 获取文件扩展名
const fileExt = computed(() => {
  const ext = props.file.extension?.toUpperCase() || 'FILE';
  return ext.length > 4 ? ext.slice(0, 4) : ext;
});

// 绘制文件后缀
const drawFileExt = () => {
  if (isImage.value || !canvasRef.value) return;

  const canvas = canvasRef.value;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  // 设置画布大小
  canvas.width = 200;
  canvas.height = 200;

  // 背景渐变
  const gradient = ctx.createLinearGradient(0, 0, 200, 200);
  gradient.addColorStop(0, '#667eea');
  gradient.addColorStop(1, '#764ba2');
  ctx.fillStyle = gradient;
  ctx.fillRect(0, 0, 200, 200);

  // 绘制文字
  ctx.fillStyle = '#ffffff';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';

  const ext = fileExt.value;
  const fontSize = ext.length <= 3 ? 48 : 36;
  ctx.font = `bold ${fontSize}px Arial`;
  ctx.fillText(ext, 100, 100);
};

onMounted(() => {
  drawFileExt();
});

const handleClick = () => {
  emit('select', props.file);
};

const handleDblClick = () => {
  emit('dblclick', props.file);
};
</script>

<template>
  <div
    class="file-card group relative cursor-pointer overflow-hidden rounded-lg border-2 transition-all hover:shadow-lg"
    :class="{
      'border-primary bg-primary/5': selected,
      'border-gray-200 hover:border-primary/50': !selected,
    }"
    @click="handleClick"
    @dblclick="handleDblClick"
  >
    <!-- 选中标记 -->
    <div
      v-if="multiple"
      class="absolute right-2 top-2 z-10"
      @click.stop="handleClick"
    >
      <Checkbox :checked="selected" />
    </div>

    <!-- 缩略图区域 -->
    <div class="flex aspect-square items-center justify-center bg-gray-50">
      <!-- 图片缩略图 -->
      <img
        v-if="isImage"
        :alt="file.filename"
        class="h-full w-full object-cover"
        :src="file.url"
      />
      <!-- 文件后缀绘制 -->
      <canvas v-else ref="canvasRef" class="h-full w-full"></canvas>
    </div>

    <!-- 文件信息 -->
    <div class="p-3">
      <div class="truncate text-sm font-medium" :title="file.filename">
        {{ file.filename }}
      </div>
      <div class="mt-1 flex items-center justify-between text-xs text-gray-500">
        <span>{{ (file.size / 1024).toFixed(2) }} KB</span>
        <span>{{ new Date(file.created_at * 1000).toLocaleDateString() }}</span>
      </div>
    </div>

    <!-- 选中遮罩 -->
    <div
      v-if="selected"
      class="pointer-events-none absolute inset-0 bg-primary/10"
    ></div>
  </div>
</template>

<style scoped>
.file-card {
  @apply transition-all duration-200;
}
</style>

