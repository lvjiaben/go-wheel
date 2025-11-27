<script lang="ts" setup>
import type { VxeTableGridOptions } from '#/adapter/vxe-table';

import { $t } from '@vben/locales';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { getOperationLogApi } from '#/api/core/auth';

interface LogItem {
  [key: string]: any;
  created_at: number;
  id: number;
  ip: string;
  status: number;
  username: string;
}

const columns: VxeTableGridOptions<LogItem>['columns'] = [
  {
    field: 'id',
    title: $t('userinfo.log.id'),
    width: 80,
  },
  {
    field: 'username',
    title: $t('userinfo.log.username'),
  },
  {
    field: 'ip',
    title: $t('userinfo.log.ip'),
  },
  {
    cellRender: {
      name: 'CellTag',
      options: [
        { color: 'success', label: $t('userinfo.log.success'), value: 1 },
        { color: 'error', label: $t('userinfo.log.failed'), value: 0 },
      ],
    },
    field: 'status',
    title: $t('userinfo.log.status'),
  },
  {
    field: 'created_at',
    formatter: ({ cellValue }) => {
      return new Date(cellValue * 1000).toLocaleString();
    },
    title: $t('userinfo.log.time'),
  },
];

const [Grid, gridApi] = useVbenVxeGrid({
  gridOptions: {
    columns,
    height: 500,
    keepSource: true,
    pagerConfig: {
      enabled: true,
      pageSize: 10,
    },
    proxyConfig: {
      ajax: {
        query: async ({ page }) => {
          return await getOperationLogApi({
            page: page.currentPage,
            page_size: page.pageSize,
          });
        },
      },
    },
    rowConfig: {
      keyField: 'id',
    },
  } as VxeTableGridOptions,
});
</script>

<template>
  <div>
    <h2 class="mb-4 text-xl font-semibold">{{ $t('userinfo.log.title') }}</h2>
    <Grid />
  </div>
</template>

