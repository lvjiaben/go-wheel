<script lang="ts" setup>

import type { Recordable } from '@vben/types';

import type { VbenFormSchema } from '#/adapter/form';

import { computed, h, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { IconifyIcon } from '@vben/icons';
import { getPopupContainer } from '@vben/utils';

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import { useVbenForm, z } from '#/adapter/form';
import {
  saveMenu,
  getMenuList,
  AdminMenuApi,
} from '#/api/admin/menu';
import { $t } from '#/locales';
import { componentKeys } from '#/router/routes';

import { getMenuTypeOptions } from '../data';

const emit = defineEmits<{
  success: [];
}>();
const formData = ref<AdminMenuApi.AdminMenu>();
const schema: VbenFormSchema[] = [
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getMenuTypeOptions(),
      optionType: 'button',
    },
    defaultValue: 'menu',
    fieldName: 'type',
    formItemClass: 'col-span-2 md:col-span-2',
    label: $t('admin.menu.type'),
  },
  {
    component: 'Input',
    fieldName: 'name',
    componentProps: {
      placeholder: $t('admin.menu.menuNamePlaceholder'),
    },
    label: $t('admin.menu.menuName'),
    rules: z
      .string()
      .min(2, $t('ui.formRules.minLength', [$t('admin.menu.menuName'), 2]))
      .max(30, $t('ui.formRules.maxLength', [$t('admin.menu.menuName'), 30]))
      ,
  },
  {
    component: 'Input',
    fieldName: 'enname',
    label: $t('admin.menu.menuNameEn'),
    componentProps: {
      placeholder: $t('admin.menu.menuNameEnPlaceholder'),
    },
    rules: z
      .string()
      .min(2, $t('ui.formRules.minLength', [$t('admin.menu.menuNameEn'), 2]))
      .max(30, $t('ui.formRules.maxLength', [$t('admin.menu.menuNameEn'), 30]))
      ,
  },
  {
    component: 'Input',
    fieldName: 'sort',
    label: $t('common.sort'),
    componentProps: {
      placeholder: $t('admin.menu.sortPlaceholder'),
    },
    rules: z.coerce.number().int($t('admin.menu.sortIntRequired'))
  },
  {
    component: 'ApiTreeSelect',
    componentProps: {
      placeholder: $t('admin.menu.parentPlaceholder'),
      api: getMenuList,
      class: 'w-full',
      filterTreeNode(input: string, node: Recordable<any>) {
        if (!input || input.length === 0) {
          return true;
        }
        const title: string = node.name ?? '';
        if (!title) return false;
        return title.includes(input) || $t(title).includes(input);
      },
      getPopupContainer,
      labelField: 'name',
      showSearch: true,
      treeDefaultExpandAll: true,
      valueField: 'id',
      childrenField: 'children',
    },
    fieldName: 'pid',
    label: $t('admin.menu.parent'),
    renderComponentContent() {
      return {
        title({ label, icon }: { label: string; icon?: string }) {
          console.log(label, icon);
          const coms = [];
          if (!label) return '';
          if (icon) {
            coms.push(h(IconifyIcon, { class: 'size-4', icon: icon }));
          }
          coms.push(h('span', { class: '' }, $t(label || '')));
          return h('div', { class: 'flex items-center gap-1' }, coms);
        },
      };
    },
  },
  {
    component: 'Input',
    dependencies: {
      show: (values) => {
        return ['menu'].includes(values.type);
      },
      triggerFields: ['type'],
    },
    componentProps: {
      placeholder: $t('admin.menu.pathPlaceholder'),
    },
    fieldName: 'path',
    label: $t('admin.menu.path'),
    rules: z
      .string()
      .min(2, $t('ui.formRules.minLength', [$t('admin.menu.path'), 2]))
      .max(100, $t('ui.formRules.maxLength', [$t('admin.menu.path'), 100]))
      .refine(
        (value: string) => {
          return value.startsWith('/');
        },
        $t('ui.formRules.startWith', [$t('admin.menu.path'), '/']),
      )
      ,
  },
  {
    component: 'IconPicker',
    componentProps: {
      prefix: 'carbon',
      placeholder: $t('admin.menu.iconPlaceholder'),
    },
    fieldName: 'icon',
    
    label: $t('admin.menu.icon'),
  },
  {
    component: 'AutoComplete',
    componentProps: {
      placeholder: $t('admin.menu.componentPlaceholder'),
      class: 'w-full',
      filterOption(input: string, option: { value: string }) {
        return option.value.toLowerCase().includes(input.toLowerCase());
      },
      options: componentKeys.map((v) => ({ value: v })),
    },
    dependencies: {
      rules: (values) => {
        return values.type === 'menu' ? 'required' : null;
      },
      show: (values) => {
        return values.type === 'menu';
      },
      triggerFields: ['type'],
    },
    fieldName: 'component',
    label: $t('admin.menu.component'),
  },
  {
    component: 'Input',
    dependencies: {
      rules: (values) => {
        return ['iframe'].includes(values.type) ? 'required' : null;
      },
      show: (values) => {
        return ['iframe'].includes(values.type);
      },
      triggerFields: ['type'],
    },
    componentProps: {
      placeholder: $t('admin.menu.iframePlaceholder'),
    },
    fieldName: 'iframe',
    label: $t('admin.menu.linkSrc'),
    rules: z.string().url($t('ui.formRules.invalidURL')),
  },
  {
    component: 'Input',
    dependencies: {
      rules: (values) => {
        return ['link'].includes(values.type) ? 'required' : null;
      },
      show: (values) => {
        return ['link'].includes(values.type);
      },
      triggerFields: ['type'],
    },
    componentProps: {
      placeholder: 'HTTP外部链接地址',
    },
    fieldName: 'external',
    label: $t('admin.menu.linkSrc'),
    rules: z.string().url($t('ui.formRules.invalidURL')),
  },
  {
    component: 'Input',
    dependencies: {
      rules: (values) => {
        return values.type === 'button' ? 'required' : null;
      },
      show: (values) => {
        return values.type === 'button';
      },
      triggerFields: ['type'],
    },
    componentProps: {
      placeholder: '权限标识，用于前端访问控制',
    },
    fieldName: 'permission',
    label: $t('admin.menu.authCode'),
  },
  {
    component: 'Input',
    dependencies: {
      rules: (values) => {
        return values.type === 'button' ? 'required' : null;
      },
      show: (values) => {
        return values.type === 'button';
      },
      triggerFields: ['type'],
    },
    componentProps: {
      placeholder: '后端API接口地址，用于后端鉴权',
    },
    fieldName: 'route',
    label: $t('admin.menu.linkSrc'),
    rules: z.string().url($t('ui.formRules.invalidURL')),
  },
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: [
        { label: $t('common.show'), value: 1 },
        { label: $t('common.hide'), value: 0 },
      ],
      optionType: 'button',
    },
    defaultValue: 1,
    fieldName: 'visible',
    label: $t('admin.menu.status'),
  },
  {
    component: 'Divider',
    dependencies: {
      show: (values) => {
        return !['button', 'link'].includes(values.type);
      },
      triggerFields: ['type'],
    },
    fieldName: 'divider1',
    formItemClass: 'col-span-2 md:col-span-2 pb-0',
    hideLabel: true,
    renderComponentContent() {
      return {
        default: () => $t('admin.menu.advancedSettings'),
      };
    },
  },
  
  {
    component: 'Switch',
    componentProps: {
      // 选中时的值
      checkedValue: 1,
      // 未选中时的值
      unCheckedValue: 0,
    },
    defaultValue: 0,
    dependencies: {
      show: (values) => {
        return !['button', 'link'].includes(values.type);
      },
      triggerFields: ['type'],
    },
    fieldName: 'fixed_tag',
    label: $t('admin.menu.affixTab'),
  },
  {
    component: 'Switch',
    componentProps: {
      // 选中时的值
      checkedValue: 1,
      // 未选中时的值
      unCheckedValue: 0,
    },
    defaultValue: 0,
    dependencies: {
      show: (values) => {
        return !['button', 'link'].includes(values.type);
      },
      triggerFields: ['type'],
    },
    fieldName: 'show_tag',
    label: $t('admin.menu.hideInTab'),
  },
];

