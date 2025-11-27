<script lang="ts" setup>
import type { AttachmentApi } from '#/api/system/attachment';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';
import { $t } from '@vben/locales';

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import {
  Button,
  Empty,
  Input,
  message,
  Modal,
  Pagination,
  Spin,
  Upload,
} from 'ant-design-vue';

import {
  deleteAttachment,
  getAttachmentList,
  getDirectories,
  uploadAttachment,
} from '#/api/system/attachment';

import FileCard from './file-card.vue';

const props = withDefaults(
  defineProps<{
    multiple?: boolean;
    selectable?: boolean;
  }>(),
  {
    multiple: true, // 默认支持多选
    selectable: false,
  },
);

const emit = defineEmits<{
  confirm: [urls: string[]];
}>();

// 响应式断点
const breakpoints = useBreakpoints(breakpointsTailwind);

// 目录列表
const directories = ref<AttachmentApi.Directory[]>([]);
const currentDirectory = ref<string>('');
const loadingDirectories = ref(false);

// 文件列表
const files = ref<AttachmentApi.Attachment[]>([]);
const selectedIds = ref<number[]>([]);
const loading = ref(false);
const keyword = ref('');

// 分页
const pagination = ref({
  current: 1,
  pageSize: 20,
  total: 0,
});

// 计算选中的文件
const selectedFiles = computed(() => {
  return files.value.filter((f) => selectedIds.value.includes(f.id));
});

// 加载目录列表
const loadDirectories = async () => {
  loadingDirectories.value = true;
  try {
    const result = await getDirectories();
    directories.value = result;
    // 默认选中第一个目录（后端返回的"全部"）
    if (result.length > 0 && !currentDirectory.value && result[0]) {
      currentDirectory.value = result[0].path;
    }
  } finally {
    loadingDirectories.value = false;
  }
};

// 加载文件列表
const loadFiles = async () => {
  loading.value = true;
  try {
    const result = await getAttachmentList({
      parent: currentDirectory.value || undefined,
      search: keyword.value || undefined,
      page: pagination.value.current,
      page_size: pagination.value.pageSize,
    });
    files.value = result.list || [];
    pagination.value.total = result.total || 0;
  } finally {
    loading.value = false;
  }
};

// 选择目录
const selectDirectory = (dir: string) => {
  currentDirectory.value = dir;
  pagination.value.current = 1;
  selectedIds.value = [];
  loadFiles();
};

// 搜索
const handleSearch = () => {
  pagination.value.current = 1;
  selectedIds.value = [];
  loadFiles();
};

// 清空搜索时自动搜索
const handleSearchOnClear = () => {
  // 当清空输入框时（keyword为空），自动触发搜索
  if (!keyword.value || keyword.value.trim() === '') {
    handleSearch();
  }
};

// 选择文件
const handleFileSelect = (file: AttachmentApi.Attachment) => {
  const index = selectedIds.value.indexOf(file.id);
  if (props.multiple) {
    if (index > -1) {
      selectedIds.value.splice(index, 1);
    } else {
      selectedIds.value.push(file.id);
    }
  } else {
    selectedIds.value = index > -1 ? [] : [file.id];
  }
};

// 双击文件（单选模式下直接确认）
const handleFileDblClick = (file: AttachmentApi.Attachment) => {
  if (!props.multiple && props.selectable) {
    emit('confirm', [file.url]);
  }
};

// 上传文件
const handleUpload = async (options: any) => {
  const { file } = options;
  const hideLoading = message.loading($t('system.attachment.uploading'), 0);
  try {
    const result = await uploadAttachment(
      file,
      currentDirectory.value || undefined,
    );
    message.success($t('system.attachment.uploadSuccess'));
    hideLoading();
    // 刷新列表
    await loadFiles();
    // 自动选中
    if (props.selectable) {
      selectedIds.value = [result.id];
    }
  } catch {
    hideLoading();
  }
};

