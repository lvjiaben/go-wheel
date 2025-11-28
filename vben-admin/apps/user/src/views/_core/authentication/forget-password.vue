<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, useTemplateRef } from 'vue';

import { AuthenticationForgetPassword, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from 'ant-design-vue';

import { useAuthStore } from '#/store';

defineOptions({ name: 'ForgetPassword' });

const authStore = useAuthStore();
const CODE_LENGTH = 6;

const forgetRef =
  useTemplateRef<InstanceType<typeof AuthenticationForgetPassword>>(
    'forgetRef',
  );

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
          const formApi = forgetRef.value?.getFormApi();
          if (!formApi) {
            throw new Error('formApi is not ready');
          }
          await formApi.validateField('mobile');
          const isMobileReady = await formApi.isFieldValid('mobile');
          if (!isMobileReady) {
            throw new Error('Mobile number is not valid');
          }
          const { mobile } = await formApi.getValues();
          await authStore.sendSmsCode(mobile, 'resetpwd');
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
    {
      component: 'VbenInputPassword',
      componentProps: {
        passwordStrength: true,
        placeholder: $t('page.auth.newPassword'),
      },
      fieldName: 'new_password',
      label: $t('page.auth.newPassword'),
      renderComponentContent() {
        return {
          strengthText: () => $t('authentication.passwordStrength'),
        };
      },
      rules: z
        .string()
        .min(6, { message: $t('page.auth.passwordMinLength') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.confirmPassword'),
      },
      dependencies: {
        rules(values) {
          const { new_password } = values;
          return z
            .string({ required_error: $t('authentication.passwordTip') })
            .min(1, { message: $t('authentication.passwordTip') })
            .refine((value) => value === new_password, {
              message: $t('authentication.confirmPasswordTip'),
            });
        },
        triggerFields: ['new_password'],
      },
      fieldName: 'confirmPassword',
      label: $t('authentication.confirmPassword'),
    },
  ];
});

async function handleSubmit(value: Recordable<any>) {
  await authStore.authResetPwd({
    mobile: value.mobile,
    code: value.code,
    new_password: value.new_password,
  });
}
</script>

<template>
  <AuthenticationForgetPassword
    ref="forgetRef"
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    :submit-button-text="$t('page.auth.resetPassword')"
    :sub-title="$t('page.auth.forgetPasswordSubtitle')"
    @submit="handleSubmit"
  />
</template>
