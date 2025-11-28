import type { OnActionClickFn, VxeTableGridOptions } from '#/adapter/vxe-table';
import type { AdminAdminApi } from '#/api/admin/admin';

import { $t } from '#/locales';

export function getAdminStatusOptions() {
  return [
    { color: 'success', label: $t('common.enable'), value: 1 },
    { color: 'error', label: $t('common.disable'), value: 0 },
  ];
}

export function useColumns(
  onActionClick: OnActionClickFn<AdminAdminApi.Admin>,
): VxeTableGridOptions<AdminAdminApi.Admin>['columns'] {
  return [
    {
      align: 'center',
      field: 'id',
      title: $t('admin.admin.id'),
      width: 80,
    },
    {
      align: 'left',
      field: 'username',
      title: $t('admin.admin.username'),
      width: 120,
    },
    {
      align: 'left',
      field: 'realname',
      title: $t('admin.admin.realname'),
      width: 120,
    },
    {
      align: 'left',
      field: 'email',
      title: $t('admin.admin.email'),
      minWidth: 180,
    },
    {
      align: 'left',
      field: 'mobile',
      title: $t('admin.admin.mobile'),
      width: 130,
    },
    {
      align: 'center',
      cellRender: { name: 'CellTag', options: getAdminStatusOptions() },
      field: 'status',
      title: $t('admin.admin.status'),
      width: 100,
    },
    {
      align: 'center',
      field: 'last_login_time',
      formatter: ({ cellValue }) => {
        if (!cellValue || cellValue === 0) return '-';
        return new Date(cellValue * 1000).toLocaleString();
      },
      title: $t('admin.admin.lastLoginTime'),
      width: 160,
    },
    {
      align: 'center',
      field: 'created_at',
      formatter: ({ cellValue }) => {
        return new Date(cellValue * 1000).toLocaleString();
      },
      title: $t('common.createTime'),
      width: 160,
    },
    {
      align: 'right',
      cellRender: {
        attrs: {
          nameField: 'username',
          onClick: onActionClick,
        },
        name: 'CellOperation',
        options: [
          {
            code: 'edit',
            show: (row: AdminAdminApi.Admin) => row.id !== 1,
          },
          {
            code: 'delete',
            show: (row: AdminAdminApi.Admin) => row.id !== 1,
          },
        ],
      },
      field: 'operation',
      headerAlign: 'center',
      showOverflow: false,
      title: $t('admin.admin.operation'),
      width: 150,
    },
  ];
}