// 删除文件
const handleDelete = () => {
  if (selectedIds.value.length === 0) {
    message.warning($t('system.attachment.selectFilesFirst'));
    return;
  }

  Modal.confirm({
    content: $t('system.attachment.deleteConfirm', [selectedIds.value.length]),
    onOk: async () => {
      const hideLoading = message.loading($t('system.attachment.deleting'), 0);
      try {
        await deleteAttachment(selectedIds.value);
        message.success($t('system.attachment.deleteSuccess'));
        selectedIds.value = [];
        await loadFiles();
      } finally {
        hideLoading();
      }
    },
    title: $t('common.confirm'),
  });
};

// 读取信息
const [InfoModal, infoModalApi] = useVbenModal({
  class: 'w-full max-w-2xl',
  confirmText: $t('common.close'),
  onConfirm: () => {
    infoModalApi.close();
  },
  showCancelButton: false,
  title: $t('system.attachment.fileInfo'),
});

const currentFileInfo = ref<AttachmentApi.Attachment | null>(null);

const handleShowInfo = () => {
  if (selectedIds.value.length === 0) {
    message.warning($t('system.attachment.selectFileFirst'));
    return;
  }
  if (selectedIds.value.length > 1) {
    message.warning($t('system.attachment.selectOnlyOneFile'));
    return;
  }
  currentFileInfo.value =
    files.value.find((f) => f.id === selectedIds.value[0]) || null;
  infoModalApi.open();
};

// 分页改变
const handlePageChange = (page: number, pageSize: number) => {
  pagination.value.current = page;
  pagination.value.pageSize = pageSize;
  loadFiles();
};

// 确认选择
const handleConfirm = () => {
  if (selectedIds.value.length === 0) {
    message.warning($t('system.attachment.selectFilesFirst'));
    return;
  }
  const urls = selectedFiles.value.map((f) => f.url);
  emit('confirm', urls);
};

onMounted(() => {
  loadDirectories();
  loadFiles();
});

defineExpose({
  refresh: loadFiles,
});
</script>

