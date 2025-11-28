import type { VxeTableGridOptions } from '#/adapter/vxe-table';
import type { UserApi } from '#/api/user';
import type { VbenFormProps } from '#/adapter/form';
import { $t } from '#/locales';
export function getUserStatusOptions() {
  return [
    { color: 'success', label: $t('common.enable'), value: 1 },
    { color: 'error', label: $t('common.disable'), value: 0 },
  ];
}

export function useSearchFormSchema(): VbenFormProps['schema'] {
  return [
    {
      component: 'Input',
      fieldName: 'id',
      label: 'ID',
    },
    {
      component: 'Input',
      fieldName: 'mobile',
      label: $t('user.mobile'),
    },
    {
      component: 'Input',
      fieldName: 'username',
      label: $t('user.username'),
    },
    {
      component: 'TableSelect',
      fieldName: 'pid',
      label: $t('user.pid'),
      componentProps:{
        multiple:false,
        config: {
          api: '/user/list',
          labelField: 'username',
          valueField: 'id',
          imageField: 'avatar',
          descField: 'id',
          pageSize: 10, 
        },
      }
    },
    {
      component: 'TableSelect',
      fieldName: 'tid',
      label: $t('user.tid'),
      componentProps:{
        multiple:false,
        config: {
          api: '/user/list',
          labelField: 'username',
          valueField: 'id',
          imageField: 'avatar',
          descField: 'id',
          pageSize: 10, 
        },
      }
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [
          {
            label: $t('common.all'),
            value: '',
          },
          {
            label: $t('common.enable'),
            value: '1',
          },
          {
            label: $t('common.disable'),
            value: '0',
          },
        ],
      },
      fieldName: 'status',
      label: $t('user.status'),
    },
    {
      component: 'RangePicker',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        timestampFormat: true,
      },
      fieldName: 'created_at',
      label: $t('common.createTime'),
    },
    {
      component: 'RangePicker',
      componentProps: {
        showTime: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        timestampFormat: true,
      },
      fieldName: 'updated_at',
      label: $t('common.updateTime'),
    },
  ];
}

export function useColumns(): VxeTableGridOptions<UserApi.User>['columns'] {
  return [
    {
      type: 'checkbox',
      width: 50,
      align: 'center',
      fixed: 'left',
    },
    {
      align: 'center',
      field: 'id',
      title: "ID",
      width: 80,
      sortable: true,
    },
    {
      align: 'center',
      field: 'pid',
      title: $t('user.pid'),
      editRender: { name: 'input' },
      slots: { default: 'pid' },
      width: 80,
    },
    {
      align: 'center',
      field: 'tid',
      title: $t('user.tid'),
      editRender: { name: 'input' },
      slots: { default: 'tid' },
      width: 80,
    },
    {
      field: 'avatar',
      slots: { default: 'avatar' },
      title: $t('user.avatar'),
      width: 50,
    },
    {
      align: 'left',
      field: 'username',
      title: $t('user.username'),
    },
    {
      align: 'left',
      field: 'mobile',
      title: $t('user.mobile'),
    },
    {
      align: 'right',
      field: 'money',
      slots: { default: 'money' },
      title: $t('user.money'),
      width: 100,
      sortable: true,
    },
    {
      align: 'right',
      field: 'score',
      title: $t('user.score'),
      width: 100,
      sortable: true,
      slots: { default: 'score' },
    },
    {
      align: 'center',
      field: 'status',
      slots: { default: 'status' },
      title: $t('user.status'),
      width: 80,
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
      align: 'center',
      field: 'updated_at',
      formatter: ({ cellValue }) => {
        return new Date(cellValue * 1000).toLocaleString();
      },
      title: $t('common.updateTime'),
      width: 160,
    },
    {
      field: 'operation',
      fixed: 'right',
      slots: { default: 'operation' },
      title: $t('common.operation'),
      width: 100,
    }
  ];
}

