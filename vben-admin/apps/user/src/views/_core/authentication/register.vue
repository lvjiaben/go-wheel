<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, h, useTemplateRef } from 'vue';

import { AuthenticationRegister, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { message } from 'ant-design-vue';

import { useAuthStore } from '#/store';

defineOptions({ name: 'Register' });

const authStore = useAuthStore();
const CODE_LENGTH = 6;

const registerRef =
  useTemplateRef<InstanceType<typeof AuthenticationRegister>>('registerRef');

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
          const formApi = registerRef.value?.getFormApi();
          if (!formApi) {
            throw new Error('formApi is not ready');
          }
          await formApi.validateField('mobile');
          const isMobileReady = await formApi.isFieldValid('mobile');
          if (!isMobileReady) {
            throw new Error('Mobile number is not valid');
          }
          const { mobile } = await formApi.getValues();
          await authStore.sendSmsCode(mobile, 'register');
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
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
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
          const { password } = values;
          return z
            .string({ required_error: $t('authentication.passwordTip') })
            .min(1, { message: $t('authentication.passwordTip') })
            .refine((value) => value === password, {
              message: $t('authentication.confirmPasswordTip'),
            });
        },
        triggerFields: ['password'],
      },
      fieldName: 'confirmPassword',
      label: $t('authentication.confirmPassword'),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('page.auth.inviteCodeTip'),
      },
      fieldName: 'invite_code',
      label: $t('page.auth.inviteCode'),
    },
    {
      component: 'VbenCheckbox',
      fieldName: 'agreePolicy',
      renderComponentContent: () => ({
        default: () =>
          h('span', [
            $t('authentication.agree'),
            h(
              'a',
              {
                class: 'vben-link ml-1 ',
                href: '',
              },
              `${$t('authentication.privacyPolicy')} & ${$t('authentication.terms')}`,
            ),
          ]),
      }),
      rules: z.boolean().refine((value) => !!value, {
        message: $t('authentication.agreeTip'),
      }),
    },
  ];
});

async function handleSubmit(value: Recordable<any>) {
  await authStore.authRegister({
    mobile: value.mobile,
    password: value.password,
    code: value.code,
    invite_code: value.invite_code,
  });
}
</script>

<template>
  <AuthenticationRegister
    ref="registerRef"
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    @submit="handleSubmit"
  />
</template>
