<script lang="ts" setup>
import type { Recordable } from '@vben/types';

import type { VbenFormSchema } from '#/adapter/form';

import { computed, nextTick, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { getPopupContainer } from '@vben/utils';

import { breakpointsTailwind, useBreakpoints } from '@vueuse/core';

import { useVbenForm, z } from '#/adapter/form';
import { saveAdmin } from '#/api/admin/admin';
import { getRoleList } from '#/api/admin/role';
import { $t } from '#/locales';

import type { AdminAdminApi } from '#/api/admin/admin';

import { getAdminStatusOptions } from '../data';

const emit = defineEmits<{
  success: [];
}>();

const formData = ref<AdminAdminApi.Admin>();

const schema: VbenFormSchema[] = [
  {
    component: 'RadioGroup',
    componentProps: {
      buttonStyle: 'solid',
      options: getAdminStatusOptions(),
      optionType: 'button',
    },
    defaultValue: 1,
    fieldName: 'status',
    label: $t('admin.admin.status'),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('admin.admin.usernamePlaceholder'),
    },
    fieldName: 'username',
    label: $t('admin.admin.username'),
    rules: z
      .string()
      .min(3, $t('ui.formRules.minLength', [$t('admin.admin.username'), 3]))
      .max(32, $t('ui.formRules.maxLength', [$t('admin.admin.username'), 32])),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('admin.admin.realnamePlaceholder'),
    },
    fieldName: 'realname',
    label: $t('admin.admin.realname'),
    rules: z
      .string()
      .min(2, $t('ui.formRules.minLength', [$t('admin.admin.realname'), 2]))
      .max(32, $t('ui.formRules.maxLength', [$t('admin.admin.realname'), 32]))
      .optional(),
  },
  {
    component: 'InputPassword',
    componentProps: {
      placeholder: $t('admin.admin.passwordPlaceholder'),
    },
    fieldName: 'password',
    label: $t('admin.admin.password'),
    rules: z
      .string()
      .min(6, $t('ui.formRules.minLength', [$t('admin.admin.password'), 6]))
      .max(32, $t('ui.formRules.maxLength', [$t('admin.admin.password'), 32]))
      .optional()
      .or(z.literal('')),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('admin.admin.emailPlaceholder'),
      type: 'email',
    },
    fieldName: 'email',
    label: $t('admin.admin.email'),
    rules: z
      .string()
      .email($t('admin.admin.invalidFormat', [$t('admin.admin.email')]))
      .optional()
      .or(z.literal('')),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('admin.admin.mobilePlaceholder'),
    },
    fieldName: 'mobile',
    label: $t('admin.admin.mobile'),
    rules: z
      .string()
      .regex(
        /^1[3-9]\d{9}$/,
        $t('admin.admin.invalidFormat', [$t('admin.admin.mobile')]),
      )
      .optional()
      .or(z.literal('')),
  },
  {
    component: 'AttachmentInput',
    componentProps: {
      multiple: false,
      placeholder: $t('admin.admin.avatarPlaceholder'),
      showPreview: true,
      showInput: false,
    },
    fieldName: 'avatar',
    formItemClass: 'col-span-2 md:col-span-2 items-baseline',
    label: $t('admin.admin.avatar'),
    rules: z.string().optional().or(z.literal('')),
  },
  {
    component: 'ApiTreeSelect',
    componentProps: {
      allowClear: true,
      api: getRoleList,
      childrenField: 'children',
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
      mode: 'multiple',
      placeholder: $t('admin.admin.rolesPlaceholder'),
      showSearch: true,
      treeCheckable: true,
      treeDefaultExpandAll: true,
      valueField: 'id',
    },
    fieldName: 'role_ids',
    formItemClass: 'col-span-2 md:col-span-2',
    label: $t('admin.admin.roles'),
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
      const data = drawerApi.getData<AdminAdminApi.Admin>();
      formApi.resetForm();

      if (data) {
        formData.value = data;
      } else {
        formData.value = undefined;
      }

      // Wait for Vue to flush DOM updates (form fields mounted)
      await nextTick();
      if (data) {
        // 设置表单值，不包含密码字段
        const formValues = {
          ...data,
          password: '', // 编辑时密码留空
        };
        formApi.setValues(formValues);
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
        Omit<
          AdminAdminApi.Admin,
          'id' | 'created_at' | 'updated_at' | 'last_login_time'
        >
      >();

    // 如果是新建且没有密码，提示错误
    if (!formData.value?.id && !data.password) {
      formApi.setFieldValue('password', '');
      drawerApi.unlock();
      return;
    }

    // 如果是编辑且密码为空，删除密码字段
    if (formData.value?.id && !data.password) {
      delete data.password;
    }

    // 确保 role_ids 字段存在
    if (!data.role_ids) {
      data.role_ids = [];
    }

    try {
      // 保存管理员
      await saveAdmin(formData.value?.id || 0, data);

      drawerApi.close();
      emit('success');
    } finally {
      drawerApi.unlock();
    }
  }
}

const getDrawerTitle = computed(() =>
  formData.value?.id
    ? $t('ui.actionTitle.edit', [$t('admin.admin.name')])
    : $t('ui.actionTitle.create', [$t('admin.admin.name')]),
);
</script>
<template>
  <Drawer class="w-full max-w-[600px]" :title="getDrawerTitle">
    <Form class="mx-4" :layout="isHorizontal ? 'horizontal' : 'vertical'" />
  </Drawer>
</template>

