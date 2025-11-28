<script lang="ts" setup>
import type { DataNode } from 'ant-design-vue/es/tree';

import type { Recordable } from '@vben/types';

import type { VbenFormSchema } from '#/adapter/form';

import { computed, nextTick, ref } from 'vue';

import { Tree, useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';
import { getPopupContainer } from '@vben/utils';

import { Spin, Tag } from 'ant-design-vue';

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import { useVbenForm, z } from '#/adapter/form';
import {
  saveRole,
  getRoleList,
  getMyMenus,
  AdminRoleApi,
} from '#/api/admin/role';
import { $t } from '#/locales';

import { getRoleStatusOptions, getRoleTypeOptions } from '../data';
import { getMenuTypeOptions } from '#/views/admin/menu/data';

const emit = defineEmits<{
  success: [];
}>();

const formData = ref<AdminRoleApi.AdminRole>();
const permissions = ref<DataNode[]>([]);
const loadingPermissions = ref(false);

const schema: VbenFormSchema[] = [
{
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getRoleStatusOptions(),
      optionType: 'button',
    },
    defaultValue: 1,
    fieldName: 'status',
    label: $t('admin.role.status'),
  },
  {
    component: 'InputNumber',
    fieldName: 'sort',
    label: $t('common.sort'),
    componentProps: {
      placeholder: $t('admin.role.sortPlaceholder'),
      min: 0,
      precision: 0,
      style: { width: '100%' },
    },
    defaultValue: 50,
    rules: z.number().min(0, $t('admin.role.sortMinValue')),
  },
  {
    component: 'Input',
    fieldName: 'name',
    componentProps: {
      placeholder: $t('admin.role.roleNamePlaceholder'),
    },
    label: $t('admin.role.roleName'),
    rules: z
      .string()
      .min(2, $t('ui.formRules.minLength', [$t('admin.role.roleName'), 2]))
      .max(64, $t('ui.formRules.maxLength', [$t('admin.role.roleName'), 64])),
  },
  {
    component: 'Input',
    fieldName: 'description',
    label: $t('admin.role.description'),
    componentProps: {
      placeholder: $t('admin.role.descriptionPlaceholder'),
      rows: 3,
    },
    rules: z
      .string()
      .max(255, $t('ui.formRules.maxLength', [$t('admin.role.description'), 255]))
      .optional(),
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getRoleTypeOptions(),
      optionType: 'button',
    },
    defaultValue: 0,
    fieldName: 'is_super',
    label: $t('admin.role.type'),
  },
  
  {
    component: 'ApiTreeSelect',
    componentProps: {
      placeholder: $t('admin.role.parentPlaceholder'),
      api: getRoleList,
      class: 'w-full',
      filterTreeNode(input: string, node: Recordable<any>) {
        if (!input || input.length === 0) {
          return true;
        }
        const title: string = node.name ?? '';
        if (!title) return false;
        return title.includes(input);
      },
      getPopupContainer,
      labelField: 'name',
      showSearch: true,
      treeDefaultExpandAll: true,
      valueField: 'id',
      childrenField: 'children',
      allowClear: true,
    },
    fieldName: 'pid',
    label: $t('admin.role.parent'),
  },
  {
    component: 'Input',
    fieldName: 'menu_ids',
    formItemClass: 'col-span-2 md:col-span-2 items-start',
    label: $t('admin.role.permissions'),
    modelPropName: 'modelValue',
  },
];

const breakpoints = useBreakpoints(breakpointsTailwind);
const isHorizontal = computed(() => breakpoints.greaterOrEqual('md').value);

const [Form, formApi] = useVbenForm({
  commonConfig: {
    colon: true,
    formItemClass: 'col-span-2 md:col-span-2',
  },
  schema,
  showDefaultActions: false,
  wrapperClass: 'grid-cols-2 gap-x-4',
});

