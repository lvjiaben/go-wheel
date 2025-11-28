<script lang="ts" setup>
import { computed, ref, watch } from 'vue';

import { $t } from '@vben/locales';
import { useDebounceFn } from '@vueuse/core';
import { Select, Spin, Pagination } from 'ant-design-vue';

import { requestClient } from '#/api/request';

interface RemoteConfig {
  api: string;
  searchField?: string;
  labelField?: string;
  valueField?: string;
  imageField?: string;
  descField?: string;
  pageSize?: number;
}

interface Props {
  value?: string | string[] | number;
  config?: string | RemoteConfig;
  multiple?: boolean;
  placeholder?: string;
  disabled?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  value: undefined,
  config: '',
  multiple: false,
  placeholder: '',
  disabled: false,
});

const emit = defineEmits<{
  'update:value': [value: string | string[] | number];
  change: [value: string | string[] | number];
}>();

interface OptionItem {
  label: string;
  value: string | number;
  image?: string;
  description?: string;
}

// 解析配置
const parseConfig = (config: string | RemoteConfig): RemoteConfig => {
  if (typeof config === 'string') {
    try {
      return JSON.parse(config);
    } catch {
      return { api: '' };
    }
  }
  return config;
};

const remoteConfig = computed(() => parseConfig(props.config));

// 选项列表
const options = ref<OptionItem[]>([]);
const fetching = ref(false);

// 分页状态
const currentPage = ref(1);
const totalCount = ref(0);
const pageSize = computed(() => remoteConfig.value.pageSize || 10);
const currentSearch = ref('');

// 已选中的选项缓存（用于显示标签）
const selectedOptionsCache = ref<Map<string, OptionItem>>(new Map());

// 本地值 - 处理各种类型
const localValue = computed({
  get: () => {
    // 如果是多选模式
    if (props.multiple) {
      if (!props.value) return [];
      if (Array.isArray(props.value)) {
        return Array.from(new Set(props.value.map(String)));
      }
      if (typeof props.value === 'string') {
        if (props.value === '') return [];
        try {
          const parsed = JSON.parse(props.value);
          if (Array.isArray(parsed)) {
            return Array.from(new Set(parsed.map(String)));
          }
        } catch {
          return Array.from(new Set(props.value.split(',').filter(Boolean)));
        }
      }
      return [];
    }
    // 单选模式
    if (props.value === null || props.value === undefined || props.value === '') {
      return undefined;
    }
    return String(props.value);
  },
  set: (val: any) => {
    if (props.multiple && Array.isArray(val)) {
      const uniqueValues = Array.from(new Set(val));
      handleChange(uniqueValues);
    } else {
      handleChange(val);
    }
  },
});

// 合并选项列表（当前页 + 已选中的选项）
const mergedOptions = computed(() => {
  const optionsMap = new Map<string, OptionItem>();

  // 先添加当前页的选项
  options.value.forEach(option => {
    optionsMap.set(String(option.value), option);
  });

  // 再添加已选中的选项（确保已选中的值能显示 label）
  const selectedValues = Array.isArray(localValue.value)
    ? localValue.value
    : (localValue.value ? [localValue.value] : []);

  selectedValues.forEach((value: string) => {
    if (!optionsMap.has(value)) {
      const cached = selectedOptionsCache.value.get(value);
      if (cached) {
        optionsMap.set(value, cached);
      }
    }
  });

  return Array.from(optionsMap.values());
});

// 获取远程数据
const fetchOptions = async (search?: string, page: number = 1) => {
  if (!remoteConfig.value.api) return;

  fetching.value = true;
  currentSearch.value = search || '';
  currentPage.value = page;

  try {
    const params: any = {
      page,
      page_size: pageSize.value,
    };
    if (search) {
      params[remoteConfig.value.searchField || 'search'] = search;
    }

    const result = await requestClient.get(remoteConfig.value.api, { params });

    const list = result.list || result.data || result || [];
    totalCount.value = result.total || list.length;

    const newOptions = list.map((item: any) => {
      const value = item[remoteConfig.value.valueField || 'id'];
      return {
        label: item[remoteConfig.value.labelField || 'name'],
        value: String(value),
        image: item[remoteConfig.value.imageField || 'image'],
        description: item[remoteConfig.value.descField || 'description'],
      };
    });

    options.value = newOptions;

    // 将新选项添加到缓存中
    newOptions.forEach((option: OptionItem) => {
      selectedOptionsCache.value.set(String(option.value), option);
    });
  } catch (error) {
    console.error($t('common.components.tableSelect.fetchError'), error);
  } finally {
    fetching.value = false;
  }
};

