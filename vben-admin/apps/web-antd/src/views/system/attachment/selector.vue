<script lang="ts" setup>
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import AttachmentIndex from './index.vue';

const props = defineProps<{
  multiple?: boolean;
}>();

const emit = defineEmits<{
  confirm: [urls: string[]];
}>();

const attachmentRef = ref<InstanceType<typeof AttachmentIndex>>();

const [Modal, modalApi] = useVbenModal({
  class: 'w-full max-w-[90vw]',
  confirmText: $t('common.close'),
  onConfirm: () => {
    modalApi.close();
  },
  showCancelButton: false,
  title: $t('system.attachment.selectFile'),
});

const handleConfirm = (urls: string[]) => {
  emit('confirm', urls);
  modalApi.close();
};

defineExpose({
  close: () => modalApi.close(),
  open: () => modalApi.open(),
});
</script>

<template>
  <Modal>
    <div class="h-[80vh] md:h-[70vh]">
      <AttachmentIndex
        ref="attachmentRef"
        :multiple="props.multiple"
        selectable
        @confirm="handleConfirm"
      />
    </div>
  </Modal>
</template>

