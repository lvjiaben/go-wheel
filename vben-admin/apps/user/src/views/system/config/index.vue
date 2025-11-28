<script lang="ts" setup>
import type { ConfigApi } from '#/api/system/config';

import { computed, onMounted, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';
import { $t } from '@vben/locales';

import {
  AutoComplete as AAutoComplete,
  Button,
  Input as AInput,
  message,
  Modal,
  Select as ASelect,
  Tabs,
  Textarea as ATextarea,
} from 'ant-design-vue';

import {
  createConfig,
  deleteConfig,
  getConfigList,
  updateConfigs,
} from '#/api/system/config';

import ConfigFormItem from './modules/config-form-item.vue';

// 配置数据
const configGroups = ref<ConfigApi.ConfigGroup[]>([]);
const loading = ref(false);
const activeTab = ref<string>('');

// 表单数据 - 存储所有配置项的值
const formData = ref<Record<string, any>>({});

// 加载配置列表
const loadConfigs = async () => {
  loading.value = true;
  try {
    const result = await getConfigList();
    configGroups.value = result.list || [];

    // 设置默认激活的tab
    if (configGroups.value.length > 0 && !activeTab.value) {
      activeTab.value = configGroups.value[0]?.dir || '';
    }

    // 初始化表单数据
    configGroups.value.forEach((group) => {
      group.children.forEach((config) => {
        formData.value[config.key] = config.value;
      });
    });
  } finally {
    loading.value = false;
  }
};

// 保存所有配置
const handleSave = async () => {
  const hideLoading = message.loading($t('system.config.saving'), 0);
  try {
    await updateConfigs(formData.value);
    message.success($t('system.config.saveSuccess'));
    await loadConfigs();
  } catch (error) {
    message.error($t('system.config.saveFailed'));
  } finally {
    hideLoading();
  }
};

// 删除配置项
const handleDelete = (config: ConfigApi.Config) => {
  Modal.confirm({
    title: $t('common.confirmDelete'),
    content: $t('system.config.deleteConfirm', [config.name]),
    onOk: async () => {
      const hideLoading = message.loading($t('system.config.deleting'), 0);
      try {
        await deleteConfig(config.id);
        message.success($t('system.config.deleteSuccess'));
        await loadConfigs();
      } catch (error) {
        message.error($t('system.config.deleteFailed'));
      } finally {
        hideLoading();
      }
    },
  });
};

// 更新表单值
const handleValueChange = (key: string, value: string) => {
  formData.value[key] = value;
};

// 获取所有 dir 选项
const dirOptions = computed(() => {
  return configGroups.value.map((group) => ({
    label: group.dir,
    value: group.dir,
  }));
});

// 创建配置表单
const createFormData = ref<Partial<ConfigApi.Config>>({
  dir: '',
  key: '',
  name: '',
  tip: '',
  type: 'input',
  value: '',
  variable: '',
});

// 配置类型选项
const typeOptions = computed(() => [
  { label: $t('system.config.typeInput'), value: 'input' },
  { label: $t('system.config.typeText'), value: 'text' },
  { label: $t('system.config.typeEditor'), value: 'editor' },
  { label: $t('system.config.typeFile'), value: 'file' },
  { label: $t('system.config.typeImage'), value: 'image' },
  { label: $t('system.config.typeImages'), value: 'images' },
  { label: $t('system.config.typeDate'), value: 'date' },
  { label: $t('system.config.typeTime'), value: 'time' },
  { label: $t('system.config.typeDatetime'), value: 'datetime' },
  { label: $t('system.config.typeDatetimerange'), value: 'datetimerange' },
  { label: $t('system.config.typeSelect'), value: 'select' },
  { label: $t('system.config.typeSelects'), value: 'selects' },
  { label: $t('system.config.typeSwitch'), value: 'switch' },
  { label: $t('system.config.typeCheckbox'), value: 'checkbox' },
  { label: $t('system.config.typeRadio'), value: 'radio' },
  { label: $t('system.config.typeArray'), value: 'array' },
  { label: $t('system.config.typeTableselect'), value: 'tableselect' },
  { label: $t('system.config.typeTableselects'), value: 'tableselects' },
]);

// 创建配置 Modal
const [CreateModal, createModalApi] = useVbenModal({
  title: $t('system.config.create'),
  class: 'w-full max-w-2xl',
  onConfirm: async () => {
    // 验证必填字段
    if (!createFormData.value.dir) {
      message.error($t('system.config.groupRequired'));
      return;
    }
    if (!createFormData.value.key) {
      message.error($t('system.config.keyRequired'));
      return;
    }
    if (!createFormData.value.name) {
      message.error($t('system.config.nameRequired'));
      return;
    }

    const hideLoading = message.loading($t('system.config.creating'), 0);
    try {
      await createConfig(createFormData.value);
      message.success($t('system.config.createSuccess'));
      createModalApi.close();
      // 重置表单
      createFormData.value = {
        dir: '',
        key: '',
        name: '',
        tip: '',
        type: 'input',
        value: '',
        variable: '',
      };
      // 刷新列表
      await loadConfigs();
    } catch (error) {
      message.error($t('system.config.createFailed'));
    } finally {
      hideLoading();
    }
  },
  onCancel: () => {
    // 重置表单
    createFormData.value = {
      dir: '',
      key: '',
      name: '',
      tip: '',
      type: 'input',
      value: '',
      variable: '',
    };
    // 关闭弹窗
    createModalApi.close();
  },
});

// 打开创建配置弹窗
const handleCreate = () => {
  createModalApi.open();
};

// 帮助弹窗
const [HelpModal, helpModalApi] = useVbenModal({
  title: $t('system.config.helpTitle'),
  class: 'w-full max-w-4xl',
  showCancelButton: false,
  confirmText: $t('system.config.helpClose'),
  onConfirm: () => {
    helpModalApi.close();
  },
});

// 打开帮助弹窗
const handleShowHelp = () => {
  helpModalApi.open();
};

onMounted(() => {
  loadConfigs();
});
</script>

<template>
  <Page
    content-class="flex flex-col"
    :description="$t('system.config.description')"
    :title="$t('system.config.title')"
  >
    <div class="flex h-full flex-col rounded-lg border bg-white">
      <Tabs v-model:activeKey="activeTab" class="flex-1 px-4">
        <Tabs.TabPane
          v-for="group in configGroups"
          :key="group.dir"
          :tab="group.dir"
        >
          <div class="space-y-4 overflow-y-auto pb-4" style="max-height: calc(100vh - 320px)">
            <ConfigFormItem
              v-for="config in group.children"
              :key="config.id"
              :config="config"
              :value="formData[config.key]"
              @delete="handleDelete(config)"
              @update:value="(val: string) => handleValueChange(config.key, val)"
            />
          </div>
        </Tabs.TabPane>
      </Tabs>

      <!-- 底部按钮 -->
      <div class="flex items-center justify-between border-t bg-white p-4">
        <!-- 创建按钮 -->
        <Button size="middle" @click="handleCreate">
          <IconifyIcon icon="mdi:plus" class="size-4" />
        </Button>

        <!-- 保存按钮 -->
        <Button type="primary" size="middle" :loading="loading" @click="handleSave">
          <IconifyIcon icon="mdi:content-save" class="mr-1 size-4" />
          {{ $t('system.config.save') }}
        </Button>
      </div>
    </div>

    <!-- 创建配置弹窗 -->
    <CreateModal>
      <div class="space-y-4 p-4">
        <div>
          <label class="mb-1 block text-sm font-medium">{{ $t('system.config.group') }} <span class="text-red-500">*</span></label>
          <a-auto-complete
            v-model:value="createFormData.dir"
            :options="dirOptions"
            :placeholder="$t('system.config.groupPlaceholder')"
            class="w-full"
          />
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium">{{ $t('system.config.key') }} <span class="text-red-500">*</span></label>
          <a-input
            v-model:value="createFormData.key"
            :placeholder="$t('system.config.keyPlaceholder')"
          />
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium">{{ $t('system.config.name') }} <span class="text-red-500">*</span></label>
          <a-input
            v-model:value="createFormData.name"
            :placeholder="$t('system.config.namePlaceholder')"
          />
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium">{{ $t('system.config.type') }}</label>
          <a-select
            v-model:value="createFormData.type"
            :options="typeOptions"
            :placeholder="$t('system.config.typePlaceholder')"
            class="w-full"
          />
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium">{{ $t('system.config.value') }}</label>
          <a-textarea
            v-model:value="createFormData.value"
            :placeholder="$t('system.config.valuePlaceholder')"
            :rows="3"
          />
        </div>

        <div>
          <label class="mb-1 block text-sm font-medium">{{ $t('system.config.tip') }}</label>
          <a-input
            v-model:value="createFormData.tip"
            :placeholder="$t('system.config.tipPlaceholder')"
          />
        </div>

        <div>
          <div class="mb-1 flex items-center justify-between">
            <label class="text-sm font-medium">{{ $t('system.config.variable') }}</label>
            <Button size="small" type="link" @click="handleShowHelp">
              <IconifyIcon icon="mdi:help-circle-outline" class="mr-1 size-4" />
              {{ $t('system.config.viewHelp') }}
            </Button>
          </div>
          <a-textarea
            v-model:value="createFormData.variable"
            :placeholder="$t('system.config.variablePlaceholder')"
            :rows="3"
          />
        </div>
      </div>
    </CreateModal>

    <!-- 帮助弹窗 -->
    <HelpModal>
      <div class="max-h-[70vh] space-y-6 overflow-y-auto p-4">
        <!-- select/selects/radio/checkbox -->
        <div>
          <h3 class="mb-3 text-base font-semibold">1. {{ $t('system.config.helpSection1Title') }}</h3>

          <div class="mb-3">
            <p class="mb-2 text-sm font-medium text-gray-700">{{ $t('system.config.helpSection1ObjectFormat') }}</p>
            <pre class="rounded bg-gray-50 p-3 text-sm"><code>{
  "option1": "选项1",
  "option2": "选项2",
  "option3": "选项3"
}</code></pre>
          </div>

          <div>
            <p class="mb-2 text-sm font-medium text-gray-700">{{ $t('system.config.helpSection1ArrayFormat') }}</p>
            <pre class="rounded bg-gray-50 p-3 text-sm"><code>[
  {"label": "选项1", "value": "option1"},
  {"label": "选项2", "value": "option2"}
]</code></pre>
          </div>
        </div>

        <!-- tableselect/tableselects -->
        <div>
          <h3 class="mb-3 text-base font-semibold">2. {{ $t('system.config.helpSection2Title') }}</h3>

          <div class="mb-3">
            <p class="mb-2 text-sm text-gray-600">{{ $t('system.config.helpSection2Desc') }}</p>
            <pre class="rounded bg-gray-50 p-3 text-sm"><code>{
  "api": "/api/products/list",
  "searchField": "keyword",
  "labelField": "name",
  "valueField": "id",
  "imageField": "image",
  "descField": "description"
}</code></pre>
          </div>

          <div class="space-y-2 text-sm">
            <p><span class="font-medium text-blue-600">api</span>: {{ $t('system.config.helpSection2Api') }}</p>
            <p><span class="font-medium text-blue-600">searchField</span>: {{ $t('system.config.helpSection2SearchField') }}</p>
            <p><span class="font-medium text-blue-600">labelField</span>: {{ $t('system.config.helpSection2LabelField') }}</p>
            <p><span class="font-medium text-blue-600">valueField</span>: {{ $t('system.config.helpSection2ValueField') }}</p>
            <p><span class="font-medium text-blue-600">imageField</span>: {{ $t('system.config.helpSection2ImageField') }}</p>
            <p><span class="font-medium text-blue-600">descField</span>: {{ $t('system.config.helpSection2DescField') }}</p>
          </div>
        </div>

        <!-- array -->
        <div>
          <h3 class="mb-3 text-base font-semibold">3. {{ $t('system.config.helpSection3Title') }}</h3>

          <div class="mb-3">
            <p class="mb-2 text-sm text-gray-600">{{ $t('system.config.helpSection3Desc') }}</p>
            <pre class="rounded bg-gray-50 p-3 text-sm"><code>{"label1": "value1", "label2": "value2"}</code></pre>
          </div>
        </div>

        <!-- 其他类型 -->
        <div>
          <h3 class="mb-3 text-base font-semibold">4. {{ $t('system.config.helpSection4Title') }}</h3>

          <div class="space-y-2 text-sm text-gray-600">
            <p><span class="font-medium">input / text / editor</span>: {{ $t('system.config.helpSection4Desc1') }}</p>
            <p><span class="font-medium">file / image / images</span>: {{ $t('system.config.helpSection4Desc2') }}</p>
            <p><span class="font-medium">date / time / datetime / datetimerange</span>: {{ $t('system.config.helpSection4Desc3') }}</p>
            <p><span class="font-medium">switch</span>: {{ $t('system.config.helpSection4Desc4') }}</p>
          </div>
        </div>
      </div>
    </HelpModal>
  </Page>
</template>
