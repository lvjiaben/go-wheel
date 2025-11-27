import type { OnActionClickFn, VxeTableGridOptions } from '#/adapter/vxe-table';
import type { AdminMenuApi } from '#/api/admin/menu';

import { $t } from '#/locales';

export function getMenuTypeOptions() {
  return [
    { color: 'default', label: $t('admin.menu.typeMenu'), value: 'menu' },
    { color: 'error', label: $t('admin.menu.typeButton'), value: 'button' },
    {
      color: 'success',
      label: $t('admin.menu.typeEmbedded'),
      value: 'iframe',
    },
    { color: 'warning', label: $t('admin.menu.typeLink'), value: 'link' },
  ];
}

export function useColumns(
  onActionClick: OnActionClickFn<AdminMenuApi.AdminMenu>,
): VxeTableGridOptions<AdminMenuApi.AdminMenu>['columns'] {
  return [
    {
      align: 'left',
      field: 'name',
      slots: { default: 'title' },
      title: $t('admin.menu.menuTitle'),
      treeNode: true,
      width: 250,
    },
    {
      align: 'center',
      cellRender: { name: 'CellTag', options: getMenuTypeOptions() },
      field: 'type',
      title: $t('admin.menu.type'),
      width: 100,
    },
    {
      align: 'left',
      field: 'sort',
      title: $t('common.sort'),
      width: 100,
    },
    {
      align: 'left',
      field: 'path',
      title: $t('admin.menu.path'),
      width: 200,
    },
    {
      align: 'left',
      field: 'component',
      formatter: ({ row }) => {
        switch (row.type) {
          case 'button': {
            return row.permission ?? '';
          }
          case 'menu': {
            return row.component ?? '';
          }
          case 'iframe': {
            return row.iframe ?? '';
          }
          case 'link': {
            return row.external ?? '';
          }
        }
        return '';
      },
      minWidth: 200,
      title: $t('admin.menu.component'),
    },
    {
      slots: { default: 'visible' },
      field: 'visible',
      title: $t('admin.menu.status'),
      width: 100,
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
            code: 'append',
            text: '新增下级',
          },
          'edit', // 默认的编辑按钮
          'delete', // 默认的删除按钮
        ],
      },
      field: 'operation',
      // fixed: 'right',
      headerAlign: 'center',
      showOverflow: false,
      title: $t('admin.menu.operation'),
      width: 200,
    },
  ];
}
