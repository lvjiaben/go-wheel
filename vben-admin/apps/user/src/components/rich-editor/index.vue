<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { $t } from '@vben/locales';
import { QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';
import type Quill from 'quill';

import AttachmentSelector from '#/views/system/attachment/selector.vue';

interface Props {
  value?: string;
  placeholder?: string;
  minHeight?: number;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  value: '',
  placeholder: '',
  minHeight: 300,
  disabled: false,
});

const emit = defineEmits<{
  'update:value': [value: string];
  change: [value: string];
}>();

// 内容
const localValue = ref(props.value || '');

// 编辑器实例引用
const editorRef = ref<InstanceType<typeof QuillEditor>>();

// 附件选择器引用
const selectorRef = ref<InstanceType<typeof AttachmentSelector>>();

// 监听外部值变化
watch(
  () => props.value,
  (newValue) => {
    if (newValue !== localValue.value) {
      localValue.value = newValue || '';
    }
  },
);

// 内容改变
const handleUpdate = (content: string) => {
  localValue.value = content;
  emit('update:value', content);
  emit('change', content);
};

// 自定义图片上传处理器 - 打开附件管理弹窗
const handleImageUpload = () => {
  selectorRef.value?.open();
};

// 处理附件选择确认
const handleAttachmentConfirm = (urls: string[]) => {
  // 获取编辑器实例
  const quill = editorRef.value?.getQuill() as Quill;
  if (!quill) return;

  // 获取当前光标位置
  const range = quill.getSelection(true);

  // 插入所有选中的图片
  urls.forEach((url, index) => {
    quill.insertEmbed(range.index + index, 'image', url);
  });

  // 移动光标到最后一张图片后面
  quill.setSelection(range.index + urls.length);
};

// 编辑器准备完成
const handleReady = (quill: Quill) => {
  // 获取工具栏
  const toolbar = quill.getModule('toolbar') as any;

  // 重写图片按钮的处理函数
  if (toolbar && toolbar.addHandler) {
    toolbar.addHandler('image', handleImageUpload);
  }
};

// Quill 编辑器配置
const editorOptions = computed(() => ({
  theme: 'snow',
  modules: {
    toolbar: [
      ['bold', 'italic', 'underline', 'strike'],
      ['blockquote', 'code-block'],
      [{ header: 1 }, { header: 2 }],
      [{ list: 'ordered' }, { list: 'bullet' }],
      [{ script: 'sub' }, { script: 'super' }],
      [{ indent: '-1' }, { indent: '+1' }],
      [{ direction: 'rtl' }],
      [{ size: ['small', false, 'large', 'huge'] }],
      [{ header: [1, 2, 3, 4, 5, 6, false] }],
      [{ color: [] }, { background: [] }],
      [{ font: [] }],
      [{ align: [] }],
      ['clean'],
      ['link', 'image', 'video'],
    ],
  },
  placeholder: props.placeholder || $t('common.components.richEditor.placeholder'),
  readOnly: props.disabled,
}));
</script>

<template>
  <div>
    <div class="rich-editor-wrapper" :style="{ '--min-height': `${minHeight}px` }">
      <QuillEditor
        ref="editorRef"
        v-model:content="localValue"
        content-type="html"
        :options="editorOptions"
        @ready="handleReady"
        @update:content="handleUpdate"
      />
    </div>

    <!-- 附件选择器 -->
    <AttachmentSelector
      ref="selectorRef"
      :multiple="true"
      @confirm="handleAttachmentConfirm"
    />
  </div>
</template>

<style scoped>
.rich-editor-wrapper {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  overflow: hidden;
  transition: border-color 0.3s;
}

.rich-editor-wrapper:hover {
  border-color: #4096ff;
}

.rich-editor-wrapper:focus-within {
  border-color: #4096ff;
  box-shadow: 0 0 0 2px rgba(5, 145, 255, 0.1);
}

.rich-editor-wrapper :deep(.ql-toolbar) {
  border: none;
  border-bottom: 1px solid #d9d9d9;
  background-color: #fafafa;
}

.rich-editor-wrapper :deep(.ql-container) {
  border: none;
  min-height: var(--min-height);
  font-size: 14px;
}

.rich-editor-wrapper :deep(.ql-editor) {
  min-height: var(--min-height);
}

.rich-editor-wrapper :deep(.ql-editor.ql-blank::before) {
  color: #bfbfbf;
  font-style: normal;
}

/* 禁用状态 */
.rich-editor-wrapper :deep(.ql-toolbar.ql-disabled) {
  background-color: #f5f5f5;
  opacity: 0.6;
}

.rich-editor-wrapper :deep(.ql-editor.ql-disabled) {
  background-color: #f5f5f5;
  cursor: not-allowed;
}
</style>

