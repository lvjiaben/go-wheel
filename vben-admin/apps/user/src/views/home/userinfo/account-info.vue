<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';

import { computed, onMounted, ref } from 'vue';

import { VbenButton, z } from '@vben/common-ui';
import { $t } from '@vben/locales';
import { useUserStore } from '@vben/stores';

import { message } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { changeMobileApi, sendSmsApi } from '#/api/core/user';

const userStore = useUserStore();
const CODE_LENGTH = 6;

const schema = computed((): VbenFormSchema[] => [
  {
    component: 'Input',
    componentProps: {
      disabled: true,
      placeholder: $t('userinfo.mobile.currentMobilePlaceholder'),
    },
    fieldName: 'mobile',
    label: $t('userinfo.mobile.currentMobile'),
  },
  {
    component: 'VbenPinInput',
    componentProps: {
      codeLength: CODE_LENGTH,
      createText: (countdown: number) => {
        return countdown > 0
          ? $t('authentication.sendText', [countdown])
          : $t('authentication.sendCode');
      },
      handleSendCode: async () => {
        // 发送验证码到当前手机号（后端会从登录用户获取）
        await sendSmsApi({ mobile: '', event: 'changemobile' });
        message.success($t('page.auth.smsSentSuccess'));
      },
      placeholder: $t('userinfo.mobile.codePlaceholder'),
    },
    fieldName: 'code',
    label: $t('userinfo.mobile.code'),
    rules: z.string().length(CODE_LENGTH, {
      message: $t('userinfo.mobile.codeRequired'),
    }),
  },
  {
    component: 'Input',
    componentProps: {
      maxlength: 11,
      placeholder: $t('userinfo.mobile.newMobilePlaceholder'),
    },
    fieldName: 'new_mobile',
    label: $t('userinfo.mobile.newMobile'),
    rules: z
      .string()
      .min(1, { message: $t('userinfo.mobile.newMobileRequired') })
      .refine((v) => /^\d{11}$/.test(v), {
        message: $t('userinfo.mobile.mobileFormatError'),
      }),
  },
]);

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
    code: string;
    mobile: string;
    new_mobile: string;
  }>();

  loading.value = true;
  try {
    await changeMobileApi(values);
    message.success($t('userinfo.mobile.changeSuccess'));

    // 更新 store 中的用户信息
    if (userStore.userInfo) {
      userStore.setUserInfo({
        ...userStore.userInfo,
        mobile: values.new_mobile,
      });
    }

    // 重置表单
    formApi.setValues({
      code: '',
      mobile: values.new_mobile,
      new_mobile: '',
    });
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  // 设置初始值
  if (userStore.userInfo) {
    formApi.setValues({
      mobile: (userStore.userInfo as any).mobile || '',
    });
  }
});
</script>

<template>
  <div>
    <h2 class="mb-6 text-xl font-semibold">
      {{ $t('userinfo.mobile.title') }}
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