const [Drawer, drawerApi] = useVbenDrawer({
  onConfirm: onSubmit,
  async onOpenChange(isOpen) {
    if (isOpen) {
      const data = drawerApi.getData<AdminRoleApi.AdminRole>();
      formApi.resetForm();

      if (data) {
        formData.value = data;
      } else {
        formData.value = undefined;
      }

      if (permissions.value.length === 0) {
        await loadPermissions();
      }
      // Wait for Vue to flush DOM updates (form fields mounted)
      await nextTick();
      if (data) {
        // 处理上级角色显示：pid 为 0 时不显示任何值
        const formValues = {
          ...data,
          pid: data.pid === 0 ? undefined : data.pid
        };
        formApi.setValues(formValues);
      }
    }
  },
});

async function loadPermissions() {
  loadingPermissions.value = true;
  try {
    const res = await getMyMenus();
    permissions.value = res as unknown as DataNode[];
  } finally {
    loadingPermissions.value = false;
  }
}

async function onSubmit() {
  const { valid } = await formApi.validate();
  if (valid) {
    drawerApi.lock();
    const data =
      await formApi.getValues<
        Omit<AdminRoleApi.AdminRole, 'id' | 'created_at' | 'updated_at'>
      >();
    
    // 处理上级角色：如果没有选择上级角色，设置为 0
    if (!data.pid) {
      data.pid = 0;
    }
    
    // 确保 menu_ids 字段存在
    if (!data.menu_ids) {
      data.menu_ids = [];
    }
    
    try {
      // 保存角色（包含菜单权限）
      await saveRole(formData.value?.id || 0, data);
      
      drawerApi.close();
      emit('success');
    } finally {
      drawerApi.unlock();
    }
  }
}

const getDrawerTitle = computed(() =>
  formData.value?.id
    ? $t('ui.actionTitle.edit', [$t('admin.role.name')])
    : $t('ui.actionTitle.create', [$t('admin.role.name')]),
);

function getNodeClass(node: Recordable<any>) {
  const classes: string[] = [];
  // 检查节点类型，可能在 node.type 或 node.value?.type 中
  const nodeType = node.type || node.value?.type;
  if (nodeType === 'button') {
    classes.push('inline-flex', 'flex-wrap');
  }
  return classes.join(' ');
}

function getMenuTypeInfo(type: string) {
  const typeOptions = getMenuTypeOptions();
  return typeOptions.find(option => option.value === type) || { color: 'default', label: type };
}
</script>
<template>
  <Drawer class="w-full max-w-[600px]" :title="getDrawerTitle">
    <Form class="mx-4" :layout="isHorizontal ? 'horizontal' : 'vertical'">
      <template #menu_ids="slotProps">
        <Spin :spinning="loadingPermissions" wrapper-class-name="w-full">
          <Tree
            :tree-data="permissions"
            multiple
            checkable
            bordered
            :default-expanded-level="2"
            :get-node-class="getNodeClass"
            v-bind="slotProps"
            value-field="id"
            label-field="name"
            icon-field="icon"
            children-field="children"
          >
            <template #node="{ value }">
              <div class="flex items-center gap-2">
                <IconifyIcon v-if="value.icon" :icon="value.icon" class="size-4" />
                <span>{{ value.name }}</span>
                <Tag v-if="value.type" :color="getMenuTypeInfo(value.type).color" size="small">
                  {{ getMenuTypeInfo(value.type).label }}
                </Tag>
              </div>
            </template>
          </Tree>
        </Spin>
      </template>
    </Form>
  </Drawer>
</template>
<style lang="css" scoped>
:deep(.ant-tree-title) {
  .tree-actions {
    display: none;
    margin-left: 20px;
  }
}

:deep(.ant-tree-title:hover) {
  .tree-actions {
    display: flex;
    flex: auto;
    justify-content: flex-end;
    margin-left: 20px;
  }
}

/* 按钮类型的节点横向显示 */
:deep(.ant-tree-treenode.inline-flex) {
  display: inline-flex !important;
  flex-wrap: wrap;
  margin-right: 8px;
  margin-bottom: 4px;
}

/* 按钮类型节点的子节点容器 */
:deep(.ant-tree-treenode.inline-flex + .ant-tree-child-tree) {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

:deep(.ant-tree-child-tree .ant-tree-treenode.inline-flex) {
  display: inline-flex !important;
  margin-right: 8px;
  margin-bottom: 4px;
}
</style>