<template>
  <Page auto-content-height>
    <div class="flex flex-col gap-4 md:h-full md:flex-row">
      <!-- 目录树 - 移动端在上，桌面端在左 -->
      <!-- 以弹窗形式打开时，移动端隐藏目录分类 -->
      <div
        v-if="!selectable || breakpoints.greaterOrEqual('md').value"
        class="w-full flex-shrink-0 overflow-hidden rounded-lg border bg-white md:w-48"
      >
        <div class="border-b p-3 font-semibold">
          {{ $t('system.attachment.directories') }}
        </div>
        <Spin :spinning="loadingDirectories">
          <div class="max-h-40 overflow-y-auto p-2 md:h-[calc(100vh-200px)]">
            <div
              v-for="(dir, index) in directories"
              :key="index"
              class="cursor-pointer rounded px-3 py-2 transition-colors hover:bg-gray-100"
              :class="{
                'bg-primary text-white hover:bg-primary':
                  currentDirectory === dir.path,
              }"
              @click="selectDirectory(dir.path)"
            >
              <div class="flex items-center justify-between">
                <span class="truncate">
                  {{ dir.name }}
                </span>
                <span
                  v-if="dir.count > 0"
                  class="text-xs"
                  :class="{
                    'text-gray-400': currentDirectory !== dir.path,
                  }"
                >
                  {{ dir.count }}
                </span>
              </div>
            </div>
          </div>
        </Spin>
      </div>

      <!-- 文件列表 - 移动端在下，桌面端在右 -->
      <div class="flex-1 overflow-hidden rounded-lg border bg-white">
        <!-- 工具栏 -->
        <div class="flex flex-wrap items-center gap-2 border-b p-3 md:gap-3">
          <div class="flex w-full gap-2 md:w-auto">
            <Input
              v-model:value="keyword"
              allow-clear
              class="flex-1 md:w-64"
              :placeholder="$t('system.attachment.searchPlaceholder')"
              @change="handleSearchOnClear"
              @press-enter="handleSearch"
            />
            <Button @click="handleSearch">
              {{ $t('common.search') }}
            </Button>
          </div>

          <div class="hidden flex-1 md:block"></div>

          <Upload
            :custom-request="handleUpload"
            :show-upload-list="false"
          >
            <Button type="primary">
              <IconifyIcon icon="mdi:upload" class="size-4" />
              {{ $t('system.attachment.upload') }}
            </Button>
          </Upload>

          <Button danger @click="handleDelete">
            <IconifyIcon icon="mdi:delete" class="size-4" />
            {{ $t('system.attachment.delete') }}
          </Button>

          <Button :disabled="selectedIds.length !== 1" @click="handleShowInfo">
            <IconifyIcon icon="mdi:information" class="size-4" />
            {{ $t('system.attachment.info') }}
          </Button>

          <Button
            v-if="selectable && selectedIds.length > 0"
            type="primary"
            class="bg-green-600 hover:bg-green-700"
            @click="handleConfirm"
          >
            {{ $t('common.confirm') }}
            <span>({{ selectedIds.length }})</span>
          </Button>
        </div>

        <!-- 文件网格 -->
        <Spin :spinning="loading">
          <div class="min-h-[300px] overflow-y-auto p-4 md:h-[calc(100vh-280px)]">
            <div
              v-if="files.length > 0"
              class="grid grid-cols-2 gap-3 sm:grid-cols-3 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6"
            >
              <FileCard
                v-for="file in files"
                :key="file.id"
                :file="file"
                :multiple="multiple"
                :selected="selectedIds.includes(file.id)"
                @dblclick="handleFileDblClick"
                @select="handleFileSelect"
              />
            </div>
            <Empty v-else :description="$t('system.attachment.noFiles')" />
          </div>
        </Spin>

        <!-- 分页 -->
        <div class="flex justify-center border-t p-2 md:p-3">
          <Pagination
            v-model:current="pagination.current"
            v-model:page-size="pagination.pageSize"
            :page-size-options="['20', '40', '60', '100']"
            show-size-changer
            :show-total="(total) => `共 ${total} 项`"
            :simple="false"
            :total="pagination.total"
            @change="handlePageChange"
          />
        </div>
      </div>
    </div>

    <!-- 文件信息 Modal -->
    <InfoModal>
      <div v-if="currentFileInfo" class="space-y-4">
        <div class="flex flex-col items-center gap-4 sm:flex-row sm:items-start">
          <div class="h-32 w-32 flex-shrink-0 overflow-hidden rounded border">
            <img
              v-if="currentFileInfo.mediatype?.startsWith('image/')"
              :alt="currentFileInfo.filename"
              class="h-full w-full object-cover"
              :src="currentFileInfo.url"
            />
            <div
              v-else
              class="flex h-full w-full items-center justify-center bg-gray-100 text-2xl font-bold text-gray-500"
            >
              {{ currentFileInfo.extension?.toUpperCase() || 'FILE' }}
            </div>
          </div>
          <div class="flex-1 space-y-2">
            <div>
              <div class="text-xs text-gray-500">
                {{ $t('system.attachment.fileName') }}
              </div>
              <div class="font-medium">{{ currentFileInfo.filename }}</div>
            </div>
            <div>
              <div class="text-xs text-gray-500">
                {{ $t('system.attachment.fileSize') }}
              </div>
              <div>{{ (currentFileInfo.size / 1024).toFixed(2) }} KB</div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <div class="text-xs text-gray-500">
              {{ $t('system.attachment.fileType') }}
            </div>
            <div class="break-words">{{ currentFileInfo.mediatype }}</div>
          </div>
          <div>
            <div class="text-xs text-gray-500">
              {{ $t('system.attachment.directory') }}
            </div>
            <div class="break-words">{{ currentFileInfo.parent || '-' }}</div>
          </div>
          <div class="sm:col-span-2">
            <div class="text-xs text-gray-500">
              {{ $t('system.attachment.uploadTime') }}
            </div>
            <div>
              {{
                new Date(currentFileInfo.created_at * 1000).toLocaleString()
              }}
            </div>
          </div>
          <div class="sm:col-span-2">
            <div class="text-xs text-gray-500">
              {{ $t('system.attachment.fileUrl') }}
            </div>
            <div class="break-all text-sm">{{ currentFileInfo.url }}</div>
          </div>
        </div>
      </div>
    </InfoModal>
  </Page>
</template>

