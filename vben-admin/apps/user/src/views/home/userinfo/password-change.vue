<script lang="ts" setup>
import { ref } from 'vue';

import { VbenButton, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { changePwdApi } from '#/api/core/user';

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

const loading = ref(false);

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

  loading.value = true;
  try {
    await changePwdApi({
      new_password: values.new_password,
      old_password: values.old_password,
    });
    message.success($t('userinfo.password.changeSuccess'));
    // 清空表单
    formApi.resetForm();
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div>
    <h2 class="mb-6 text-xl font-semibold">
      {{ $t('userinfo.password.title') }}
    </h2>
    <Form class="max-w-2xl" />
    <div class="mt-6 flex max-w-2xl">
      <div class="flex-1"></div>
      <VbenButton :loading="loading" @click="handleSubmit">
        {{ $t('userinfo.account.save') }}
      </VbenButton>
    </div>
  </div>
</template>