const breakpoints = useBreakpoints(breakpointsTailwind);
const isHorizontal = computed(() => breakpoints.greaterOrEqual('md').value);

const [Form, formApi] = useVbenForm({
  commonConfig: {
    colon: true,
    formItemClass: 'col-span-2 md:col-span-1',
  },
  schema,
  showDefaultActions: false,
  wrapperClass: 'grid-cols-2 gap-x-4',
});

const [Drawer, drawerApi] = useVbenDrawer({
  onConfirm: onSubmit,
  onOpenChange(isOpen) {
    if (isOpen) {
      const data = drawerApi.getData<AdminMenuApi.AdminMenu>();
      if (data) {
        formData.value = data;
        formApi.setValues(formData.value);
      } else {
        formApi.resetForm();
      }
    }
  },
});

async function onSubmit() {
  const { valid } = await formApi.validate();
  if (valid) {
    drawerApi.lock();
    const data =
      await formApi.getValues<
        Omit<AdminMenuApi.AdminMenu, 'children' | 'id'>
      >();
    
    // 转换数据类型
    
    try {
      await saveMenu(formData.value?.id, data);
      drawerApi.close();
      emit('success');
    } finally {
      drawerApi.unlock();
    }
  }
}
const getDrawerTitle = computed(() =>
  formData.value?.id
    ? $t('ui.actionTitle.edit', [$t('admin.menu.name')])
    : $t('ui.actionTitle.create', [$t('admin.menu.name')]),
);
</script>
<template>
  <Drawer class="w-full max-w-[800px]" :title="getDrawerTitle">
    <Form class="mx-4" :layout="isHorizontal ? 'horizontal' : 'vertical'" />
  </Drawer>
</template>
