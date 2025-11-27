import type { OnActionClickFn, VxeTableGridOptions } from '#/adapter/vxe-table';
import type { AdminRoleApi } from '#/api/admin/role';

import { $t } from '#/locales';

export function getRoleStatusOptions() {
  return [
    { color: 'success', label: $t('common.enable'), value: 1 },
    { color: 'error', label: $t('common.disable'), value: 0 },
  ];
}

export function getRoleTypeOptions() {
  return [
    { color: 'default', label: $t('admin.role.typeNormal'), value: 0 },
    { color: 'warning', label: $t('admin.role.typeSuper'), value: 1 },
  ];
}

export function useColumns(
  onActionClick: OnActionClickFn<AdminRoleApi.AdminRole>,
): VxeTableGridOptions<AdminRoleApi.AdminRole>['columns'] {
  return [
    {
      align: 'left',
      field: 'name',
      treeNode: true,
      title: $t('admin.role.name'),
    },
    {
      align: 'center',
      cellRender: { name: 'CellTag', options: getRoleTypeOptions() },
      field: 'is_super',
      title: $t('admin.role.type'),
    },
    {
      align: 'center',
      cellRender: { name: 'CellTag', options: getRoleStatusOptions() },
      field: 'status',
      title: $t('admin.role.status'),
    },
    {
      align: 'center',
      field: 'sort',
      title: $t('common.sort'),
    },
    {
      align: 'center',
      field: 'created_at',
      formatter: ({ cellValue }) => {
        return new Date(cellValue * 1000).toLocaleString();
      },
      title: $t('common.createTime'),
    },
    {
      align: 'right',
      cellRender: {
        attrs: {
          nameField: 'name',
          onClick: onActionClick,
        },
        name: 'CellOperation',
        options: [
          {
            code: 'edit',
            show: (row: AdminRoleApi.AdminRole) => row.id !== 1,
          },
          {
            code: 'delete', 
            show: (row: AdminRoleApi.AdminRole) => row.id !== 1,
          },
        ],
      },
      field: 'operation',
      headerAlign: 'center',
      showOverflow: false,
      title: $t('admin.role.operation'),
    },
  ];
}
