<script lang="ts" setup>
import type {
  VxeTableGridOptions,
} from '#/adapter/vxe-table';
import { onMounted, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { VbenButton, VbenButtonGroup ,Page, useVbenDrawer } from '@vben/common-ui';
import { Plus,Trash,Check,X,IconifyIcon } from '@vben/icons';
import { $t } from '@vben/locales';

import { message, Popconfirm,Tag,Image, InputSearch } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import {
  deleteUser,
  getUserList,
  operateUser,
} from '#/api/user';

import type { UserApi } from '#/api/user';

import { useColumns, useSearchFormSchema } from './data';
import FormDrawerComponent from './modules/form-drawer.vue';
import MoneyScoreDrawerComponent from './modules/money-score-drawer.vue';

const route = useRoute();
const router = useRouter();

// 选中的行数
const selectedCount = ref(0);

// 模糊搜索内容
const searchValue = ref<string>("");

// 获取选中的行数据（公共方法）
const getSelectedRecords = (): UserApi.User[] => {
  return (gridApi.grid?.getCheckboxRecords() as UserApi.User[]) || [];
}

// 更新选中行数量
const updateSelectedCount = () => {
  requestAnimationFrame(() => {
    selectedCount.value = getSelectedRecords().length;
  });
}

// 清空搜索
const onClearSearch = () => {
  searchValue.value = "";
  onRefresh();
}

// 搜索表单
const onSearch = (content: string) => {
  searchValue.value = content;
  onRefresh();
}

const [FormDrawer, formDrawerApi] = useVbenDrawer({
  connectedComponent: FormDrawerComponent,
  destroyOnClose: true,
});

const [MoneyScoreDrawer, moneyScoreDrawerApi] = useVbenDrawer({
  connectedComponent: MoneyScoreDrawerComponent,
  destroyOnClose: true,
});

const [Grid, gridApi] = useVbenVxeGrid({
  formOptions: { // 表单配置
    collapsed: true,// 默认true 关闭展开
    schema: useSearchFormSchema(),
    showCollapseButton: true,// 控制表单是否显示折叠按钮
    submitOnChange: false,// 是否在字段值改变时提交表单
    submitOnEnter: false,// 按下回车时是否提交表单
  },
  gridOptions: { // 表格配置
    columns: useColumns(),
    keepSource: true,
    pagerConfig: {
      enabled: true,
    },
    proxyConfig: {
      ajax: {
        query: async ({ page, sort }, formValues) => {
          const params: UserApi.ListParams = {
            page: page.currentPage,
            page_size: page.pageSize,
            sort_by: sort.field ? sort.field : 'id',
            sort_order: sort.order === 'asc' ? 'asc' : 'desc',
            filter: JSON.stringify(formValues),
            search: searchValue.value || undefined,
          };
          return await getUserList(params);
        },
      },
      sort: true,
    },
    rowConfig: {
      keyField: 'id',
      isHover: true,
    },
    sortConfig: {
      remote: true,
      trigger: 'cell',
    },
    stripe: false,
    toolbarConfig: {
      custom: true,
      refresh: true,
      zoom: true,
    },
    checkboxConfig: {
      trigger: 'cell',
      highlight: true,
      reserve: false,
    },
    editConfig: {
      trigger: 'dblclick',
      mode: 'cell',
      showStatus: true,
    },
  } as VxeTableGridOptions,
  gridEvents: {
    checkboxChange: updateSelectedCount,
    checkboxAll: updateSelectedCount,
    editClosed: async ({ row, column }: { row: any; column: any }) => {
      // 双击编辑完成后的回调
      if (column.field === 'tid' || column.field === 'pid') {
        const newValue = row[column.field];
        // 调用 operateApi 更新字段
        await operateApi(newValue, row, column.field);
      }
    },
  },
});

// 通用字段操作方法（支持 status、tid、pid 等字段）
const operateApi = async (value: number, row?: UserApi.User, field: string = 'status') => {
  let ids: number[] = [];

  if (row) {
    ids = [row.id!];
  } else {
    const selectRecords = getSelectedRecords();
    if (selectRecords.length === 0) {
      message.warning($t('common.tableSelectTip'));
      return;
    }
    ids = selectRecords.map((item) => item.id!);
  }

  try {
    await operateUser({
      ids: ids,
      field,
      value,
    });
    message.success($t('common.success'));
    onRefresh();
    return true;
  } catch (error) {
    message.error('error');
    return false;
  }
}

const onRefresh = () => {
  gridApi.grid?.clearCheckboxRow();
  gridApi.query();
  selectedCount.value = 0;
}

const onEdit = (row: UserApi.User) => {
  formDrawerApi.setData(row).open();
}

const onCreate = () => {
  formDrawerApi.setData({}).open();
}

const onUpdateMoney = (row: UserApi.User) => {
  moneyScoreDrawerApi.setData({ ...row, type: 'money' }).open();
}

const onUpdateScore = (row: UserApi.User) => {
  moneyScoreDrawerApi.setData({ ...row, type: 'score' }).open();
}

const onDelete = (row?: UserApi.User) => {
  let ids: number[] = [];

  if (row) {
    ids = [row.id!];
  } else {
    const selectRecords = getSelectedRecords();
    if (selectRecords.length === 0) {
      message.warning($t('common.tableSelectTip'));
      return;
    }
    ids = selectRecords.map((item) => item.id!);
  }

  const hideLoading = message.loading({
    content: $t('ui.actionMessage.deleting'),
    duration: 0,
    key: 'action_process_msg',
  });
  deleteUser({ ids: ids })
    .then(() => {
      message.success({
        content: $t('ui.actionMessage.deleteSuccess'),
        key: 'action_process_msg',
      });
      onRefresh();
    })
    .catch(() => {
      hideLoading();
    });
}

// 获取搜索表单的字段名列表
const searchFormSchema = useSearchFormSchema();
const searchFormFields = searchFormSchema
  ? searchFormSchema.map((item) => item.fieldName).filter((name): name is string => !!name)
  : [];

// 页面加载时从 URL 参数填充搜索表单
onMounted(() => {
  const query = route.query;
  if (Object.keys(query).length > 0) {
    const formValues: Record<string, any> = {};
    for (const [key, value] of Object.entries(query)) {
      // 只填充存在于搜索表单中的字段
      if (searchFormFields.includes(key) && value) {
        formValues[key] = value;
      }
    }
    if (Object.keys(formValues).length > 0) {
      // 设置表单值并触发搜索
      gridApi.formApi.setValues(formValues);
    }
  }
});
</script>
<template>
  <Page>
    <FormDrawer @success="onRefresh" />
    <MoneyScoreDrawer @success="onRefresh" />
    <Grid>
      <template #pid="{ row }">
        <span
          class="text-blue-600 underline cursor-pointer hover:text-blue-800"
          @click="router.push({ path: '/user/list', query: { id: row.pid } })"
      >
          {{ row.pid }}
      </span>
      </template>
      <template #tid="{ row }">
        <span
          class="text-blue-600 underline cursor-pointer hover:text-blue-800"
          @click="router.push({ path: '/user/list', query: { id: row.tid } })"
      >
          {{ row.tid }}
      </span>
      </template>
      <template #money="{ row }">
        <span
          class="text-blue-600 cursor-pointer hover:underline"
          @click="onUpdateMoney(row)"
        >
          ¥{{ Number(row.money || 0).toFixed(2) }}
        </span>
      </template>
      <template #score="{ row }">
        <span
          class="text-blue-600 cursor-pointer hover:underline"
          @click="onUpdateScore(row)"
        >
          {{ Number(row.score || 0).toFixed(2) }}
        </span>
      </template>
      <template #avatar="{ row }">
        <Image
          :src="row.avatar"
          :width="30"
          :height="30"
          :preview="true"
          loading="lazy"
        />
      </template>
      <template #status="{ row }">
        <Popconfirm
          :title="$t('common.confirmTip')"
          @confirm="operateApi(row.status === 1 ? 0 : 1, row, 'status')"
        >
          <Tag :color="row.status === 1 ? 'success' : 'error'" class="ml-2 text-center">
            {{ row.status === 1 ? $t('common.enable') : $t('common.disable') }}
          </Tag>
        </Popconfirm>
      </template>
      <template #operation="{ row }">
        <div class="flex gap-2 justify-center">
          <VbenButtonGroup border>
            <VbenButton variant="icon" @click="onEdit(row)">
              <IconifyIcon
                  class="size-6 outline-none text-blue-600"
                  icon="ant-design:edit-outlined"
                />
            </VbenButton>
            <Popconfirm
              :title="$t('ui.actionTitle.delete')"
              @confirm="onDelete(row)"
            >
              <VbenButton variant="icon">
                <X class="size-6 text-red-600" />
              </VbenButton>
            </Popconfirm>
          </VbenButtonGroup>
        </div>
      </template>
      <template #toolbar-actions>
        <VbenButton size="sm" @click="onCreate" variant="outline" class="mr-3">
          <Plus class="size-3" />
          {{ $t('common.create') }}
        </VbenButton>
        <VbenButtonGroup v-show="selectedCount > 0" border>
          <Popconfirm
            :title="$t('common.confirmDeleteSelected', [selectedCount])"
            :ok-text="$t('common.confirm')"
            :cancel-text="$t('common.cancel')"
            @confirm="onDelete()"
          >
            <VbenButton variant="outline">
              <Trash class="size-3" />
              {{ $t('common.delete') }}
            </VbenButton>
          </Popconfirm>
          <Popconfirm
            :title="$t('common.confirmEnableSelected', [selectedCount])"
            :ok-text="$t('common.confirm')"
            :cancel-text="$t('common.cancel')"
            @confirm="operateApi(1, undefined, 'status')"
          >
            <VbenButton variant="outline">
              <Check class="size-3" />
              {{ $t('common.enable') }}
            </VbenButton>
          </Popconfirm>
          <Popconfirm
            :title="$t('common.confirmDisableSelected', [selectedCount])"
            :ok-text="$t('common.confirm')"
            :cancel-text="$t('common.cancel')"
            @confirm="operateApi(0, undefined, 'status')"
          >
            <VbenButton variant="outline">
              <X class="size-3" />
              {{ $t('common.disable') }}
            </VbenButton>
          </Popconfirm>
        </VbenButtonGroup>
      </template>
      <template #toolbar-tools>
        <InputSearch @search="onSearch" @clear="onClearSearch" allow-clear :placeholder="$t('common.fuzzySearch')" />
      </template>
    </Grid>
  </Page>
</template>
