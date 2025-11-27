<script lang="ts" setup>
import { $t } from '@vben/locales';

import { message } from 'ant-design-vue';

import { useVbenForm, z } from '#/adapter/form';
import { changePasswordApi } from '#/api/core/auth';

const schema = [
  {
    component: 'InputPassword',
    componentProps: {
      placeholder: $t('userinfo.password.oldPasswordPlaceholder'),
    },
    fieldName: 'old_password',
    label: $t('userinfo.password.oldPassword'),
    rules: z.string().min(1, $t('userinfo.password.oldPasswordRequired')),
  },
  {
    component: 'InputPassword',
    componentProps: {
      placeholder: $t('userinfo.password.newPasswordPlaceholder'),
    },
    fieldName: 'new_password',
    label: $t('userinfo.password.newPassword'),
    rules: z.string().min(6, $t('userinfo.password.passwordMinLength')),
  },
  {
    component: 'InputPassword',
    componentProps: {
      placeholder: $t('userinfo.password.confirmPasswordPlaceholder'),
    },
    fieldName: 'confirm_password',
    label: $t('userinfo.password.confirmPassword'),
    rules: z.string().min(6, $t('userinfo.password.passwordMinLength')),
  },
];

const [Form, formApi] = useVbenForm({
  commonConfig: {
    colon: true,
  },
  layout: 'horizontal',
  schema,
  showDefaultActions: false,
});

// 提交修改
const handleSubmit = async () => {
  const { valid } = await formApi.validate();
  if (!valid) return;

  const values = await formApi.getValues<{
    confirm_password: string;
    new_password: string;
    old_password: string;
  }>();

  // 验证两次密码是否一致
  if (values.new_password !== values.confirm_password) {
    message.error($t('userinfo.password.passwordMismatch'));
    return;
  }

  const hideLoading = message.loading($t('userinfo.password.changing'), 0);
  try {
    await changePasswordApi({
      new_password: values.new_password,
      old_password: values.old_password,
    });
    message.success($t('userinfo.password.changeSuccess'));
    // 清空表单
    formApi.resetForm();
    hideLoading();
  } catch {
    hideLoading();
  }
};
</script>

<template>
  <div>
    <h2 class="mb-6 text-xl font-semibold">{{ $t('userinfo.password.title') }}</h2>
    <Form class="max-w-2xl" />
    <div class="mt-6 flex max-w-2xl">
      <div class="flex-1"></div>
      <button
        class="rounded bg-blue-600 px-6 py-2 text-white hover:bg-blue-700"
        type="button"
        @click="handleSubmit"
      >
        {{ $t('userinfo.account.save') }}
      </button>
    </div>
  </div>
</template>

