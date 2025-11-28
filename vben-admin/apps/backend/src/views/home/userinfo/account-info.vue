<script lang="ts" setup>
import { onMounted } from 'vue';

import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { message } from 'ant-design-vue';

import { useVbenForm, z } from '#/adapter/form';
import { updateProfileApi } from '#/api/core/auth';

const userStore = useUserStore();

const schema = [
  {
    component: 'AttachmentInput',
    componentProps: {
      multiple: false,
      placeholder: $t('userinfo.account.avatarPlaceholder'),
      showInput: true,
      showPreview: true,
    },
    fieldName: 'avatar',
    formItemClass: 'items-baseline',
    label: $t('userinfo.account.avatar'),
    rules: z.string().url($t('ui.formRules.invalidFormat', [$t('userinfo.account.avatar')])).optional().or(z.literal('')),
  },
  {
    component: 'Input',
    componentProps: {
      disabled: true,
      placeholder: $t('userinfo.account.usernameDisabled'),
    },
    fieldName: 'username',
    label: $t('userinfo.account.username'),
  },
  {
    component: 'Input',
    componentProps: {
      placeholder: $t('userinfo.account.emailPlaceholder'),
      type: 'email',
    },
    fieldName: 'email',
    label: $t('userinfo.account.email'),
    rules: z.string().email($t('ui.formRules.invalidFormat', [$t('userinfo.account.email')])),
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

// 保存
const handleSubmit = async () => {
  const { valid } = await formApi.validate();
  if (!valid) return;

  const values = await formApi.getValues<{ avatar: string; email: string }>();

  const hideLoading = message.loading($t('userinfo.account.saving'), 0);
  try {
    await updateProfileApi(values);
    message.success($t('userinfo.account.saveSuccess'));

    // 更新 store 中的用户信息
    if (userStore.userInfo) {
      userStore.setUserInfo({
        ...userStore.userInfo,
        avatar: values.avatar,
        email: values.email,
      });
    }

    hideLoading();
  } catch {
    hideLoading();
  }
};

onMounted(() => {
  // 设置初始值
  if (userStore.userInfo) {
    formApi.setValues({
      avatar: userStore.userInfo.avatar || '',
      email: userStore.userInfo.email || '',
      username: userStore.userInfo.username || '',
    });
  }
});
</script>

<template>
  <div>
    <h2 class="mb-6 text-xl font-semibold">{{ $t('userinfo.account.title') }}</h2>
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

