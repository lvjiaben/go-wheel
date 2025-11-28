<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, useTemplateRef } from 'vue';

import { AuthenticationCodeLogin, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from 'ant-design-vue';

import { useAuthStore } from '#/store';

defineOptions({ name: 'CodeLogin' });

const authStore = useAuthStore();
const CODE_LENGTH = 6;

const loginRef =
  useTemplateRef<InstanceType<typeof AuthenticationCodeLogin>>('loginRef');

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.mobile'),
      },
      fieldName: 'mobile',
      label: $t('authentication.mobile'),
      rules: z
        .string()
        .min(1, { message: $t('authentication.mobileTip') })
        .refine((v) => /^\d{11}$/.test(v), {
          message: $t('authentication.mobileErrortip'),
        }),
    },
    {
      component: 'VbenPinInput',
      componentProps: {
        codeLength: CODE_LENGTH,
        createText: (countdown: number) => {
          const text =
            countdown > 0
              ? $t('authentication.sendText', [countdown])
              : $t('authentication.sendCode');
          return text;
        },
        handleSendCode: async () => {
          const formApi = loginRef.value?.getFormApi();
          if (!formApi) {
            throw new Error('formApi is not ready');
          }
          await formApi.validateField('mobile');
          const isMobileReady = await formApi.isFieldValid('mobile');
          if (!isMobileReady) {
            throw new Error('Mobile number is not valid');
          }
          const { mobile } = await formApi.getValues();
          await authStore.sendSmsCode(mobile, 'login');
          message.success($t('page.auth.smsSentSuccess'));
        },
        placeholder: $t('authentication.code'),
      },
      fieldName: 'code',
      label: $t('authentication.code'),
      rules: z.string().length(CODE_LENGTH, {
        message: $t('authentication.codeTip', [CODE_LENGTH]),
      }),
    },
  ];
});

/**
 * 异步处理登录操作
 * Asynchronously handle the login process
 * @param values 登录表单数据
 */
async function handleLogin(values: Recordable<any>) {
  await authStore.authMobileLogin({
    mobile: values.mobile,
    code: values.code,
  });
}
</script>

<template>
  <AuthenticationCodeLogin
    ref="loginRef"
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    @submit="handleLogin"
  />
</template>