// 分页变化
const handlePageChange = (page: number) => {
  fetchOptions(currentSearch.value, page);
};

// 处理选择变化
const handleChange = (val: any) => {
  // 当选中某个选项时，将其添加到缓存
  if (val) {
    const selectedValues = Array.isArray(val) ? val : [val];
    selectedValues.forEach((value: string) => {
      const option = options.value.find(opt => String(opt.value) === String(value));
      if (option) {
        selectedOptionsCache.value.set(String(value), option);
      }
    });
  }
  emit('update:value', val);
  emit('change', val);
};

// 防抖搜索（搜索时重置到第一页）
const handleSearch = useDebounceFn((value: string) => {
  fetchOptions(value, 1);
}, 300);

// 根据初始值加载对应的选项（用于显示 label）
const fetchInitialOptions = async () => {
  if (!props.value || !remoteConfig.value.api) return;

  // 获取需要查询的值
  const values: string[] = [];
  if (Array.isArray(props.value)) {
    values.push(...props.value.map(String));
  } else if (typeof props.value === 'string' && props.value !== '') {
    try {
      const parsed = JSON.parse(props.value);
      if (Array.isArray(parsed)) {
        values.push(...parsed.map(String));
      } else {
        values.push(props.value);
      }
    } catch {
      values.push(...props.value.split(',').filter(Boolean));
    }
  } else {
    values.push(String(props.value));
  }

  // 检查哪些值不在缓存中
  const missingValues = values.filter(v => !selectedOptionsCache.value.has(v));
  if (missingValues.length === 0) return;

  // 查询缺失的选项
  try {
    const params: any = {
      filter: JSON.stringify({ [remoteConfig.value.valueField || 'id']: missingValues }),
    };
    const result = await requestClient.get(remoteConfig.value.api, { params });
    const list = result.list || result.data || result || [];

    list.forEach((item: any) => {
      const value = item[remoteConfig.value.valueField || 'id'];
      const option: OptionItem = {
        label: item[remoteConfig.value.labelField || 'name'],
        value: String(value),
        image: item[remoteConfig.value.imageField || 'image'],
        description: item[remoteConfig.value.descField || 'description'],
      };
      selectedOptionsCache.value.set(String(value), option);
    });
  } catch (error) {
    console.error('获取初始选项失败', error);
  }
};

// 初始加载
watch(
  () => props.config,
  () => {
    fetchOptions();
  },
  { immediate: true },
);

// 监听初始值变化，加载对应的 label
watch(
  () => props.value,
  () => {
    fetchInitialOptions();
  },
  { immediate: true },
);
</script>

<template>
  <Select
    v-model:value="localValue"
    :disabled="disabled"
    :filter-option="false"
    :loading="fetching"
    :mode="multiple ? 'multiple' : undefined"
    :not-found-content="fetching ? undefined : null"
    :options="mergedOptions"
    :placeholder="placeholder || $t('common.components.tableSelect.placeholder')"
    show-search
    @search="handleSearch"
  >
    <template #notFoundContent>
      <Spin v-if="fetching" size="small" />
      <span v-else>{{ $t('common.components.tableSelect.noData') }}</span>
    </template>

    <!-- 自定义选项显示 -->
    <template #option="{ label, image, description }">
      <div class="flex items-center gap-2">
        <img
          v-if="image"
          :alt="label"
          :src="image"
          class="h-8 w-8 rounded object-cover"
        />
        <div class="flex-1">
          <div class="font-medium">{{ label }}</div>
          <div v-if="description" class="text-xs text-gray-500">
            {{ description }}
          </div>
        </div>
      </div>
    </template>

    <!-- 底部分页器 -->
    <template #dropdownRender="{ menuNode }">
      <div>
        <component :is="menuNode" />
        <div
          v-if="totalCount > pageSize"
          class="border-t border-gray-200 px-2 py-2 flex justify-center"
        >
          <Pagination
            :current="currentPage"
            :page-size="pageSize"
            :total="totalCount"
            :show-size-changer="false"
            :show-total="(total: number) => `共 ${total} 条`"
            size="small"
            show-quick-jumper
            @change="handlePageChange"
            @mousedown.stop
          />
        </div>
      </div>
    </template>
  </Select>
</template>


